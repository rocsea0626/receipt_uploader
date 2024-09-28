package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	auth := Auth(testHandler)

	t.Run("succeed, token=user_123", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/some-endpoint", nil)
		req.Header.Set("username_token", "username_token")

		rr := httptest.NewRecorder()
		auth.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should fail, empty token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/some-endpoint", nil)
		req.Header.Set("username_token", "")

		rr := httptest.NewRecorder()
		auth.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("should fail, token has only whitespaces", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/some-endpoint", nil)
		req.Header.Set("username_token", "   ")

		rr := httptest.NewRecorder()
		auth.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})

	t.Run("should fail, invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/some-endpoint", nil)
		req.Header.Set("username_token", "invalid_token")

		rr := httptest.NewRecorder()
		auth.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}
