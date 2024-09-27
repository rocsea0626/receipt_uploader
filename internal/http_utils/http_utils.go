package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/models/http_responses"
	"strings"
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

func SendImageDownloadResponse(w http.ResponseWriter, fileName string, fileBytes *[]byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(*fileBytes)))

	reader := bytes.NewReader(*fileBytes)
	_, err := io.Copy(w, reader)
	if err != nil {
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		SendErrorResponse(w, &resp, http.StatusInternalServerError)
	}
}

func ValidateGetImageRequest(r *http.Request) (string, string, error) {
	receiptID := strings.TrimPrefix(r.URL.Path, "/receipts/")

	size := r.URL.Query().Get("size")
	if size != "small" && size != "medium" && size != "large" {
		return "", "", fmt.Errorf("invalid size parameter, size: %s", size)
	}

	return receiptID, size, nil
}

func sendResponse(w http.ResponseWriter, response *map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*response)
}
