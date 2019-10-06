package main

//AppConfig structure of application configuration
type AppConfig struct {
	Port           int
	NavisURL       string
	RoutingServURL string
}

//AppConfigO object model of AppConfig
var AppConfigO AppConfig

//DBconfig describes the config params
type DBconfig struct {
	DBName string `json:"dbName"`
	DBPass string `json:"dbPass"`
	DBUser string `json:"dbUSer"`
	DBHost string `json:"dbHost"`
	DBPort int    `json:"dbPort"`
}

//DBConfigInterface interface of config params
var DBConfigInterface DBconfig

//InputNewPlate - handling information about new plates
type InputNewPlate struct {
	Ansid       int
	Plate       string `json:"plateNumber"`
	Destination string
}

//InputNewPlateO - object of InputNewPlate
var InputNewPlateO InputNewPlate

type RoutingData struct {
	PlateNumber      string
	Direction        string
	DestinationPoint string
	IsLastAns        bool
}

var RoutingDataO RoutingData
