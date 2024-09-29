package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"receipt_uploader/constants"
	"receipt_uploader/internal/images"
	images_mock "receipt_uploader/internal/images/mock"
	"receipt_uploader/internal/models/configs"
	"receipt_uploader/internal/test_utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadReceiptHandler(t *testing.T) {
	config := configs.Config{
		ImagesDir: "./mock-images",
	}
	userToken := ""

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(config.ImagesDir)

	imagesService := images.NewService()

	t.Run("succeed, POST, 1200x1200 image", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImage(fileName, 1200, 1200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusCreated, status)
	})

	t.Run("should fail, POST, too small image", func(t *testing.T) {
		fileName := "test_image_save_upload.jpg"

		createErr := test_utils.CreateTestImage(fileName, 300, 200)
		assert.Nil(t, createErr)
		defer os.Remove(fileName)

		req, reqErr := test_utils.GenerateUploadRequest(t, "/receipts", fileName, userToken)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("should fail, GenerateImages() failed", func(t *testing.T) {
		mockConfig := configs.Config{
			ImagesDir: "mock_generate_images_failed",
		}
		mockImagesService := images_mock.ServiceMock{}

		req, reqErr := http.NewRequest(http.MethodPost, "/receipts", nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		handler := UploadReceipt(&mockConfig, &mockImagesService)

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
		handler := UploadReceipt(&config, imagesService)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		assert.Equal(t, http.StatusMethodNotAllowed, status)
	})
}
