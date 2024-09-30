package http_requests

import (
	"fmt"
	"io"
	"net/http"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
)

// UploadRequest represents the incoming request for uploading an image
type UploadRequest struct {
	ContentType string `json:"contentType"`
	Payload     []byte `json:"payload"`
}

type DownloadRequest struct {
	ReceiptId string `json:"receiptId"`
	Size      string `json:"size"`
	Username  string `json:"username"`
}

func ParseUploadRequest(r *http.Request) (*UploadRequest, error) {

	file, header, fromErr := r.FormFile("receipt")
	if fromErr != nil {
		return nil, fmt.Errorf("r.FormFile() failed: %w", fromErr)
	}
	logging.Debugf("content-type: %s", header.Header.Get("Content-Type"))
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return &UploadRequest{
		ContentType: header.Header.Get("Content-Type"),
		Payload:     fileBytes,
	}, nil
}

func ParseDownloadRequest(r *http.Request, dimensions *configs.Dimensions) (*DownloadRequest, error) {

	receiptId, size, err := http_utils.ValidateGetImageRequest(r, dimensions)
	if err != nil {
		return nil, fmt.Errorf("http_utils.ValidateGetImageRequest() failed, err: %s", err.Error())
	}
	username := r.Header.Get("username_token")

	return &DownloadRequest{
		ReceiptId: receiptId,
		Size:      size,
		Username:  username,
	}, nil
}
