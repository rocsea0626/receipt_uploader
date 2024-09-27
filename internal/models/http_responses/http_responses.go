package http_responses

type ErrorResponse struct {
	Error string `json:"error"`
}

type UploadResponse struct {
	ReceiptID string `json:"receiptId"`
}
