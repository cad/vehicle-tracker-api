package endpoints

import (
	"github.com/gorilla/mux"
	//	valid "github.com/asaskevich/govalidator"
	"encoding/json"
	"net/http"

	"github.com/cad/vehicle-tracker-api/repository"
	//	"log"
	//	"fmt"
)

type GPSData struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
	TS  string `json:"ts"`
}

// swagger:parameters FilterAgents
type FilterAgentsParams struct {

	// AgentState
	//
	// AgentState to be filtered.
	// "ASSIGNED" or "UNASSIGNED"
	//
	//
	// in: query
	// required: false
	// enum: ASSIGNED,UNASSIGNED
	AgentState string `json:"state"`
}

// swagger:route GET /agent/ Agents FilterAgents
// Filter vehicles.
//
//
//   Responses:
//     default: ErrorMsg
//     200: AgentSuccessAgentsResponse
func FilterAgents(w http.ResponseWriter, req *http.Request) {
	var agents = []repository.Agent{}
	var state string
	if len(req.URL.Query()["state"]) > 0 {
		state = req.URL.Query()["state"][0]
	}
	params := FilterAgentsParams{AgentState: state}

	agents = repository.FilterAgents(params.AgentState)

	j, err := json.Marshal(agents)
	checkErr(w, err)
	sendContentType(w, "application/json")
	w.Write(j)
}

// swagger:parameters SyncAgent
type SyncAgentParams struct {

	// UUID is an unique identifier across agents
	// in: path
	// required: true
	UUID string `json:"uuid" validate:"required"`

	// Data represents the x,y location of the agent at ts time.
	// in: body
	// required: true
	Data GPSData
}

// swagger:route POST /agent/{uuid}/sync Agents SyncAgent
// Send GPS data from agent.
//
//
//   Responses:
//     default: ErrorMsg
//     200: AgentEmptyResponse
func SyncAgent(w http.ResponseWriter, req *http.Request) {
	params := SyncAgentParams{UUID: mux.Vars(req)["uuid"]}
	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params.Data); err != nil {
		sendErrorMessage(w, "Error decoding the input", http.StatusBadRequest)
		return
	}

	err := repository.SyncAgentByUUID(
		params.UUID,
		params.Data.Lat,
		params.Data.Lon,
		params.Data.TS,
	)
	if err != nil {
		sendErrorMessage(w, "Agent Sync Error", http.StatusBadRequest)
		return
	}

	sendContentType(w, "application/json")
}
