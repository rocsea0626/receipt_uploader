package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"receipt_uploader/constants"
	"receipt_uploader/internal/images"
	images_mock "receipt_uploader/internal/images/mock"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReceiptsPostHandler(t *testing.T) {
	config := configs.Config{
		DIR_TMP:    "./mock-tmp",
		DIR_IMAGES: "./mock-images",
	}

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(config.DIR_TMP)
	defer os.RemoveAll(config.DIR_IMAGES)

	imagesService := images.NewService(&config)

	t.Run("succeed, POST, 1200x1200 image", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImage(fileName, 1200, 1200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusCreated, status)
	})

	t.Run("should fail, POST, too small image", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImage(fileName, 300, 200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("should fail, GenerateImages() failed", func(t *testing.T) {
		mockConfig := configs.Config{
			DIR_IMAGES: "mock_generate_images_failed",
		}
		mockImagesService := images_mock.ServiceMock{
			Config: &mockConfig,
		}

		req, reqErr := http.NewRequest(http.MethodPost, "/receipts", nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&mockConfig, &mockImagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		body := rr.Body.String()
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Contains(t, body, constants.HTTP_ERR_MSG_500)
	})

	t.Run("should fail, not allowed method", func(t *testing.T) {
		req, reqErr := http.NewRequest(http.MethodDelete, "/receipts", nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusMethodNotAllowed, status)
	})
}

func TestReceiptsGetHandler(t *testing.T) {
	config := configs.Config{
		DIR_TMP:    "./mock-get-tmp",
		DIR_IMAGES: "./mock-get-images",
	}

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(config.DIR_TMP)
	defer os.RemoveAll(config.DIR_IMAGES)

	imagesService := images.NewService(&config)
	t.Run("return 200, size=small", func(t *testing.T) {
		receiptId := "testgetimage"
		size := "small"
		fileName := receiptId + "_" + size + ".jpg"
		fPath := filepath.Join(config.DIR_IMAGES, fileName)

		createErr := test_utils.CreateTestImage(fPath, 300, 300)
		assert.Nil(t, createErr)

		url := "/receipts/" + receiptId + "?size=small"
		log.Printf("url: %s", url)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&config, imagesService)

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
		handler := ReceiptsHandler(&config, imagesService)

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
		handler := ReceiptsHandler(&config, imagesService)

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
		handler := ReceiptsHandler(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("return 500, GetImage() failed", func(t *testing.T) {
		mockConfig := configs.Config{
			DIR_IMAGES: "mock_get_image_failed",
		}
		mockImagesService := images_mock.ServiceMock{
			Config: &mockConfig,
		}

		receiptId := "mockgetimagefailed"
		size := "medium"

		url := fmt.Sprintf("/receipts/%s?size=%s", receiptId, size)
		req, reqErr := http.NewRequest(http.MethodGet, url, nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&mockConfig, &mockImagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
