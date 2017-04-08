// Package endpoints Vehicle Tracker API
//
// Vehicle Tracker API provides necessary methods and means to
// track ground and water vehicles positions.
//
//   Title: Vehicle Tracker API
//   Schemes: http, https
//   Host: localhost
//   BasePath: /
//   Version: 0.0.1
//   License: MIT http://opensource.org/licenses/MIT
//   Contact: Mustafa Arici<mustafa.arici@neu.edu.tr>
//   Terms of Service: http://tracker.neu.edu.tr/
//
//   Consumes:
//   - application/json
//
//   Produces:
//   - application/json
//
//
// swagger:meta
package endpoints

import (
	"net/http"
	"encoding/json"
	"log"
)


func sendErrorMessage(w http.ResponseWriter, message string, status int) {
	j, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}


func sendContentType(w http.ResponseWriter, contentType string){
	w.Header().Set("Content-Type", contentType)
}


func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		sendErrorMessage(w, err.Error(), http.StatusInternalServerError)
	}
}
