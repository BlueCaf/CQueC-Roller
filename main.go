package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type QueueProcessor struct {
	redisClient *redis.Client
	queueKey    string
	processKey  string
	ttl         time.Duration
	ticker      *time.Ticker
	userCount   int
	stopChan    chan struct{}
}

func NewQueueProcessor(redisClient *redis.Client, queueKey string, processKey string, ttl time.Duration, userCount int) *QueueProcessor {
	qp := &QueueProcessor{
		redisClient: redisClient,
		queueKey:    queueKey,
		processKey:  processKey,
		ttl:         ttl,
		userCount:   userCount,
		stopChan:    make(chan struct{}),
	}
	qp.ticker = time.NewTicker(1 * time.Second)
	return qp
}

func (qp *QueueProcessor) Start() {
	go func() {
		for {
			select {
			case <-qp.ticker.C:
				qp.ProcessQueue(context.Background(), qp.userCount)
			case <-qp.stopChan:
				qp.ticker.Stop()
				InfoLogger.Println("QueueProcessor Stopped...")
				return
			}
		}
	}()
}

func (qp *QueueProcessor) Stop() {
	close(qp.stopChan)
}

func (qp *QueueProcessor) ProcessQueue(ctx context.Context, userCount int) {
	// debug, _ := qp.redisClient.ZRangeWithScores(ctx, qp.queueKey, 0, -1).Result()
	// fmt.Println(debug)
	items, err := qp.redisClient.ZRangeWithScores(ctx, qp.queueKey, 0, int64(userCount-1)).Result()
	if err != nil {
		ErrorLogger.Println("ZRangeWithScores 오류:", err)
		return
	}

	if len(items) != 0 {
		// minScore := items[0].Score
		// maxScore := items[len(items)-1].Score

		// err = qp.redisClient.ZRem(ctx, qp.queueKey, strconv.FormatFloat(minScore, 'f', -1, 64), strconv.FormatFloat(maxScore, 'f', -1, 64)).Err()
		// if err != nil {
		// 	ErrorLogger.Println("ZRemRangeByScore 오류:", err)
		// 	return
		// }

		for _, item := range items {
			err = qp.redisClient.ZRem(ctx, qp.queueKey, item.Member).Err()
			if err != nil {
				ErrorLogger.Println("ZRem error:", err)
				return
			}
		}

		var wg sync.WaitGroup
		for _, item := range items {
			wg.Add(1)
			go func(item redis.Z) {
				defer wg.Done()
				key := fmt.Sprintf("%s%s", qp.processKey, item.Member)
				// value := fmt.Sprintf("%.0f", item.Score)
				// err := qp.redisClient.Set(ctx, key, value, qp.ttl).Err()
				err := qp.redisClient.Set(ctx, key, "1", qp.ttl).Err()
				if err != nil {
					ErrorLogger.Println("Set 오류: ", err.Error())
				} else {
					InfoLogger.Println("Process된 항목 확인:", key, item.Score)
				}
			}(item)
		}
		wg.Wait()
	}
}

func main() {
	ctx := context.Background()

	// Load configuration from YAML file
	var config Config
	config, err := LoadConfig("config.yaml")
	if err != nil {
		InfoLogger.Println("Error loading config:", err)
		return
	}

	OpenLogFile(config.LogFile)
	defer CloseLogFile()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		ErrorLogger.Println("Redis에 연결하는 데 오류가 발생했습니다:", err)
		return
	}
	InfoLogger.Println("Redis에 연결되었습니다:", pong)

	queueKey := config.QueueKey
	processKey := config.ProcessKey
	ttl, err := time.ParseDuration(config.TTL)
	if err != nil {
		ErrorLogger.Println("Invalid TTL format:", err)
		return
	}

	userCount := config.UserCount

	processor := NewQueueProcessor(redisClient, queueKey, processKey, ttl, userCount)
	processor.Start()

	select {}
}
