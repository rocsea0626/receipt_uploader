package handlers

import (
	"net/http"
	"os"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_requests"
	"receipt_uploader/internal/models/http_responses"
	"receipt_uploader/internal/models/image_meta"
)

func DownloadReceipt(config *configs.Config, imagesService images.ServiceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logging.Infof("received request, %s, %s", r.Method, r.URL.Path)

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
	logging.Debugf("handleGet(), path: %s", r.URL.Path)

	downloadReq, parseErr := http_requests.ParseDownloadRequest(r, &config.Dimensions)
	if parseErr != nil {
		logging.Errorf("http_requests.ParseDownloadRequest() failed, err: %s", parseErr.Error())
		resp := http_responses.ErrorResponse{
			Error: constants.HTTP_ERR_MSG_400,
		}
		http_utils.SendErrorResponse(w, &resp, http.StatusBadRequest)
		return
	}

	imageMeta := image_meta.FromGetRequset(downloadReq.ReceiptId, downloadReq.Size, downloadReq.Username, config.ResizedDir)
	fileBytes, fileName, getErr := imagesService.GetImage(imageMeta)
	if getErr != nil {
		logging.Errorf("images.GetImage() failed, err: %s", getErr.Error())
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

	logging.Infof("response with image: %s", imageMeta.FileName)
	http_utils.SendGetImageResponse(w, fileName, &fileBytes)
}
