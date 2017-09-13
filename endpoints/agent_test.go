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

func errorMsg(what, shouldBe, was string) string {
	//debug.PrintStack()
	return fmt.Sprintf("expected %s \"%s\" but got \"%s\"", what, shouldBe, was)
}

func TestGetAllAgentsEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	_, _ = repository.CreateNewAgent("test")

	// Execute
	req, _ := http.NewRequest("GET", "/agent/", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	var agents []repository.Agent
	err := json.Unmarshal([]byte(res.Body.String()), &agents)
	if err != nil {
		t.Error(errorMsg("Agents", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	if agents[0].UUID != "test" {
		t.Error(errorMsg("UUID", "test", agents[0].UUID))
		return
	}
}

func TestFilterAgentsEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	_, _ = repository.CreateNewAgent("test")
	_, _ = repository.CreateNewAgent("assigned")
	_ = repository.CreateVehicle("testvehicle", "assigned", []int{}, "SCHOOL-BUS")

	// Execute
	req, _ := http.NewRequest("GET", "/agent/?state=ASSIGNED", nil)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	var agents []repository.Agent
	err := json.Unmarshal([]byte(res.Body.String()), &agents)
	if err != nil {
		t.Error(errorMsg("Agents", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	if agents[0].UUID != "assigned" {
		t.Error(errorMsg("UUID", "assigned", agents[0].UUID))
		return
	}

	// Execute
	req, _ = http.NewRequest("GET", "/agent/?state=UNASSIGNED", nil)
	res = httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	err = json.Unmarshal([]byte(res.Body.String()), &agents)
	if err != nil {
		t.Error(errorMsg("Agents", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	if agents[0].UUID != "test" {
		t.Error(errorMsg("UUID", "assigned", agents[0].UUID))
		return
	}

	// Execute
	req, _ = http.NewRequest("GET", "/agent/", nil)
	res = httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	err = json.Unmarshal([]byte(res.Body.String()), &agents)
	if err != nil {
		t.Error(errorMsg("Agents", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	if count := len(agents); count != 2 {
		t.Error(errorMsg("len(agents)", "2", fmt.Sprintf("%d", count)))
		return
	}

}

func TestSyncAgentEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	_, _ = repository.CreateNewAgent("test")

	// Execute
	params := GPSData{Lat: "40", Lon: "40", TS: "40"}
	params_json, err := json.Marshal(&params)
	if err != nil {
		t.Error(errorMsg("AgentStruct", "Marshallable", "UnMarshallable"))
	}
	body := bytes.NewBuffer(params_json)
	req, _ := http.NewRequest("POST", "/agent/test/sync", body)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d", res.Code)))
		return
	}

	agent, err := repository.GetAgentByUUID("test")
	if err != nil {
		t.Error(errorMsg("Agent", "ToBeFound", "NotFound"))
		return
	}

	if agent.Lat != "40" {
		t.Error(errorMsg("Lat", "40", agent.Lat))
		return
	}
	if agent.Lon != "40" {
		t.Error(errorMsg("Lon", "40", agent.Lon))
		return
	}
	if agent.TS != "40" {
		t.Error(errorMsg("TS", "40", agent.TS))
		return
	}

}
