package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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
		t.Error(errorMsg("Agent.UUID", fmt.Sprintf("%s", agent.UUID), fmt.Sprintf("%s", vehicles[0].Agent.UUID)))
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
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()

	//agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	vehicle := Ident{
		PlateID: "test",
		Groups:  []int{int(groupID)},
		Type:    "SCHOOL-BUS",
	}
	vehicleJSON, err := json.Marshal(&vehicle)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Marshallable", "NotMarshallable"))
	}
	body := bytes.NewBuffer(vehicleJSON)

	// Execute
	req, _ := http.NewRequest("POST", "/vehicle/", body)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	agent, _ := repository.CreateNewAgent("string")

	vehicle = Ident{
		PlateID:   "test2",
		AgentUUID: agent.UUID,
		Groups:    []int{int(groupID)},
		Type:      "SCHOOL-BUS",
	}
	vehicleJSON, err = json.Marshal(&vehicle)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Marshallable", "NotMarshallable"))
	}
	body = bytes.NewBuffer(vehicleJSON)

	// Execute
	req, _ = http.NewRequest("POST", "/vehicle/", body)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res = httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	// Ensure vehicle creation without providing agent; but
	// when an agent is already present in the system.
	vehicle2 := Ident{
		PlateID: "test22",
		Type:    "SCHOOL-BUS",
	}
	vehicleJSON, err = json.Marshal(&vehicle2)
	if err != nil {
		t.Error(errorMsg("Vehicle2", "Marshallable", "NotMarshallable"))
	}
	body = bytes.NewBuffer(vehicleJSON)

	// Execute
	req, _ = http.NewRequest("POST", "/vehicle/", body)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res = httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
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
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
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
		t.Error(errorMsg("Agent.UUID", fmt.Sprintf("%s", agent.UUID), fmt.Sprintf("%s", vehicle.Agent.UUID)))
		return
	}

}

func TestDeleteVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()
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
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

// Test Delete Vehicle with empty input
// Ensure BUG #13 does not regress.
func TestDeleteVehicleEndpointEmptyInput(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()
	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	_ = repository.CreateVehicle(
		"test",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)

	// Execute
	req, _ := http.NewRequest("DELETE", "/vehicle/+", nil)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 404 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	vehicles := repository.GetAllVehicles()
	if len(vehicles) <= 0 {
		t.Error(errorMsg("len(vehicles)", "> 0", fmt.Sprintf("%d", len(vehicles))))
		return
	}
}

// Test Set Agent
func TestSetVehicleAgentEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()

	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	_ = repository.CreateVehicle(
		"test",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)
	newAgent, _ := repository.CreateNewAgent("another-string")
	params := VehicleSetAgentParams{
		PlateID: "test",
		Agent: struct {
			UUID string `json:"uuid" valid:"required"`
		}{
			UUID: newAgent.UUID,
		},
	}
	vehicleJSON, err := json.Marshal(&params.Agent)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Marshallable", "NotMarshallable"))
	}
	body := bytes.NewBuffer(vehicleJSON)

	// Execute
	req, _ := http.NewRequest("POST", fmt.Sprintf("/vehicle/%s/agent", params.PlateID), body)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	vehicle, err := repository.GetVehicleByPlateID(params.PlateID)
	if err != nil {
		t.Error(errorMsg("Vehicle", "ToBeAbleToGet", "CannotGet"))
	}

	if vehicle.Agent.UUID != params.Agent.UUID {
		t.Error(errorMsg("vehicle.Agent.UUID", params.Agent.UUID, vehicle.Agent.UUID))
	}
}

// Test Unset Agent
func TestUnsetVehicleAgentEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()

	agent, _ := repository.CreateNewAgent("string")
	groupID, _ := repository.CreateNewGroup("string")
	err := repository.CreateVehicle(
		"test",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)

	if err != nil {
		t.Fatalf("can not create vehicle: %v", err)
	}

	// Execute
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/vehicle/%s/agent", "test"), nil)
	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	if err != nil {
		t.Fatalf("can not create vehicle: %v", err)
	}

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	vehicle, err := repository.GetVehicleByPlateID("test")
	if err != nil {
		t.Error(errorMsg("Vehicle", "ToBeAbleToGet", "CannotGet"))
	}

	if vehicle.Agent != nil {
		fmt.Printf("%+v\n", vehicle.Agent)
		t.Error(errorMsg("vehicle.Agent", "nil", "not-nil"))
	}
}

// Test Filter
func TestFilterVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	agent, _ := repository.CreateNewAgent("string")
	newAgent, _ := repository.CreateNewAgent("another-string")
	groupID, _ := repository.CreateNewGroup("string")
	newGroupID, _ := repository.CreateNewGroup("stringg")
	_ = repository.CreateVehicle(
		"test1",
		agent.UUID,
		[]int{int(groupID)},
		"SCHOOL-BUS",
	)
	_ = repository.CreateVehicle(
		"test2",
		newAgent.UUID,
		[]int{int(newGroupID)},
		"SCHOOL-BUS",
	)

	// Execute
	req, _ := http.NewRequest("GET", fmt.Sprintf("/vehicle/filter?vehicle_group_id=%d", newGroupID), nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	var vehicles []repository.Vehicle
	err := json.Unmarshal([]byte(res.Body.String()), &vehicles)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Unmarshallable", "NotUnmarshallable"))
		return
	}
	if vehicles[0].Groups[0].ID != newGroupID {
		t.Error(errorMsg("Vehicle", fmt.Sprintf("%d", newGroupID), fmt.Sprintf("%d", vehicles[0].Groups[0].ID)))
		return

	}
}

// Test Group Create
func TestCreateVehicleGroupEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	user, _ := repository.CreateNewUser("test@test.com", "1234")
	token, _ := user.RenewToken()

	params := CreateNewGroupParams{
		Group: struct {
			Name string `json:"name" valid:"required"`
		}{Name: "string"},
	}
	_, _ = repository.CreateNewAgent("string")
	payloadJSON, err := json.Marshal(&params.Group)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Marshallable", "NotMarshallable"))
	}
	body := bytes.NewBuffer(payloadJSON)

	// Execute
	req, _ := http.NewRequest("POST", "/vehicle/group/", body)

	// Authenticate
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}
	var group repository.Group
	err = json.Unmarshal([]byte(res.Body.String()), &group)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	createdGroup, _ := repository.GetGroupByName("string")

	if createdGroup.ID != group.ID {
		t.Error(errorMsg("createdGroup.ID", fmt.Sprintf("%d", group.ID), fmt.Sprintf("%d", createdGroup.ID)))
		return
	}
}

func TestGetAllVehicleGroupsEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	repository.CreateNewGroup("string")

	// Execute
	req, _ := http.NewRequest("GET", "/vehicle/group/", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}
	var groups []repository.Group
	err := json.Unmarshal([]byte(res.Body.String()), &groups)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	createdGroups := repository.GetAllGroups()

	if len(createdGroups) != len(groups) {
		t.Error(errorMsg("len(createdGroups)", fmt.Sprintf("%d", len(groups)), fmt.Sprintf("%d", len(createdGroups))))
		return
	}

	if createdGroups[0].ID != groups[0].ID {
		t.Error(errorMsg("createdGroups[0].ID", fmt.Sprintf("%d", groups[0].ID), fmt.Sprintf("%d", createdGroups[0].ID)))
		return
	}

	if createdGroups[0].Name != groups[0].Name {
		t.Error(errorMsg("createdGroups[0].Name", fmt.Sprintf("%s", groups[0].Name), fmt.Sprintf("%s", createdGroups[0].Name)))
		return
	}

}

func TestGetAllVehicleTypesEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare

	// Execute
	req, _ := http.NewRequest("GET", "/vehicle/type/", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}
	var types []string
	err := json.Unmarshal([]byte(res.Body.String()), &types)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	createdTypes := repository.GetAllTypes()

	if len(createdTypes) != len(types) {
		t.Error(errorMsg("len(createdTypes)", fmt.Sprintf("%d", len(types)), fmt.Sprintf("%d", len(createdTypes))))
		return
	}

	if createdTypes[0] != types[0] {
		t.Error(errorMsg("createdTypes[0]", fmt.Sprintf("%s", types[0]), fmt.Sprintf("%s", createdTypes[0])))
		return
	}
}
