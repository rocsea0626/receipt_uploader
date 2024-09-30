package middlewares

import (
	"net/http"
	"receipt_uploader/internal/constants"
	"receipt_uploader/internal/http_utils"
	"receipt_uploader/internal/models/http_responses"
	"regexp"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usernameToken := r.Header.Get("username_token")

		if !isValidUsernameToken(usernameToken) {
			resp := http_responses.ErrorResponse{
				Error: constants.HTTP_ERR_MSG_403,
			}
			http_utils.SendErrorResponse(w, &resp, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidUsernameToken(token string) bool {
	p := regexp.MustCompile("^[a-z0-9_]+$")
	return p.MatchString(token)
}
