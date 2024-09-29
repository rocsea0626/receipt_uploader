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
	config := configs.Config{
		ImagesDir: "./mock-get-images",
	}

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(config.ImagesDir)

	imagesService := images.NewService()
	t.Run("return 200, size=small", func(t *testing.T) {
		receiptId := "testgetimage"
		size := "small"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(config.ImagesDir, fileName)

		createErr := test_utils.CreateTestImage(fPath, 300, 300)
		assert.Nil(t, createErr)

		url := "/receipts/" + receiptId + "?size=small"
		logging.Debugf("url: %s", url)

		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := DownloadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusOK, status)
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
			ImagesDir: "mock_get_image_failed",
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
