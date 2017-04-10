package endpoints

import (
	"github.com/gorilla/mux"
	valid "github.com/asaskevich/govalidator"
	"net/http"
	"encoding/json"
	"github.com/cad/vehicle-tracker-api/repository"
	"log"
)

// VehicleResponse
// swagger:response
type VehicleResponse struct {
	// Vehicle
	// in: body
	Body struct {
		Vehicle repository.Vehicle `json:"vehicle,required"`
	}
}


// CoordinatePair represents a location on earth
//
// swagger:model
type CoordinatePair struct {

	// Latitude
	//
	// required: true
	Lat float64 `json:"lat"`

	// Longtitude
	//
	// required: true
	Lon float64 `json:"lon"`
}


// swagger:parameters GetVehicle
type GetVehicleParams struct {

	// PlateID is a unique identifier across the vehicles
	// in: path
	// required: true
	PlateID string `json:"plate_id"`

}

// swagger:route GET /vehicles/{plate_id} Vehicles GetVehicle
// Get a vehicle from database.
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleResponse
func GetVehicle(w http.ResponseWriter, req *http.Request) {
	params := GetVehicleParams{PlateID: mux.Vars(req)["plate_id"]}

	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if vehicle == (repository.Vehicle{}) {
		sendErrorMessage(w, "Not found", 404)
		return
	}

	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:route GET /vehicles/ Vehicles GetAllVehicles
// Get all vehicles in the database.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehiclesResponse
func GetAllVehicles(w http.ResponseWriter, req *http.Request) {
	var vehicles []repository.Vehicle
	vehicles = repository.GetAllVehicles()
	log.Println("Vehicles: ", vehicles)

	j, err := json.Marshal(vehicles)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:parameters CreateNewVehicle
type CreateNewVehicleParams struct {

	// Ident represents the identity definition of the  Vehicle
	// in: body
	// required: true
	Ident struct {

		// PlateID
		//
		// required: true
		PlateID string `json:"plate_id" valid:"required"`

	}

}

// swagger:route POST /vehicles/ Vehicles CreateNewVehicle
// Create a new vehicle record.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleResponse
func CreateNewVehicle(w http.ResponseWriter, req *http.Request) {
	var params CreateNewVehicleParams

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	repository.CreateVehicle(params.Ident.PlateID)
	sendContentType(w, "application/json")
}


// swagger:parameters SyncVehicle
type SyncVehicleParams struct {

	// PlateID is an unique identifier across vehicles
	// in: path
	// required: true
	PlateID string `json:"plate_id" validate:"required"`

	// Location represents the x,y location of the Vehicle
	// in: body
	// required: true
	Location CoordinatePair `json:"location" validate:"required"`

}

// swagger:route POST /vehicles/{plate_id}/sync Vehicles SyncVehicle
// Set position.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleResponse
func SyncVehicle(w http.ResponseWriter, req *http.Request) {
	params := SyncVehicleParams{PlateID: mux.Vars(req)["plate_id"]}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params.Location); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	repository.SyncVehicleByPlateID(params.PlateID, params.Location.Lat, params.Location.Lon)
	sendContentType(w, "application/json")
}


// swagger:parameters DeleteVehicle
type DeleteVehicleParams struct {

	// PlateID is an unique identifier across vehicles
	// in: path
	// required: true
	PlateID string `json:"plate_id" validate:"required"`

}

// swagger:route DELETE /vehicles/{plate_id} Vehicles DeleteVehicle
// Delete a vehicle.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleResponse
func DeleteVehicle(w http.ResponseWriter, req *http.Request) {
	params := DeleteVehicleParams{PlateID: mux.Vars(req)["PlateID"]}

	error := repository.DeleteVehicleByPlateID(params.PlateID)

	if error!=nil {
		sendErrorMessage(w, "There is no vehicle with that ID", http.StatusNotFound)
		return
	}
	sendContentType(w, "application/json")
}
