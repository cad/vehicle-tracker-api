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
	req, _ := http.NewRequest("GET", "/agents/", nil)
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


func TestSyncAgentEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	_, _ = repository.CreateNewAgent("test")

	// Execute
	params := GPSData{Lat: "40", Lon: "40", TS:"40"}
	params_json, err := json.Marshal(&params)
	if err != nil {
		t.Error(errorMsg("AgentStruct", "Marshallable", "UnMarshallable"))
	}
	body := bytes.NewBuffer(params_json)
	req, _ := http.NewRequest("POST", "/agents/test/sync", body)
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d",res.Code)))
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
