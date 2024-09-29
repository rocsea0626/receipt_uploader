package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/http_responses"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	logging.Debugf("r.URL.Path: %s", r.URL.Path)

	receiptID := strings.TrimPrefix(r.URL.Path, "/receipts/")
	logging.Debugf("receiptID: %s", receiptID)

	re := regexp.MustCompile(`^[a-z0-9]+$`)
	if !re.MatchString(receiptID) {
		return "", "", fmt.Errorf("invalid receiptId")
	}

	size := r.URL.Query().Get("size")
	if size != constants.IMAGE_SIZE_SMALL && size != constants.IMAGE_SIZE_MEDIUM && size != constants.IMAGE_SIZE_LARGE {
		return "", "", fmt.Errorf("invalid size parameter, size: %s", size)
	}

	return receiptID, size, nil
}

func sendResponse(w http.ResponseWriter, response *map[string]string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(*response)
}

func ParseUploadResponse(t *testing.T, resp *http.Response) *http_responses.UploadResponse {
	respBody, readErr := io.ReadAll(resp.Body)
	assert.Nil(t, readErr)

	var uploadResp http_responses.UploadResponse
	unmarshalErr := json.Unmarshal(respBody, &uploadResp)
	assert.Nil(t, unmarshalErr)

	return &uploadResp
}

func ParseErrorResponse(t *testing.T, resp *http.Response) *http_responses.ErrorResponse {
	respBody, readErr := io.ReadAll(resp.Body)
	assert.Nil(t, readErr)

	var errorResp http_responses.ErrorResponse
	unmarshalErr := json.Unmarshal(respBody, &errorResp)
	assert.Nil(t, unmarshalErr)

	return &errorResp
}

func ParseResponseBody(t *testing.T, resp *http.Response, response interface{}) {
	respBody, readErr := io.ReadAll(resp.Body)
	assert.Nil(t, readErr)

	unmarshalErr := json.Unmarshal(respBody, response)
	assert.Nil(t, unmarshalErr)
}

func ParseDownloadResponseHeader(resp *http.Response) (*http_responses.DownloadResponseHeader, error) {
	contentLen, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	cp := resp.Header.Get("Content-Disposition")
	fileName := strings.TrimPrefix(cp, "attachment; filename=")

	return &http_responses.DownloadResponseHeader{
		Filename:      fileName,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: int64(contentLen),
	}, nil
}
