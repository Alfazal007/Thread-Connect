package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	if code > 499 {
		log.Println("Client side error, Responding with 5xx", message)
	}
	type ErrType struct {
		Error string `json:"error"`
	}
	RespondWithJson(w, code, ErrType{Error: message})
}
