package endpoints

import (
	"fmt"
	"os"
	"net/http"
	"net/http/httptest"
//	"encoding/json"
	"testing"
//	"bytes"
	"github.com/cad/vehicle-tracker-api/repository"
)


func TestCORSPreflight(t *testing.T) {
	// Init
	repository.ConnectDB("sqlite3", "/tmp/test.db")
	defer repository.CloseDB()
	defer os.Remove("/tmp/test.db")

	// Prepare


	// Execute
	req, _ := http.NewRequest("OPTIONS", "/vehicle/", nil)
	req.Header.Set("Origin", "example.com")
	res := httptest.NewRecorder()
	GetRouter().ServeHTTP(res, req)

	// Test
	if res.Header().Get("Access-Control-Allow-Methods") != "POST, GET, OPTIONS, PUT, DELETE" {
		fmt.Println(res.Header())
		t.Error(errorMsg("Accss-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE", res.Header().Get("Access-Control-Allow-Methods")))
		return
	}

	if res.Header().Get("Access-Control-Allow-Origin") != "example.com" {
		fmt.Println(res.Header())
		t.Error(errorMsg("Accss-Control-Allow-Origin", "example.com", res.Header().Get("Access-Control-Allow-Origin")))
		return
	}

}
