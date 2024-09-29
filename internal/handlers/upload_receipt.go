package handlers

import (
	"net/http"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
	"receipt_uploader/internal/utils"
)

func UploadReceipt(config *configs.Config, imagesService images.ServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Debugf("r.Method: %s", r.Method)

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
		logging.Debugf("http_utils.DecodeImage() failed, err: %s", decodeErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_400,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusBadRequest)
		return
	}

	orgFilePath, saveErr := imagesService.SaveUpload(&bytes, config.ImagesDir)
	if saveErr != nil {
		logging.Errorf("utils.SaveUpload() failed, err: %s", saveErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	destDir := filepath.Join(config.ImagesDir, username)
	genErr := imagesService.GenerateImages(&bytes, orgFilePath, destDir)
	if genErr != nil {
		logging.Errorf("images.GenerateImages() failed, err: %s", genErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_500,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusInternalServerError)
		return
	}

	receiptID := utils.GetFileName(orgFilePath)
	resp := http_responses.UploadResponse{
		ReceiptID: receiptID,
	}
	http_utils.SendUploadResponse(w, &resp)
}
