package controller

import (
	"net/http"

	"todo/internal/entity"

	"github.com/go-chi/render"
)

// Login handles the HTTP request for user login.
// It parses the request body to get the password,
// calls the Login method of the TaskController's use case,
// and returns a JSON response with the generated token.
func (tc *TaskController) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := render.Decode(r, &req); err != nil {
		tc.sendError(w, r, entity.ErrUnmarshal)
		return
	}

	token, err := tc.uc.Login(req.Password)
	if err != nil {
		tc.sendError(w, r, err)
		return
	}

	render.JSON(w, r, map[string]string{"token": token})
}
