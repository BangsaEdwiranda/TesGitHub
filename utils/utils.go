package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func ErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	var payload Payload
	payload.Error = true
	payload.Message = message
	r, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Response failed during JSON encoding with error: %s\n", err.Error())
		return
	}

	w.WriteHeader(httpStatusCode)
	_, err = w.Write(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Response failed with error: %s\n", err.Error())
		return
	}
}

func SuccessResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	var payload Payload
	payload.Error = false
	payload.Message = message
	r, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Response failed during JSON encoding with error: %s\n", err.Error())
		return
	}

	w.WriteHeader(httpStatusCode)
	_, err = w.Write(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Response failed with error: %s\n", err.Error())
		return
	}
}
