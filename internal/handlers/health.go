package handlers

import (
	"net/http"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/models/http_responses"
)

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := http_responses.HealthResponse{
			Message: "hello world!",
		}
		http_utils.SendHealthResponse(w, &resp, http.StatusOK)
	}
}
