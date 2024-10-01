package handlers

import (
	"net/http"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
)

func UploadReceipt(config *configs.Config, imagesService images.ServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Infof("received request, %s, %s", r.Method, r.URL.Path)

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
	logging.Debugf("handlePost()")
	username := r.Header.Get("username_token")

	bytes, decodeErr := imagesService.ParseImage(r)
	if decodeErr != nil {
		logging.Debugf("imagesService.ParseImage() failed, err: %s", decodeErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_400,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusBadRequest)
		return
	}
	logging.Debugf("len(bytes): %d", len(bytes))

	imgFile, saveErr := imagesService.SaveUpload(&bytes, username, config.UploadsDir)
	if saveErr != nil {
		logging.Errorf("utils.SaveUpload() failed, err: %s", saveErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}
	logging.Infof("image has been saved to path: %s", imgFile.Path)

	receiptID := imgFile.ReceiptID
	resp := http_responses.UploadResponse{
		ReceiptID: receiptID,
	}
	http_utils.SendUploadResponse(w, &resp)
}
