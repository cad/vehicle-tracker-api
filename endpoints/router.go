package endpoints

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "github.com/cad/vehicle-tracker-api/statik"
	"github.com/cad/vehicle-tracker-api/config"
	"github.com/cad/statik/fs"
	"io/ioutil"
	"encoding/json"
)

func GetRouter() http.Handler {
	router := mux.NewRouter()

	// Agents
	router.HandleFunc("/agent/", GetAllAgents).Methods("GET")
	router.HandleFunc("/agent/{uuid}/sync", SyncAgent).Methods("POST")

	// Vehicles
	router.HandleFunc("/vehicle/", GetAllVehicles).Methods("GET")
	router.HandleFunc("/vehicle/filter", FilterVehicles).Methods("GET")
	router.HandleFunc("/vehicle/", CreateNewVehicle).Methods("POST")
	router.HandleFunc("/vehicle/group/", GetAllGroups).Methods("GET")
	router.HandleFunc("/vehicle/group/", CreateNewGroup).Methods("POST")
	router.HandleFunc("/vehicle/{plate_id}/agent", VehicleSetAgent).Methods("POST")
	router.HandleFunc("/vehicle/{plate_id}", GetVehicle).Methods("GET")
	router.HandleFunc("/vehicle/{plate_id}", DeleteVehicle).Methods("DELETE")


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
		var data map[string]interface{}
		err = json.Unmarshal(contents, &data)
		if err != nil {
			log.Fatal(err)
		}

		// Override
		data["host"] = r.Host

		s, ok := data["info"].(map[string]interface{})
		if !ok {
			log.Fatal("unknown type", s, ok)
		}
		s["version"] = config.VERSION
		data["info"] = s

		encoded_data, err := json.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		sendContentType(w, "application/json")
		w.Write(encoded_data)
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
