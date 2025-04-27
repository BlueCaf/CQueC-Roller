package main

import (
	"fmt"
	"log"
	"os"
)

var (
	LogFile     *os.File
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func OpenLogFile(filename string) {
	if filename == "" {
		filename = "cquec_rolloer.log"
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("Log file not found:", filename)
		err = os.WriteFile(filename, []byte{}, 0644)
		if err != nil {
			fmt.Printf("Error createing file: %v\n", err)
			panic(fmt.Errorf("error createing file: %v", err))
		}
	}

	var err error
	LogFile, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}

	// 각 레벨별 로거 만들기
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarnLogger = log.New(LogFile, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func CloseLogFile() {
	if LogFile != nil {
		LogFile.Close()
	}
}
