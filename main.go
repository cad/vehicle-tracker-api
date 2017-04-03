package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"github.com/cad/vehicle-tracker-api/endpoints"
	"github.com/cad/vehicle-tracker-api/repository"
	"github.com/cad/vehicle-tracker-api/config"
	"fmt"
	"os"
	"flag"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args)>0 {
		switch args[0] {
		case "serve":
			executeServer()
		default:
			fmt.Println("Invalid command")
			os.Exit(1)
		}
	} else {
		fmt.Println("Invalid command, please retry with the following format $ main serve")
		os.Exit(1)
	}

}


func GetServer() http.Handler {
	router := mux.NewRouter()

	// Vehicle
	router.HandleFunc("/vehicles/", endpoints.GetAllVehicles).Methods("GET")
	router.HandleFunc("/vehicles/", endpoints.CreateNewVehicle).Methods("POST")
	router.HandleFunc("/vehicles/{plateID}", endpoints.GetVehicle).Methods("GET")
	router.HandleFunc("/vehicles/{plateID}", endpoints.DeleteVehicle).Methods("DELETE")
	router.HandleFunc("/vehicles/{plateID}/sync", endpoints.SyncVehicle).Methods("POST")


	return router

}


func executeServer() {
	if err := config.LoadConfigFile("./config.json"); err != nil {
		fmt.Printf("Error: %s loading configuration file: %s\n", "./config.json", err)
		os.Exit(1)
	}

	repository.ConnectDB("sqlite3", "data/devel.db")

	router := GetServer()
	fmt.Println("API server is listening on port :5000")
	log.Fatal(http.ListenAndServe(":5000", router))

	defer repository.CloseDB()
}
