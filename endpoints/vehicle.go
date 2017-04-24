package endpoints

import (
	"github.com/gorilla/mux"
	valid "github.com/asaskevich/govalidator"
	"net/http"
	"encoding/json"
	"github.com/cad/vehicle-tracker-api/repository"
	//"strings"
	"strconv"
	//	"log"
//	"fmt"
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

// VehicleGroupResponse
// swagger:response
type VehicleGroupsResponse struct {
	// VehicleGroup
	// in: body
	Body struct {
		Groups []repository.Group `json:"group,required"`
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

// swagger:route GET /vehicle/{plate_id} Vehicles GetVehicle
// Get a vehicle from database.
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleResponse
func GetVehicle(w http.ResponseWriter, req *http.Request) {
	params := GetVehicleParams{PlateID: mux.Vars(req)["plate_id"]}

	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if vehicle.ID == 0 {
		sendErrorMessage(w, "Not found", 404)
		return
	}

	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:route GET /vehicle/ Vehicles GetAllVehicles
// Get all vehicles in the database.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehiclesResponse
func GetAllVehicles(w http.ResponseWriter, req *http.Request) {
	var vehicles []repository.Vehicle
	vehicles = repository.GetAllVehicles()
//	log.Println("Vehicles: ", vehicles)

	j, err := json.Marshal(vehicles)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:parameters FilterVehicles
type FilterVehiclesParams struct {

	// VehicleType
	//
	// VehicleType to be filtered.
	// e.g: "SCHOOL-BUS"
	//
	//
	// in: query
	// required: false
	VehicleType string `json:"vehicle_type"`

	// VehicleGroup
	//
	// VehicleGroup id to be filtered.
	// e.g: 3
	//
	// in: query
	// required: false
	VehicleGroupID int `json:"vehicle_group_id"`

}

// swagger:route GET /vehicle/filter Vehicles FilterVehicles
// Filter vehicles in the database.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehiclesResponse
func FilterVehicles(w http.ResponseWriter, req *http.Request) {
	groupID, err := strconv.Atoi(req.URL.Query().Get("vehicle_group_id"))
	if err != nil {
		sendErrorMessage(w, "vehicle_group_id should be int", http.StatusBadRequest)
		return

	}
	params := FilterVehiclesParams{
		VehicleType: req.URL.Query().Get("vehicle_type"),
		VehicleGroupID: groupID,
	}
	var vehicles []repository.Vehicle
	vehicles = repository.FilterVehicles(params.VehicleType, uint(params.VehicleGroupID))

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

		// AgentID
		//
		// required: false
		AgentUUID string `json:"agent_uuid"`

		// Groups
		//
		// required: false
		Groups []int `json:"groups"`

		// Type
		//
		// required: true
		Type string `json:"type" valid:"required"`

	}

}

// swagger:route POST /vehicle/ Vehicles CreateNewVehicle
// Create a new vehicle record.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleResponse
func CreateNewVehicle(w http.ResponseWriter, req *http.Request) {
	var params CreateNewVehicleParams

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params.Ident); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.CreateVehicle(
		params.Ident.PlateID,
		params.Ident.AgentUUID,
		params.Ident.Groups,
		params.Ident.Type,
	)

	if err != nil  {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	vehicle, err := repository.GetVehicleByPlateID(params.Ident.PlateID)
	if err != nil  {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}


	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:parameters VehicleSetAgent
type VehicleSetAgentParams struct {
	// PlateID is a unique identifier across the vehicles
	// in: path
	// required: true
	PlateID string `json:"plate_id"`

	// Agent represents an agent
	// in: body
	// required: true
	Agent struct {

		// UUID
		//
		// required: true
		UUID string `json:"uuid" valid:"required"`

	}

}

// swagger:route POST /vehicle/{plate_id}/agent Vehicles VehicleSetAgent
// Set vehicle agent.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleResponse
func VehicleSetAgent(w http.ResponseWriter, req *http.Request) {
	var params VehicleSetAgentParams

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params.Agent); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	params.PlateID = mux.Vars(req)["plate_id"]
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.VehicleSetAgent(params.PlateID, params.Agent.UUID)
	if err != nil  {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil  {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}


	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:parameters DeleteVehicle
type DeleteVehicleParams struct {

	// PlateID is an unique identifier across vehicles
	// in: path
	// required: true
	PlateID string `json:"plate_id" validate:"required"`

}

// swagger:route DELETE /vehicle/{plate_id} Vehicles DeleteVehicle
// Delete a vehicle.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleResponse
func DeleteVehicle(w http.ResponseWriter, req *http.Request) {
	params := DeleteVehicleParams{PlateID: mux.Vars(req)["PlateID"]}

	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil {
		sendErrorMessage(w, "There is no vehicle with that ID", http.StatusNotFound)
		return
	}

	error := repository.DeleteVehicleByPlateID(params.PlateID)
	if error!=nil {
		sendErrorMessage(w, "There is no vehicle with that ID", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}



// swagger:parameters CreateNewGroup
type CreateNewGroupParams struct {

	// Group
	// in: body
	// required: true
	Group struct {

		// Name
		//
		// required: true
		Name string `json:"name" valid:"required"`
	}

}

// swagger:route POST /vehicle/group/ Vehicles CreateNewGroup
// Create a new vehicle group.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleGroupResponse
func CreateNewGroup(w http.ResponseWriter, req *http.Request) {
	var params CreateNewGroupParams

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params.Group); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupID, err := repository.CreateNewGroup(
		params.Group.Name,
	)

	if err != nil  {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	group, err := repository.GetGroupByID(groupID)

	j, err := json.Marshal(group)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}


// swagger:route GET /vehicle/group/ Vehicles GetAllGroups
// Get all vehicle groups in the database.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleGroupsResponse
func GetAllGroups(w http.ResponseWriter, req *http.Request) {
	var groups []repository.Group
	groups = repository.GetAllGroups()
	j, err := json.Marshal(groups)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}
