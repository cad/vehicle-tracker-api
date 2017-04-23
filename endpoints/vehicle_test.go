package endpoints

import (
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"testing"
	"bytes"
	"github.com/cad/vehicle-tracker-api/repository"
)


func TestGetAllVehiclesEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	_ = repository.CreateVehicle(
		"test",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)

	// Execute
	req, _ := http.NewRequest("GET", "/vehicle/", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	var vehicles []repository.Vehicle
	err := json.Unmarshal([]byte(res.Body.String()), &vehicles)
	if err != nil {
		t.Error(errorMsg("Vehicles", "Marshallable", "NotUnmarshallable"))
		return
	}

	if vehicles[0].PlateID != "test" {
		t.Error(errorMsg("PlateID", "test", vehicles[0].PlateID))
		return
	}

	if vehicles[0].Groups[0].Name != "string" {
		t.Error(errorMsg("Groups[0].Name", "test", vehicles[0].Groups[0].Name))
		return
	}

	if vehicles[0].Type != "SCHOOL-BUS" {
		t.Error(errorMsg("Type", "SCHOOL-BUS", vehicles[0].Type))
		return
	}

	if vehicles[0].Agent.UUID != agent.UUID {
		t.Error(errorMsg("Agent.UUID", fmt.Sprintf("%d", agent.UUID), fmt.Sprintf("%d", vehicles[0].Agent.UUID)))
		return
	}



}


func TestCreateNewVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	type Ident struct {

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

	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	vehicle := Ident{
		PlateID: "test",
		AgentUUID: agent.UUID,
		Groups: []int{int(groupID)},
		Type: "SCHOOL-BUS",
	}
	vehicle_json, err := json.Marshal(&vehicle)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Marshallable", "NotMarshallable"))
	}
	body := bytes.NewBuffer(vehicle_json)

	// Execute
	req, _ := http.NewRequest("POST", "/vehicle/", body)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d",res.Code)))
		return
	}
}


func TestGetVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	_ = repository.CreateVehicle(
		"test",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)


	// Execute
	req, _ := http.NewRequest("GET", "/vehicle/test", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d",res.Code)))
		return
	}

	var vehicle repository.Vehicle
	err := json.Unmarshal([]byte(res.Body.String()), &vehicle)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Unmarshallable", "NotUnmarshallable"))
		return
	}


	if vehicle.PlateID != "test" {
		t.Error(errorMsg("PlateID", "test", vehicle.PlateID))
		return
	}

	if vehicle.Groups[0].Name != "string" {
		t.Error(errorMsg("Groups[0].Name", "test", vehicle.Groups[0].Name))
		return
	}

	if vehicle.Type != "SCHOOL-BUS" {
		t.Error(errorMsg("Type", "SCHOOL-BUS", vehicle.Type))
		return
	}

	if vehicle.Agent.UUID != agent.UUID {
		t.Error(errorMsg("Agent.UUID", fmt.Sprintf("%d", agent.UUID), fmt.Sprintf("%d", vehicle.Agent.UUID)))
		return
	}



}


func TestDeleteVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	_ = repository.CreateVehicle(
		"test",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)


	// Execute
	req, _ := http.NewRequest("DELETE", "/vehicle/test", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	vehicles := repository.GetAllVehicles()
	if len(vehicles) > 0 {
		t.Error(errorMsg("len(vehicles)", "0", fmt.Sprintf("%d", len(vehicles))))
		return
	}
}
// Test Set Agent
// Test Filter
// Test Group Create
// Test Group GetAll
