package main

import (
	"fmt"
	"log"
	"os"
)

var loggers map[string]*log.Logger

var logFiles []*os.File

func InitLog(fileNames []string) {
	logFiles = []*os.File{}
	loggers = map[string]*log.Logger{}
	for _, fileName := range fileNames {
		logFile, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln("open file error!")
		}
		logFiles = append(logFiles, logFile)
		loggers[fileName] = log.New(logFile, "", log.LstdFlags)
	}
}

func CloseLog() {
	fmt.Printf("Log Closed\n")
	for _, logFile := range logFiles {
		logFile.Close()
	}
}

func WriteLog(fileName, log string) {
	logger := loggers[fileName]
	logger.Println(log)
}