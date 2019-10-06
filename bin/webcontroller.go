package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func welcomeInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `Automarshal integration service, version: 0.1`)
}

func newPlate(w http.ResponseWriter, r *http.Request) {
	//decoder := json.NewDecoder(r.Body)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(`Body parrse err: `, err)

		return
	}
	fmt.Println(string(body))
	// fmt.Println(string(bytes.TrimRightFunc(body, []byte("\r")))
	// fmt.Println(string(bytes.TrimRightFunc(body, unicode.IsPunct)))
	// fmt.Println(string(bytes.TrimRightFunc(body, unicode.IsNumber)))
	//body = bytes.TrimPrefix(body, []byte("\r")) // Or []byte{239, 187, 191}

	//err := decoder.Decode(&InputNewPlateO)
	err = json.Unmarshal(body, &InputNewPlateO)
	if err != nil {
		PrintlnToLogFile(`Couldn't parse plate number information json: ` + err.Error())
		fmt.Println(err)

		return
	}
	defer r.Body.Close()

	var destination string
	//Проверка на наличие номера в базе данных
	plateNumber := cutnumber(InputNewPlateO.Plate)
	have := istherenumber(plateNumber)

	fullPlate, err := SelectFullPlate(plateNumber)
	if err != nil {
		fullPlate = InputNewPlateO.Plate
	}
	ent, lastans, _ := SelectFromDBEntrance(plateNumber, strconv.Itoa(InputNewPlateO.Ansid))
	if have == true {

		destination, err = SelectDestFromLocalDB(plateNumber)
		if err != nil {
			PrintlnToLogFile(`Couldn't get destination point: ` + InputNewPlateO.Plate + err.Error())
			fmt.Fprintln(w, `Couldn't get destination point`, http.StatusInternalServerError)

			return
		}
		UpdateDBinformation(plateNumber, ``, destination, lastans, InputNewPlateO.Ansid)
	} else {
		PrintlnToLogFile(`Couldn't get destination point: ` + InputNewPlateO.Plate)
		fmt.Fprintln(w, `Couldn't get destination point`, http.StatusInternalServerError)
		return
		// InsertToDBNewPlate(InputNewPlateO.Ansid, InputNewPlateO.Plate, destination, false)
		// InsertNewDestination(InputNewPlateO.Plate, destination)
	}

	err = PostTORoutingService(InputNewPlateO.Ansid, InputNewPlateO.Plate, destination, ent, plateNumber)
	if err != nil {
		PrintlnToLogFile(`Couldn't make request to Routing service: ` + err.Error())
		fmt.Println(err)

		fmt.Fprintln(w, `Couldn't make request to Routing service: `, err)
		return
	}

	// strjson := `
	// 	{
	// 		"plateNumber": "` + RoutingDataO.PlateNumber + `",
	// 		"direction": "` + RoutingDataO.Direction + `",
	// 		"destinationPoint":"` + RoutingDataO.DestinationPoint + `"
	// 	}
	// `

	if lastans == true {
		RoutingDataO.DestinationPoint = `EXIT`
	}
	RoutingDataO.PlateNumber = fullPlate
	js, err := json.Marshal(RoutingDataO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		PrintlnToLogFile(`Couldn't unmarshal json of results, ` + err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func commingPlate(w http.ResponseWriter, r *http.Request) {
	number, ok := r.URL.Query()["number"]
	if !ok || len(number[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'number' is missing")
		return
	}

	destination, ok := r.URL.Query()["destination"]
	if !ok || len(destination[0]) < 1 {
		fmt.Fprintln(w, "Url Param 'destination' is missing")
		return
	}

	cutplate := cutnumber(number[0])

	InsertToDBNewPlate(InputNewPlateO.Ansid, number[0], destination[0], false)
	InsertNewDestination(number[0], destination[0], cutplate)
}
