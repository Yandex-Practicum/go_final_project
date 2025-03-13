package controller

import (
	"net/http"

	"todo/internal/entity"
)

// ValidateTokenMiddleware is a middleware function that validates the token in the request cookie.
// If the token is valid, it allows the request to proceed to the next handler.
// If the token is invalid or missing, it sends an authentication required error response.
func (tc *TaskController) ValidateTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tc.uc.ShouldCheckToken() {
			next.ServeHTTP(w, r)
			return
		}

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			tc.sendError(w, r, entity.ErrAuthenticationRequired)
			return
		}

		if err = tc.uc.ValidateToken(tokenCookie.Value); err != nil {
			tc.sendError(w, r, entity.ErrAuthenticationRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}
