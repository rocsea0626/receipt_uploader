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

func TestReceiptsHandler(t *testing.T) {
	config := configs.Config{
		DIR_TMP:    "./mock-tmp",
		DIR_IMAGES: "./mock-images",
	}

	test_utils.InitTestServer(&config)
	defer os.RemoveAll(config.DIR_TMP)
	defer os.RemoveAll(config.DIR_IMAGES)

	imagesService := images.NewService()

	t.Run("GET method should return 501 Not Implemented", func(t *testing.T) {
		t.Skip()
		req, err := http.NewRequest(http.MethodGet, "/receipts", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := ReceiptsHandler(&config, imagesService)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotImplemented {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotImplemented)
		}

		expected := "GET method is not yet implemented"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})

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
		req, reqErr := http.NewRequest(http.MethodPost, "/receipts", nil)
		assert.Nil(t, reqErr)

		rr := httptest.NewRecorder()
		mockImagesService := images_mock.ServiceMock{}
		mockConfig := configs.Config{
			DIR_IMAGES: "mock_generate_images_failed",
		}
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
