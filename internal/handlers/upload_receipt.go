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

func UploadReceiptHandler(config *configs.Config, imagesService images.ServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("r.Method: ", r.Method)
		if http.MethodPost != r.Method {
			resp := http_responses.ErrorResponse{
				Error: constants.HTTP_ERR_MSG_405,
			}
			http_utils.SendErrorResponse(w, &resp, http.StatusMethodNotAllowed)
			return
		}

		handlePost(w, r, config, imagesService)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request, config *configs.Config, imagesService images.ServiceType) {
	log.Println("handlePost()")
	// username := r.Header.Get("username_token")

	bytes, decodeErr := imagesService.DecodeImage(r)
	if decodeErr != nil {
		log.Printf("http_utils.DecodeImage() failed, err: %s", decodeErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_400,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusBadRequest)
		return
	}

	tmpFilePath, saveErr := imagesService.SaveUpload(bytes, config.UploadedDir)
	if saveErr != nil {
		log.Printf("utils.SaveUpload() failed, err: %s", saveErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	destDir := config.GeneratedDir
	genErr := imagesService.GenerateImages(tmpFilePath, destDir)
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
