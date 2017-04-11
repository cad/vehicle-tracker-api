//go:generate swagger generate spec -o data/swagger.json
package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"github.com/cad/vehicle-tracker-api/endpoints"
	"github.com/cad/vehicle-tracker-api/repository"
	"github.com/cad/vehicle-tracker-api/config"
	"fmt"
	"os"
	"flag"
	//"github.com/go-swagger/go-swagger/cmd/swagger/commands/generate"
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


// GetServer creates a mux.Router for Vehicle Tracker API
func GetServer() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/vehicles/", endpoints.GetAllVehicles).Methods("GET")

	router.HandleFunc("/vehicles/", endpoints.CreateNewVehicle).Methods("POST")
	router.HandleFunc("/vehicles/{plate_id}", endpoints.GetVehicle).Methods("GET")
	router.HandleFunc("/vehicles/{plate_id}", endpoints.DeleteVehicle).Methods("DELETE")
	router.HandleFunc("/vehicles/{plate_id}/sync", endpoints.SyncVehicle).Methods("POST")

	router.HandleFunc("/spec", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "data/swagger.json")
	})

	router.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return router

}


func executeServer() {
	if err := config.LoadConfigFile("./config.json"); err != nil {
		fmt.Printf("Error: %s loading configuration file: %s\n", "./config.json", err)
		os.Exit(1)
	}

	//specFile := generate.NewSpecFile("github.com/cad/vehicle-tracker-api", "data/swagger.json")
	//specFile.Execute([]string{})
	repository.ConnectDB(config.C.DB.Type , config.C.DB.URL)

	router := GetServer()
	router = handlers.LoggingHandler(os.Stdout, router)
	fmt.Println("API server version", config.VERSION, "is listening on port", config.C.Server.Port)
	log.Fatal(http.ListenAndServe(config.C.Server.Port, router))

	defer repository.CloseDB()
}
