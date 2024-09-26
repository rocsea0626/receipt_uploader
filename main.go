package main

import (
	"fmt"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/models/http_responses"
	"receipt_uploader/internal/utils"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, "Hello, World!")
}

func receiptsHandler(w http.ResponseWriter, r *http.Request) {
	const uploadPath = "./uploads/"

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName, saveErr := utils.SaveUploadImage(r)
	if saveErr != nil {
		http_responses.SendErrorResponse(w, &http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}, http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("File uploaded successfully: %s", fileName)))
}

func main() {
	http.HandleFunc("/health", helloHandler)
	http.HandleFunc("/receipts", receiptsHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
