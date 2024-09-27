package http_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/models/http_responses"
)

func DecodeImage(r *http.Request) ([]byte, error) {

	parseErr := r.ParseMultipartForm(10 << 20) // Maximum 10 MB
	if parseErr != nil {
		return nil, fmt.Errorf("r.ParseMultipartForm() failed, err: %s", parseErr.Error())
	}

	file, header, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return nil, fmt.Errorf("r.FormFile() failed: %w", fromErr)
	}
	log.Printf("content-type: %s", header.Header.Get("Content-Type"))
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	img, format, decodeErr := image.Decode(bytes.NewReader(fileBytes))
	if decodeErr != nil {
		return nil, fmt.Errorf("image.Decode() failed, err: %s", decodeErr.Error())
	}
	if img.Bounds().Dx() < constants.IMAGE_SIZE_MIN_W || img.Bounds().Dy() < constants.IMAGE_SIZE_MIN_H {
		return nil, fmt.Errorf("invalid image size, minHeight=%d, minWidth=%d", constants.IMAGE_SIZE_MIN_H, constants.IMAGE_SIZE_MIN_W)
	}

	if format != "jpeg" {
		return nil, fmt.Errorf("invalid image format, format=%s, only jpeg format is allowed", format)
	}

	return fileBytes, nil
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
