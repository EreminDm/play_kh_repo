package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

var db *sql.DB
var err error

func init() {
	ReadDBConfig()
	ConnectTODatabase()
	ReadAPPConfig()
	go CleanDB()
}

func main() {
	appport := strconv.Itoa(AppConfigO.Port)
	http.HandleFunc("/sis/", welcomeInfo)
	http.HandleFunc("/sis/plate", newPlate)
	http.HandleFunc("/sis/commingPlate", commingPlate)
	ListenAndServeErr := http.ListenAndServe(":"+appport, nil)
	if ListenAndServeErr != nil {
		writeFatalToLogFile("Port is in use, please try another port" + ListenAndServeErr.Error())
		fmt.Println("Port is in use, please try another port" + ListenAndServeErr.Error())
	}
}
