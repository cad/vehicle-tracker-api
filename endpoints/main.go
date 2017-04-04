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
