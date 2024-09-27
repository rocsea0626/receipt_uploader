package http_utils

import (
	"encoding/json"
	"net/http"
	"receipt_uploader/internal/models/http_responses"
)

func SendHealthResponse(w http.ResponseWriter, resp *http_responses.HealthResponse, status int) {
	errMap := map[string]string{
		"message": resp.Message,
	}
	sendResponse(w, &errMap, status)
}

func SendErrorResponse(w http.ResponseWriter, resp *http_responses.ErrorResponse, status int) {
	errMap := map[string]string{
		"error": resp.Error,
	}
	sendResponse(w, &errMap, status)
}

func SendUploadResponse(w http.ResponseWriter, resp *http_responses.UploadResponse) {
	respMap := map[string]string{
		"receiptId": resp.ReceiptID,
	}
	sendResponse(w, &respMap, http.StatusCreated)
}

func sendResponse(w http.ResponseWriter, response *map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*response)
}
