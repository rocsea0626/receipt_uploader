package main

import (
	"fmt"
	"log"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/futils"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/models/http_responses"
	"receipt_uploader/internal/utils"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintln(w, "Hello, World!")
}

func receiptsHandler(w http.ResponseWriter, r *http.Request) {
	tmpDir := "tmp"

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath, saveErr := utils.SaveUploadedImage(r, tmpDir)
	if saveErr != nil {
		log.Printf("utils.SaveUploadedImage() failed, err: %s", saveErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	genErr := utils.GenerateImages(filePath, constants.OUTPUT_DIR)
	if genErr != nil {
		log.Printf("utils.GenerateImages() failed, err: %s", genErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	receiptID := futils.GetFileName(filePath)
	resp := http_responses.UploadResponse{
		ReceiptID: receiptID,
	}
	http_utils.SendUploadResponse(w, &resp)
}

func main() {
	log.Println("initializing server")
	utils.InitServer()

	http.HandleFunc("/health", helloHandler)
	http.HandleFunc("/receipts", receiptsHandler)

	fmt.Printf("Starting server on %s", constants.PORT)
	if err := http.ListenAndServe(constants.PORT, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
