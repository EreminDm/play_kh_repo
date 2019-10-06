package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ReadAPPConfig take params to start
func ReadAPPConfig() {
	//	PrintToLog("Start reading application config file")
	filename := "../conf/application.json"
	plan, _ := ioutil.ReadFile(filename)
	plan = bytes.TrimPrefix(plan, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
	err := json.Unmarshal(plan, &AppConfigO)
	if err != nil {
		errString := `Error configuration.go  ReadAPPConfig()` + err.Error()
		writeFatalToLogFile(errString)
		fmt.Println(errString)

	} else {
		PrintlnToLogFile("Application config readed successfully")
	}
}

//ReadDBConfig func to set params to connection to DB
func ReadDBConfig() {
	filename := "../conf/dbconfig.json"
	linuxfileNname := "./conf/dbconfig.json"

	plan, _ := ioutil.ReadFile(linuxfileNname)
	if len(plan) == 0 {
		plan, _ = ioutil.ReadFile(filename)
	}
	plan = bytes.TrimPrefix(plan, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}
	err = json.Unmarshal(plan, &DBConfigInterface)
	if err != nil {
		errstr := "Cannot unmarshal the json" + err.Error()
		writeFatalToLogFile(errstr)
		fmt.Println(errstr)
	}
	fmt.Println(DBConfigInterface)
	fmt.Println(DBConfigInterface.DBHost, DBConfigInterface.DBName, DBConfigInterface.DBPass, DBConfigInterface.DBPort, DBConfigInterface.DBUser)

}

//CreateDirIfNotExist creating folder if it isn't exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			writeFatalToLogFile("error on create folder " + err.Error())
		} else {
			PrintlnToLogFile("Folder created: " + dir)
		}
	}
}
