package endpoints

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "github.com/cad/vehicle-tracker-api/statik"
	"github.com/cad/statik/fs"
	"io/ioutil"
)

func GetRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/vehicles/", GetAllVehicles).Methods("GET")

	router.HandleFunc("/vehicles/", CreateNewVehicle).Methods("POST")
	router.HandleFunc("/vehicles/{plate_id}", GetVehicle).Methods("GET")
	router.HandleFunc("/vehicles/{plate_id}", DeleteVehicle).Methods("DELETE")
	router.HandleFunc("/vehicles/{plate_id}/sync", SyncVehicle).Methods("POST")
	dataFS, err := fs.New("/")
	if err != nil {
		log.Fatalf(err.Error())
	}

	router.HandleFunc("/spec", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, "data/swagger.json")
		specFile, err := dataFS.Open("/swagger.json")
		if err != nil {
			log.Fatal(err)
		}
		contents, err := ioutil.ReadAll(specFile)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(contents)
	})
	router.HandleFunc("/docs/", func(w http.ResponseWriter, r *http.Request) {
		specFile, err := dataFS.Open("/static/index.html")
		if err != nil {
			log.Fatal(err)
		}
		contents, err := ioutil.ReadAll(specFile)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(contents)
	})


	statikFS, err := fs.New("/static/")
	if err != nil {
		log.Fatalf(err.Error())
	}
	//router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(statikFS)))

	return router
}
