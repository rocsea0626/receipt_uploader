package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/images"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/resize_queue/resize_queue_mock"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadReceiptHandler(t *testing.T) {
	config := configs.Config{
		UploadsDir: "./mock-uploads",
		ResizedDir: "./mock-images",
		Dimensions: configs.AllowedDimensions,
	}
	userToken := ""

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(config.ResizedDir)
	defer os.RemoveAll(config.UploadsDir)

	imagesService := images.NewService(&config.Dimensions)
	mockResizeQueue := &resize_queue_mock.ServiceMock{}

	t.Run("succeed, POST, 1200x1200 image", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImageJPG(fileName, 1200, 1200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusCreated, status)
	})

	t.Run("should fail, POST, too small width", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImageJPG(fileName, constants.IMAGE_SIZE_MIN_W-1, constants.IMAGE_SIZE_MIN_H)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("should fail, POST, too small height", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImageJPG(fileName, constants.IMAGE_SIZE_MIN_W, constants.IMAGE_SIZE_MIN_H-1)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("should fail, POST, too big image", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImageJPG(fileName, 4000, 4000)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		fileBytes, readErr := os.ReadFile(fileName)
		assert.Nil(t, readErr)
		assert.Greater(t, len(fileBytes), 10<<20)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("should fail, POST, wrong format", func(t *testing.T) {
		fileName := "test_image_save_upload.png"

		createErr := test_utils.CreateTestImageWithFormat(fileName, 4000, 4000, "png")
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		fileBytes, readErr := os.ReadFile(fileName)
		assert.Nil(t, readErr)
		assert.Greater(t, len(fileBytes), 10<<20)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("should fail, not allowed method", func(t *testing.T) {
		req, reqErr := http.NewRequest(http.MethodGet, "/receipts", nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusMethodNotAllowed, status)
	})

	t.Run("should fail, enqueue() failed", func(t *testing.T) {
		fileName := "test_image_enqueue_failed.jpg"
		mockConfig := configs.Config{
			UploadsDir: "./mock-uploads",
			ResizedDir: "./test_image_enqueue_failed",
			Dimensions: configs.AllowedDimensions,
		}
		test_utils.InitTestServer(&mockConfig)
		defer os.RemoveAll(mockConfig.ResizedDir)
		defer os.RemoveAll(mockConfig.UploadsDir)

		userToken := ""

		createErr := test_utils.CreateTestImageJPG(fileName, 1200, 1200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&mockConfig, imagesService, mockResizeQueue)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
