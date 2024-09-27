package http_responses

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Message string `json:"message"`
}

type UploadResponse struct {
	ReceiptID string `json:"receiptId"`
}
