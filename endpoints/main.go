// Package endpoints Vehicle Tracker API
//
// Vehicle Tracker API provides necessary methods and means to
// track ground and water vehicles positions.
//
//
// Terms Of Service:
// https://api.vehicles.neu.edu.tr/
//
//
//   Schemes: http, https
//   BasePath: /
//
//   Consumes:
//   - application/json
//
//   Produces:
//   - application/json
//
//   SecurityDefinitions:
//     Bearer:
//       type: apiKey
//       name: Authorization
//       in: header
//
// swagger:meta
package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
)

type GenericError struct {
	// Error Message
	//
	// Required: true
	Message string `json:"message"`
}

// Generic Error
// swagger:response ErrorMsg
type ErrorMsg struct {
	// in:body
	Body GenericError
}

func sendErrorMessage(w http.ResponseWriter, message string, status int) {
	errorMsg := ErrorMsg{
		Body: GenericError{Message: message},
	}
	j, err := json.Marshal(errorMsg.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

func sendContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Printf("Error: %s", err.Error())
		sendErrorMessage(w, err.Error(), http.StatusInternalServerError)
	}
}
