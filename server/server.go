package server


import (
	"log"
	"net/http"
	"github.com/cad/vehicle-tracker-api/endpoints"
	"github.com/cad/vehicle-tracker-api/config"
	"github.com/cad/vehicle-tracker-api/repository"
	"fmt"
	"os"
	"github.com/gorilla/handlers"
)

func getRouter() http.Handler {
	return endpoints.GetRouter()
}


func GetServer() http.Handler {
	return getRouter()
}


func ExecuteServer() {
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
