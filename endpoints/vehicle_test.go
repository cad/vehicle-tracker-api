package endpoints

import (
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"testing"
	"bytes"
//	"runtime/debug"
	api "github.com/cad/vehicle-tracker-api"
	"github.com/cad/vehicle-tracker-api/repository"
	"github.com/cad/vehicle-tracker-api/endpoints"
)


func errorMsg(what, shouldBe, was string) string {
	//debug.PrintStack()
	return fmt.Sprintf("expected %s \"%s\" but got \"%s\"", what, shouldBe, was)
}


func TestGetAllVehiclesEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	repository.CreateVehicle("test")

	// Execute
	req, _ := http.NewRequest("GET", "/vehicles/", nil)
	res := httptest.NewRecorder()
	api.GetServer().ServeHTTP(res, req)

	// Test
	var vehicles []repository.Vehicle
	err := json.Unmarshal([]byte(res.Body.String()), &vehicles)
	if err != nil {
		t.Error(errorMsg("Vehicles", "Unmarshallable", "NotUnmarshallable"))
		return
	}

	if vehicles[0].PlateID != "test" {
		t.Error(errorMsg("PlateID", "test", vehicles[0].PlateID))
		return
	}
}


func TestCreateNewVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	type vehicleStruct struct {
		PlateID string `json:"plate_id"`
	}
	vehicle := vehicleStruct{PlateID: "test"}
	vehicle_json, err := json.Marshal(&vehicle)
	if err != nil {
		t.Error(errorMsg("Vehicle", "Marshallable", "NotMarshallable"))
	}
	body := bytes.NewBuffer(vehicle_json)

	// Execute
	req, _ := http.NewRequest("POST", "/vehicles/", body)
	res := httptest.NewRecorder()
	api.GetServer().ServeHTTP(res, req)

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
	repository.CreateVehicle("test")

	// Execute
	req, _ := http.NewRequest("GET", "/vehicles/test", nil)
	res := httptest.NewRecorder()
	api.GetServer().ServeHTTP(res, req)

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

}


func TestSyncVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	repository.CreateVehicle("test")

	// Execute
	vehicleSyncStruct := endpoints.VehicleSyncStruct{Lat: 40.0, Lon: 40.0}
	vehicle_json, err := json.Marshal(&vehicleSyncStruct)
	if err != nil {
		t.Error(errorMsg("VehicleStruct", "Marshallable", "UnMarshallable"))
	}
	body := bytes.NewBuffer(vehicle_json)
	req, _ := http.NewRequest("POST", "/vehicles/test/sync", body)
	res := httptest.NewRecorder()
	api.GetServer().ServeHTTP(res, req)

	// Test
	if res.Code != 200 {
		t.Error(errorMsg("StatusCode", "200", fmt.Sprintf("%d",res.Code)))
		return
	}

	vehicle, err := repository.GetVehicleByPlateID("test")
	if err != nil {
		t.Error(errorMsg("Vehicle", "ToBeFound", "NotFound"))
		return
	}

	if vehicle.Lat != 40.0 {
		t.Error(errorMsg("Lat", "40.0", fmt.Sprintf("%d", vehicle.Lat)))
		return
	}
	if vehicle.Lon != 40.0 {
		t.Error(errorMsg("Lon", "40.0", fmt.Sprintf("%d", vehicle.Lon)))
		return
	}
}


func TestDeleteVehicleEndpoint(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare
	repository.CreateVehicle("test")

	// Execute
	req, _ := http.NewRequest("DELETE", "/vehicles/test", nil)
	res := httptest.NewRecorder()
	api.GetServer().ServeHTTP(res, req)

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
