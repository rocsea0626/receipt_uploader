package http_responses

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type UploadResponse struct {
	ReceiptID string `json:"receiptId"`
	ImageURL  string `json:"imageUrl"`
}

func SendErrorResponse(w http.ResponseWriter, resp *ErrorResponse, status int) {
	errMap := map[string]string{
		"error": resp.Error,
	}
	SendResponse(w, &errMap, status)
}

func SendUploadResponse(w http.ResponseWriter, resp *UploadResponse, status int) {
	errMap := map[string]string{
		"receiptId": resp.ReceiptID,
		"imageUrl":  resp.ImageURL,
	}
	SendResponse(w, &errMap, status)
}

func SendResponse(w http.ResponseWriter, response *map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*response)
}
