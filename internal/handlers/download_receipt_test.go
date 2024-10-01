package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"receipt_uploader/internal/images"
	images_mock "receipt_uploader/internal/images/mock"
	"receipt_uploader/internal/logging"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadReceiptHandler(t *testing.T) {
	baseDir := "test-get"
	config := configs.Config{
		ResizedDir: filepath.Join(baseDir, "resized"),
		UploadsDir: filepath.Join(baseDir, "uploads"),
		Dimensions: configs.AllowedDimensions,
	}

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(baseDir)

	imagesService := images.NewService(&config.Dimensions)
	t.Run("return 200, size=small", func(t *testing.T) {
		username := "test-user-get"
		receiptId := "testrecieptid"
		size := "small"
		fileName := receiptId + "_" + size + ".jpg"

		userDir := filepath.Join(config.ResizedDir, username)
		os.MkdirAll(userDir, 0755)
		fPath := filepath.Join(userDir, fileName)

		createErr := test_utils.CreateTestImageJPG(fPath, 300, 300)
		assert.Nil(t, createErr)

		url := "/receipts/" + receiptId + "?size=small"
		logging.Debugf("url: %s", url)

		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)
		req.Header.Set("username_token", username)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusOK, status)
	})

	t.Run("return 200, size is empty", func(t *testing.T) {
		username := "test-user-get"
		receiptId := "testrecieptid"
		fileName := receiptId + ".jpg"

		userDir := filepath.Join(config.ResizedDir, username)
		os.MkdirAll(userDir, 0755)
		fPath := filepath.Join(userDir, fileName)

		createErr := test_utils.CreateTestImageJPG(fPath, 300, 300)
		assert.Nil(t, createErr)

		url := "/receipts/" + receiptId + "?size"
		logging.Debugf("url: %s", url)

		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)
		req.Header.Set("username_token", username)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusOK, status)

		url1 := "/receipts/" + receiptId + "?size="
		logging.Debugf("url: %s", url)

		req1, reqErr1 := http.NewRequest(http.MethodGet, url1, nil)
		assert.Nil(t, reqErr1)
		req.Header.Set("username_token", username)
		handler.ServeHTTP(rr, req1)

		status1 := rr.Code
		assert.Equal(t, http.StatusOK, status1)
	})

	t.Run("return 404, not found by receiptId", func(t *testing.T) {
		receiptId := "notfound"
		size := "medium"

		url := fmt.Sprintf("/receipts/%s?size=%s", receiptId, size)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("return 400, invalid receiptId, receiptId=12.34", func(t *testing.T) {
		receiptId := "12.34"
		size := "larage"

		url := fmt.Sprintf("/receipts/%s?size=%s", receiptId, size)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("return 400, invalid receiptId, receiptId=Ab1234", func(t *testing.T) {
		receiptId := "Ab1234"
		size := "larage"

		url := fmt.Sprintf("/receipts/%s?size=%s", receiptId, size)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("return 400, invalid szie, size=Xs", func(t *testing.T) {
		receiptId := "1234"
		size := "Xs"

		url := fmt.Sprintf("/receipts/%s?size=%s", receiptId, size)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("return 500, GetImage() failed", func(t *testing.T) {
		mockConfig := configs.Config{
			ResizedDir: "mock_get_image_failed",
			Dimensions: configs.AllowedDimensions,
		}
		mockImagesService := images_mock.ServiceMock{}

		receiptId := "mockgetimagefailed"
		size := "medium"

		url := fmt.Sprintf("/receipts/%s?size=%s", receiptId, size)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&mockConfig, &mockImagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
