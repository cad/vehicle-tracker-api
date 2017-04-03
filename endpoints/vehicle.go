package endpoints

import (
	"github.com/gorilla/mux"
	valid "github.com/asaskevich/govalidator"
	"net/http"
	"encoding/json"
	"github.com/cad/vehicle-tracker-api/repository"
	"log"
)


type VehicleCreationStruct struct {

	PlateID string `json:"plate_id" valid:"required"`

}

type VehicleSyncStruct struct {

	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`

}


func GetVehicle(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	plateID := params["plateID"]

	//id, _ := strconv.Atoi(p)

	vehicle, err := repository.GetVehicleByPlateID(plateID)
	if vehicle == (repository.Vehicle{}) {
		sendErrorMessage(w, "Not found", 404)
		return
	}

	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


func GetAllVehicles(w http.ResponseWriter, req *http.Request) {
	var vehicles []repository.Vehicle
	vehicles = repository.GetAllVehicles()
	log.Println("Vehicles: ", vehicles)

	j, err := json.Marshal(vehicles)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


func CreateNewVehicle(w http.ResponseWriter, req *http.Request) {
	var vehicleStruct VehicleCreationStruct

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&vehicleStruct); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	_, err := valid.ValidateStruct(vehicleStruct)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	repository.CreateVehicle(vehicleStruct.PlateID)
	sendContentType(w, "application/json")
}


func SyncVehicle(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	plateID := params["plateID"]

	//id, _ := strconv.Atoi(p)

	var vehicleSyncStruct VehicleSyncStruct

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&vehicleSyncStruct); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	repository.SyncVehicleByPlateID(plateID, vehicleSyncStruct.Lat, vehicleSyncStruct.Lon)
	sendContentType(w, "application/json")
}



func DeleteVehicle(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	plateID, _ := params["plateID"]

	error := repository.DeleteVehicleByPlateID(plateID)

	if error!=nil {
		sendErrorMessage(w, "There is no vehicle with that ID", http.StatusNotFound)
		return
	}
	sendContentType(w, "application/json")
}
