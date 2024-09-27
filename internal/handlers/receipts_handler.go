// receipts_handler.go

package handlers

import (
	"log"
	"net/http"
	"receipt_uploader/constants"
	"receipt_uploader/internal/futils"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
)

func ReceiptsHandler(config *configs.Config, imagesService images.ServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r, config)
		case http.MethodPost:
			handlePost(w, r, imagesService)
		default:
			resp := http_responses.ErrorResponse{
				Error: constants.HTTP_ERR_MSG_405,
			}
			http_utils.SendErrorResponse(w, &resp, http.StatusMethodNotAllowed)
		}
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, imagesService images.ServiceType) {
	bytes, decodeErr := imagesService.DecodeImage(r)
	if decodeErr != nil {
		log.Printf("http_utils.DecodeImage() failed, err: %s", decodeErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_400,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusBadRequest)
	}

	tmpFilePath, saveErr := imagesService.SaveUpload(bytes)
	if saveErr != nil {
		log.Printf("utils.SaveUpload() failed, err: %s", saveErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	genErr := imagesService.GenerateImages(tmpFilePath)
	if genErr != nil {
		log.Printf("images.GenerateImages() failed, err: %s", genErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	receiptID := futils.GetFileName(tmpFilePath)
	resp := http_responses.UploadResponse{
		ReceiptID: receiptID,
	}
	http_utils.SendUploadResponse(w, &resp)
}

func handleGet(w http.ResponseWriter, r *http.Request, config *configs.Config) {
	// Implement the logic to handle GET requests here
	// For example, you might return a list of receipts, or a specific receipt's information
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("GET method is not yet implemented"))
}
