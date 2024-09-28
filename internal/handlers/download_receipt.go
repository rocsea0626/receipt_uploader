package handlers

import (
	"log"
	"net/http"
	"os"
	"receipt_uploader/constants"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
)

func DownloadReceipt(config *configs.Config, imagesService images.ServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("r.Method: ", r.Method)
		if http.MethodGet != r.Method {
			resp := http_responses.ErrorResponse{
				Error: constants.HTTP_ERR_MSG_405,
			}
			http_utils.SendErrorResponse(w, &resp, http.StatusMethodNotAllowed)
			return
		}

		handleGet(w, r, config, imagesService)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, config *configs.Config, imagesService images.ServiceType) {
	log.Printf("handleGet(), path: %s", r.URL.Path)

	receiptId, size, validateErr := http_utils.ValidateGetImageRequest(r)
	if validateErr != nil {
		log.Printf("http_utils.ValidateGetImageRequest() failed, err: %s", validateErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_400,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusBadRequest)
		return
	}

	fileBytes, fileName, getErr := imagesService.GetImage(receiptId, size, config.GeneratedDir)
	if getErr != nil {
		log.Printf("images.GetImage() failed, err: %s", getErr.Error())
		if os.IsNotExist(getErr) {
			resp := http_responses.ErrorResponse{
				Error: constants.HTTP_ERR_MSG_404,
			}
			http_utils.SendErrorResponse(w, &resp, http.StatusNotFound)
		}
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	http_utils.SendImageDownloadResponse(w, fileName, &fileBytes)
}
