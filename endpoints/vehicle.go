package endpoints

import (
	"github.com/gorilla/mux"
	valid "github.com/asaskevich/govalidator"
	"net/http"
	"encoding/json"
	"github.com/cad/vehicle-tracker-api/repository"
	"strings"
	"strconv"
//	"log"
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

	// VehicleTypes
	//
	// Comma seperated list of vehicle types.
	// e.g. "SCHOOL-BUS,SOLAR-CAR"
	//
	// in: query
	// required: false
	VehicleTypes string `json:"vehicle_types"`

	// VehicleGroups
	//
	// Comma seperated list of vehicle group ids.
	// e.g. "1,3"
	//
	// in: query
	// required: false
	VehicleGroups string `json:"vehicle_groups"`

}

// swagger:route GET /vehicle/filter Vehicles FilterVehicles
// Get all vehicles in the database.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehiclesResponse
func FilterVehicles(w http.ResponseWriter, req *http.Request) {
	params := FilterVehiclesParams{
		VehicleTypes: req.URL.Query().Get("vehicle_types"),
		VehicleGroups: req.URL.Query().Get("vehicle_types"),
	}
	types := strings.Split(params.VehicleTypes, ",")
	groups := strings.Split(params.VehicleGroups, ",")

	var groupIDs []uint
	index := 0
	for _, element := range groups {
		val, err := strconv.Atoi(element)
		if err != nil {
			//sendErrorMessage(w, err.Error(), http.StatusBadRequest)
			continue
		}
		groupIDs[index] = uint(val)
		index = index  + 1
	}

	var vehicles []repository.Vehicle

	vehicles = repository.FilterVehicles(types, groupIDs)

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

	// Ident represents the identity definition of the Group	       // in: body
	// required: true
	Ident struct {

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

	if err := decoder.Decode(&params.Ident); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	groupID, err := repository.CreateNewGroup(
		params.Ident.Name,
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
