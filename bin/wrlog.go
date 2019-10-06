package main

import (
	"log"
	"os"
	"time"
)

func writeFatalToLogFile(logInfo string) {
	CreateDirIfNotExist("../log/errors")
	t := time.Now().Format("2006-01-02")
	f, err := os.OpenFile("../log/errors/log_"+t+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(logInfo)
}
func writeToMAINLogFile(logINfo string) {
	CreateDirIfNotExist("../log/software_log")
	t := time.Now().Format("2006-01-02")
	f, err := os.OpenFile("../log/software_log/software_log_"+t+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(logINfo)
}
func PrintlnToLogFile(logInfo string) {
	CreateDirIfNotExist("../log/everyday_log")
	t := time.Now().Format("2006-01-02")
	f, err := os.OpenFile("../log/everyday_log/everyday_log_"+t+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(logInfo)
}
