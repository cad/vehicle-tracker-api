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

func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

func GetRouter() http.Handler {
	router := mux.NewRouter()

	// Apply CORS to all preflight (OPTIONS) request.
	router.Methods("OPTIONS").HandlerFunc( func(w http.ResponseWriter, r *http.Request){
		doCORS(w, r)
	})

	// Users
	router.HandleFunc("/user/", use(GetAllUsers, TokenAuthMiddleware, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/user/", use(CreateNewUser, TokenAuthMiddleware, CORSMiddleware)).Methods("POST")
	router.HandleFunc("/user/{uuid}", use(GetUser, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/user/{uuid}", use(DeleteUser, TokenAuthMiddleware, CORSMiddleware)).Methods("DELETE")

	// Auth
	router.HandleFunc("/auth/", use(CheckAuth, TokenAuthMiddleware, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/auth/", use(Authorize, CORSMiddleware)).Methods("POST")

	// Agents
	router.HandleFunc("/agent/", use(GetAllAgents, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/agent/{uuid}/sync", use(SyncAgent, CORSMiddleware)).Methods("POST")
	router.HandleFunc("/agents/{uuid}/sync", use(SyncAgent, CORSMiddleware)).Methods("POST")  // NOTE(cad): this line added for backwards compatibility

	// Vehicles
	router.HandleFunc("/vehicle/", use(GetAllVehicles, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/vehicle/filter", use(FilterVehicles, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/vehicle/", use(CreateNewVehicle, TokenAuthMiddleware, CORSMiddleware)).Methods("POST")
	router.HandleFunc("/vehicle/group/", use(GetAllGroups, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/vehicle/group/", use(CreateNewGroup, TokenAuthMiddleware, CORSMiddleware)).Methods("POST")
	router.HandleFunc("/vehicle/{plate_id}/agent", use(VehicleSetAgent, TokenAuthMiddleware, CORSMiddleware)).Methods("POST")
	router.HandleFunc("/vehicle/{plate_id}", use(GetVehicle, CORSMiddleware)).Methods("GET")
	router.HandleFunc("/vehicle/{plate_id}", use(DeleteVehicle, TokenAuthMiddleware, CORSMiddleware)).Methods("DELETE")

	// WebSocket
	router.HandleFunc("/ws/vehicle/filter", use(FilterVehiclesWS, CORSMiddleware)).Methods("GET")

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
