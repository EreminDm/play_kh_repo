package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NavisIntegration(ansId int, plate string) (err error, destination string) {
	destination, err = SelectDestFromLocalDB(plate)
	if err != nil {
		return err, ``
	}
	return nil, destination
}

func PostTORoutingService(ansId int, plate string, destination string, entrance int, plate_numbe string) (err error) {
	//entrance = 1 - въезд; 0 - выезд;
	jsonstring := `{"entrance":` + strconv.Itoa(entrance) + `, "plateNumber":"` + plate + `", "uniqueAnshlagId":` + strconv.Itoa(ansId) + `, "destination":"` + destination + `"}`
	fmt.Println(jsonstring)
	jsonbyte := []byte(jsonstring)
	req, err := http.NewRequest("POST", AppConfigO.RoutingServURL, bytes.NewBuffer(jsonbyte))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		PrintlnToLogFile(`Couldn't make request function: ` + err.Error())
		fmt.Println(`Couldn't make request function: `, err)
		return err
	}
	defer resp.Body.Close()
	// need to test
	if resp.StatusCode != 200 {
		err1 := errors.New("Houston, we have a problem with server, code: " + resp.Status)
		PrintlnToLogFile(err1.Error())
		fmt.Println(err1)
		return err1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		PrintlnToLogFile(`Couldn't handle responce body: ` + err.Error())
		return err
	}

	err = json.Unmarshal(body, &RoutingDataO)
	if err != nil {
		fmt.Println(`Couldn't unmarshal json of Routing data: `, err)
		PrintlnToLogFile(`Couldn't unmarshal json of Routing data: ` + err.Error())
		return err
	}
	//fmt.Println(RoutingDataO.IsLastAns)
	err = UpdateDBinformation(plate_numbe, RoutingDataO.Direction, RoutingDataO.DestinationPoint, RoutingDataO.IsLastAns, ansId)
	if err != nil {
		fmt.Println(`Couldn't update information of vehicle routes: `, err)
		PrintlnToLogFile(`Couldn't update information of vehicle routes: ` + err.Error())
		return err
	}

	return nil
}
