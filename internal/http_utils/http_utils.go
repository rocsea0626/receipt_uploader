package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
	"regexp"
	"strings"
)

func SendHealthResponse(w http.ResponseWriter, resp *http_responses.HealthResponse, status int) {
	errMap := map[string]string{
		"message": resp.Message,
	}
	sendJSONResponse(w, &errMap, status)
}

func SendErrorResponse(w http.ResponseWriter, resp *http_responses.ErrorResponse, status int) {
	errMap := map[string]string{
		"error": resp.Error,
	}
	sendJSONResponse(w, &errMap, status)
}

func SendUploadResponse(w http.ResponseWriter, resp *http_responses.UploadResponse) {
	respMap := map[string]string{
		"receiptId": resp.ReceiptID,
	}
	sendJSONResponse(w, &respMap, http.StatusCreated)
}

func SendGetImageResponse(w http.ResponseWriter, fileName string, fileBytes *[]byte) {
	w.Header().Set("Content-Type", "image/jpeg")
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

func ValidateGetImageRequest(r *http.Request, dimensions *configs.Dimensions) (string, string, error) {
	logging.Debugf("ValidateGetImageRequest(r.URL.Path: %s)", r.URL.Path)

	receiptID := strings.TrimPrefix(r.URL.Path, "/receipts/")
	logging.Debugf("receiptID: %s", receiptID)

	re := regexp.MustCompile(`^[a-z0-9]+$`)
	if !re.MatchString(receiptID) {
		return "", "", fmt.Errorf("invalid receiptId")
	}

	size := r.URL.Query().Get("size")

	if size != "" {
		isValidSize := false
		for _, dimension := range *dimensions {
			if dimension.Name == size {
				isValidSize = true
				break
			}
		}
		if !isValidSize {
			return "", "", fmt.Errorf("invalid size parameter, size: %s", size)
		}
	}

	for key := range r.URL.Query() {
		if key != "size" {
			return "", "", fmt.Errorf("unrecognized parameter: %s", key)
		}
	}

	return receiptID, size, nil
}

func sendJSONResponse(w http.ResponseWriter, response *map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*response)
}
