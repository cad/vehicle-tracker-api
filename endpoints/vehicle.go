package endpoints

import (
	"encoding/json"
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/cad/vehicle-tracker-api/event"
	"github.com/cad/vehicle-tracker-api/repository"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	//"strings"
	"log"
	"strconv"
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

	// AgentState
	//
	// AgentState to be filtered.
	// "ASSIGNED" or "UNASSIGNED"
	//
	//
	// in: query
	// required: false
	// enum: ASSIGNED,UNASSIGNED
	AgentState string `json:"agent_state"`
}

// swagger:route GET /vehicle/filter Vehicles FilterVehicles
// Filter vehicles in the database.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehiclesResponse
func FilterVehicles(w http.ResponseWriter, req *http.Request) {
	var groupID int

	groupIDStr := req.URL.Query().Get("vehicle_group_id")
	if groupIDStr != "" {
		var err error
		groupID, err = strconv.Atoi(groupIDStr)
		if err != nil {
			sendErrorMessage(w, "vehicle_group_id should be int", http.StatusBadRequest)
			return
		}
	}

	params := FilterVehiclesParams{
		VehicleType:    req.URL.Query().Get("vehicle_type"),
		VehicleGroupID: groupID,
		AgentState:     req.URL.Query().Get("agent_state"),
	}
	var vehicles []repository.Vehicle
	vehicles = repository.FilterVehicles(params.VehicleType, uint(params.VehicleGroupID), params.AgentState)

	j, err := json.Marshal(vehicles)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// swagger:parameters FilterVehiclesWS
type FilterVehiclesWSParams struct {

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

// swagger:route GET /ws/vehicle/filter WebSocket FilterVehiclesWS
// WebSocket Endpoint for filter vehicles.
//
// e.g. wss://api.vehicles.neu.edu.tr/ws/vehicle/filter?vehicle_type=SCHOOL-BUS&vehicle_group_id=2
//
//   Responses:
//     200: VehicleSuccessVehicleResponse
//
func FilterVehiclesWS(w http.ResponseWriter, req *http.Request) {
	groupID, err := strconv.Atoi(req.URL.Query().Get("vehicle_group_id"))
	if err != nil {
		sendErrorMessage(w, "vehicle_group_id should be int", http.StatusBadRequest)
		return

	}
	params := FilterVehiclesParams{
		VehicleType:    req.URL.Query().Get("vehicle_type"),
		VehicleGroupID: groupID,
	}

	c, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	newAgentEvent := event.MakeKind(repository.NEW_AGENT)
	handler := func(e *event.Event) {
		agent, ok := e.Payload.(repository.Agent)
		if !ok {
			log.Println("Cannot assert type Agent. Ignoring.")
			return
		}

		vehicle, err := repository.GetVehicleByAgentUUID(agent.UUID)
		if err != nil {
			// Ignore update
			log.Printf("(%s)", err.Error())
			log.Println("[WS-EXPORT] Vehicle not found for the Agent streamed from channel. Ignoring.", "Agent.UUID", agent.UUID)
			return
		}
		if params.VehicleType != "" && vehicle.Type != params.VehicleType {
			// Ignore update
			return
		}

		found := false
		for _, group := range vehicle.Groups {
			if group.ID == uint(params.VehicleGroupID) {
				found = true
			}
		}
		if found != true {
			// Ignore update
			return
		}

		err = c.WriteJSON(vehicle)
		if err != nil {
			log.Println("[WS-EXPORT] Can't write to WS Connection!. Ignoring.")
			return
		}
	}
	newAgentEvent.Register(&handler)
	defer newAgentEvent.UnRegister(&handler)

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("Client disconnected")
			break
		}
	}

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
//   Security:
//       Bearer:
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

	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	vehicle, err := repository.GetVehicleByPlateID(params.Ident.PlateID)
	if err != nil {
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
//   Security:
//       Bearer:
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
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(vehicle)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters VehicleUnsetAgent
type VehicleUnsetAgentParams struct {
	// PlateID is a unique identifier across the vehicles
	// in: path
	// required: true
	PlateID string `json:"plate_id"`
}

// swagger:route DELETE /vehicle/{plate_id}/agent Vehicles VehicleUnsetAgent
// Unset vehicle agent.
//
//   Security:
//       Bearer:
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleResponse
func VehicleUnsetAgent(w http.ResponseWriter, req *http.Request) {
	var params VehicleUnsetAgentParams

	params.PlateID = mux.Vars(req)["plate_id"]
	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.VehicleUnsetAgent(params.PlateID)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil {
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
//   Security:
//       Bearer:
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleResponse
func DeleteVehicle(w http.ResponseWriter, req *http.Request) {
	params := DeleteVehicleParams{PlateID: mux.Vars(req)["plate_id"]}

	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil {
		sendErrorMessage(w, "There is no vehicle with that ID", http.StatusNotFound)
		return
	}

	error := repository.DeleteVehicleByPlateID(params.PlateID)
	if error != nil {
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
//   Security:
//       Bearer:
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

	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	group, err := repository.GetGroupByID(groupID)

	j, err := json.Marshal(group)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters DeleteVehicleGroup
type DeleteGroupParams struct {

	// GroupID
	// in: path
	// required: true
	ID string `json:"group_id"`
}

// swagger:route DELETE /vehicle/group/{group_id} Vehicles DeleteVehicleGroup
// Delete a group definition.
//
//   Security:
//       Bearer:
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleGroupResponse
func DeleteGroup(w http.ResponseWriter, req *http.Request) {
	params := DeleteGroupParams{ID: mux.Vars(req)["group_id"]}

	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}
	groupID, err := strconv.Atoi(params.ID)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	group, err := repository.GetGroupByID(uint(groupID))
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = repository.DeleteGroup(uint(groupID))
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	j, err := json.Marshal(group)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters SetVehicleGroups
type SetVehicleGroupsParams struct {

	// PlateID
	// in: path
	// required: true
	PlateID string `json:"plate_id"`

	// Groups
	// in: body
	// required: true
	Ident struct {
		Groups []int `json:"groups"`
	}
}

// swagger:route PUT /vehicle/{plate_id}/groups Vehicles SetVehicleGroups
// Set vehicle's groups.
//
//   Security:
//       Bearer:
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleGroupsResponse
func SetVehicleGroups(w http.ResponseWriter, req *http.Request) {
	params := SetVehicleGroupsParams{PlateID: mux.Vars(req)["plate_id"]}

	_, err := valid.ValidateStruct(params)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&params.Ident); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}

	err = repository.SetVehicleGroups(params.PlateID, params.Ident.Groups)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil {
		sendErrorMessage(w, err.Error(), http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(vehicle.Groups)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:route GET /vehicle/type/ Vehicles GetAllTypes
// Get possible vehicle types defined in the system.
//
//
//   Responses:
//     default: ErrorMsg
//     200: VehicleSuccessVehicleTypesResponse
func GetAllTypes(w http.ResponseWriter, req *http.Request) {
	var types []string
	types = repository.GetAllTypes()
	j, err := json.Marshal(types)
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
