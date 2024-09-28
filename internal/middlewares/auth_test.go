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
		req.Header.Set("username_token", "invalid-token")

		rr := httptest.NewRecorder()
		auth.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}

func TestIsValidUsernameToken(t *testing.T) {
	t.Run("valid_token", func(t *testing.T) {
		valid := isValidUsernameToken("valid_token")
		assert.True(t, valid)
	})

	t.Run("user123", func(t *testing.T) {
		valid := isValidUsernameToken("user123")
		assert.True(t, valid)
	})

	t.Run("user_name", func(t *testing.T) {
		valid := isValidUsernameToken("user_name")
		assert.True(t, valid)
	})

	t.Run("user_name123", func(t *testing.T) {
		valid := isValidUsernameToken("user_name123")
		assert.True(t, valid)
	})

	t.Run("empty_token", func(t *testing.T) {
		valid := isValidUsernameToken("")
		assert.False(t, valid)
	})

	t.Run("uppercase_token", func(t *testing.T) {
		valid := isValidUsernameToken("INVALID_TOKEN")
		assert.False(t, valid)
	})

	t.Run("special_character", func(t *testing.T) {
		valid := isValidUsernameToken("username!")
		assert.False(t, valid)
	})

	t.Run("token_with_space", func(t *testing.T) {
		valid := isValidUsernameToken("user name")
		assert.False(t, valid)
	})

	t.Run("another_special_character", func(t *testing.T) {
		valid := isValidUsernameToken("user@name")
		assert.False(t, valid)
	})

	t.Run("only_digits", func(t *testing.T) {
		valid := isValidUsernameToken("123456")
		assert.True(t, valid)
	})

	t.Run("leading_and_trailing_underscores", func(t *testing.T) {
		valid := isValidUsernameToken("__username__")
		assert.True(t, valid)
	})
}
