package http_utils_test

import (
	"net/http/httptest"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/models/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateGetImageRequest(t *testing.T) {

	t.Run("succeed, small=size", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/12345?size=small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.Nil(t, err)
		assert.Equal(t, "12345", receiptID)
		assert.Equal(t, "small", size)
	})

	t.Run("succeed, small=medium", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/67890?size=medium", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.Nil(t, err)
		assert.Equal(t, "67890", receiptID)
		assert.Equal(t, "medium", size)
	})

	t.Run("succeed, small=large", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/246e80?size=large", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.Nil(t, err)
		assert.Equal(t, "246e80", receiptID)
		assert.Equal(t, "large", size)
	})

	t.Run("succeed, no size parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/246e80", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.Nil(t, err)
		assert.Equal(t, "246e80", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, escape /", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts%2F456a8?size=Small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, missing receiptId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/?size=small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, invalid size parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/10293?size=extra-large", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)

	})

	t.Run("should fail, invalid query parameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/45678?resolution=small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, invalid receiptId", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/456-78?size=small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, path with multiple slashes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts///45678?size=small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, path with space", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/456%2078?size=small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, case sensitive, size=Small", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/45678?size=Small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})

	t.Run("should fail, case sensitive, uuid=A45678", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/receipts/A45678å?size=Small", nil)
		receiptID, size, err := http_utils.ValidateGetImageRequest(req, &configs.AllowedDimensions)

		assert.NotNil(t, err)
		assert.Equal(t, "", receiptID)
		assert.Equal(t, "", size)
	})
}
