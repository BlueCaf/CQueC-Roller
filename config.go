package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LogFile    string `yaml:"log_file"`
	RedisAddr  string `yaml:"redis_addr"`
	QueueKey   string `yaml:"queue_key"`
	ProcessKey string `yaml:"process_key"`
	TTL        string `yaml:"ttl"`
	UserCount  int    `yaml:"user_count"`
}

func LoadConfig(filename string) (Config, error) {
	var config Config
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Config file not found:", filename)
		defaultConfig := Config{
			LogFile:    "cquec_roller.log",
			RedisAddr:  "localhost:6379",
			QueueKey:   "queue:waiting:zset",
			ProcessKey: "queue:processed:",
			TTL:        "30m",
			UserCount:  10,
		}

		data, err := yaml.Marshal(&defaultConfig)
		if err != nil {
			fmt.Printf("failed to read config file: %v\n", err)
			return config, fmt.Errorf("failed to read config file: %w", err)
		}

		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Printf("failed to read config file: %v\n", err)
			return config, fmt.Errorf("failed to read config file: %w", err)
		}

		fmt.Printf("Create config file: %s.\n", filename)
	}

	f, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(f, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return config, nil
}
