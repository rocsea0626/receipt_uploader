package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/models/http_responses"
	"receipt_uploader/internal/test_utils"
	"receipt_uploader/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	baseDir := "integ-test-images"
	config := &configs.Config{
		Port:         ":8080",
		UploadedDir:  filepath.Join(baseDir, "uploaded"),
		GeneratedDir: filepath.Join(baseDir, "generated"),
	}
	baseUrl := "http://localhost" + config.Port
	url := baseUrl + "/receipts"
	defer os.RemoveAll(baseDir)

	go utils.StartServer(config)
	client := &http.Client{}

	t.Run("return 200, /health", func(t *testing.T) {
		resp, err := http.Get(baseUrl + "/health")
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("return 201, POST /receipts", func(t *testing.T) {
		uploadFilePath := "./integ-test.jpg"
		userToken := "valid_user"

		test_utils.CreateTestImage(uploadFilePath, 1000, 1200)
		defer os.Remove(uploadFilePath)
		req, reqErr := test_utils.GenerateUploadRequest(t, url, uploadFilePath, userToken)
		assert.Nil(t, reqErr)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var uploadResp http_responses.UploadResponse
		http_utils.ParseResponseBody(t, resp, &uploadResp)
		assert.NotEmpty(t, uploadResp.ReceiptID)
	})

	t.Run("return 400, POST /receipts, empty image", func(t *testing.T) {
		uploadFilePath := "./integ-test.jpg"
		userToken := "valid_user"

		test_utils.CreateTestImage(uploadFilePath, 0, 0)
		defer os.Remove(uploadFilePath)
		req, reqErr := test_utils.GenerateUploadRequest(t, url, uploadFilePath, userToken)
		assert.Nil(t, reqErr)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var errorResp http_responses.ErrorResponse
		http_utils.ParseResponseBody(t, resp, &errorResp)
		assert.Equal(t, constants.HTTP_ERR_MSG_400, errorResp.Error)
	})

	t.Run("return 403, username_token missing", func(t *testing.T) {
		uploadFilePath := "./integ-test.jpg"
		userToken := ""

		test_utils.CreateTestImage(uploadFilePath, 1000, 1200)
		defer os.Remove(uploadFilePath)
		req, reqErr := test_utils.GenerateUploadRequest(t, url, uploadFilePath, userToken)
		assert.Nil(t, reqErr)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("return 200, GET /receipts/{receiptId}?size=large", func(t *testing.T) {
		uploadFilePath := "./integ-test.jpg"
		size := "large"
		userToken := "valid_user"

		test_utils.CreateTestImage(uploadFilePath, 1000, 1200)
		defer os.Remove(uploadFilePath)
		req, reqErr := test_utils.GenerateUploadRequest(t, url, uploadFilePath, userToken)
		assert.Nil(t, reqErr)

		orgFileSize, sizeErr := test_utils.GetFileSize(uploadFilePath)
		assert.Nil(t, sizeErr)

		resp, err := client.Do(req)
		assert.Nil(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var uploadResp http_responses.UploadResponse
		http_utils.ParseResponseBody(t, resp, &uploadResp)
		assert.NotEmpty(t, uploadResp.ReceiptID)

		getUrl := fmt.Sprintf("%s/%s?size=%s", url, uploadResp.ReceiptID, size)
		log.Println("getUrl: ", getUrl)
		getReq, getReqErr := http.NewRequest(http.MethodGet, getUrl, nil)
		getReq.Header.Set("username_token", userToken)
		assert.Nil(t, getReqErr)

		getResp, getErr := client.Do(getReq)
		assert.Nil(t, getErr)
		defer getResp.Body.Close()

		assert.Equal(t, http.StatusOK, getResp.StatusCode)
		header, headerErr := http_utils.ParseDownloadResponseHeader(getResp)
		assert.Nil(t, headerErr)
		assert.Equal(t, "application/octet-stream", header.ContentType)

		fileName := uploadResp.ReceiptID + "_" + size + ".jpg"
		assert.Equal(t, fileName, header.Filename)

		getRespBody, getReadErr := io.ReadAll(getResp.Body)
		assert.Nil(t, getReadErr)
		assert.Equal(t, orgFileSize, int64(len(getRespBody)))
		assert.Equal(t, header.ContentLength, int64(len(getRespBody)))
	})

	t.Run("return 403, GET /receipts/{receiptId}?size=large, username_token missing", func(t *testing.T) {

		getUrl := fmt.Sprintf("%s/%s?size=%s", url, "fakereceiptId", "small")
		log.Println("getUrl: ", getUrl)
		getReq, getReqErr := http.NewRequest(http.MethodGet, getUrl, nil)
		assert.Nil(t, getReqErr)

		getResp, getErr := client.Do(getReq)
		assert.Nil(t, getErr)
		defer getResp.Body.Close()

		assert.Equal(t, http.StatusForbidden, getResp.StatusCode)
	})

	t.Run("return 403, GET /receipts/{receiptId}?size=large, token has wrong key", func(t *testing.T) {

		getUrl := fmt.Sprintf("%s/%s?size=%s", url, "fakereceiptId", "small")
		log.Println("getUrl: ", getUrl)
		getReq, getReqErr := http.NewRequest(http.MethodGet, getUrl, nil)
		getReq.Header.Set("wrong_token_key", "username_token_val")

		assert.Nil(t, getReqErr)

		getResp, getErr := client.Do(getReq)
		assert.Nil(t, getErr)
		defer getResp.Body.Close()

		assert.Equal(t, http.StatusForbidden, getResp.StatusCode)
	})
}
