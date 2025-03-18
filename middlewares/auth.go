package middlewares

import (
	"net/http"

	"todo_restapi/internal/myfunctions"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(write http.ResponseWriter, request *http.Request) {

		if err := myfunctions.ValidateJWT(request); err != nil {
			http.Error(write, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(write, request)
	})
}
