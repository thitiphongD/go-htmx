package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func PrintJSONResponse(data interface{}, statusCode int) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Printf("Response (Status: %d): %s\n", statusCode, string(jsonData))
}
