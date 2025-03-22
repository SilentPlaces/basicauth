package controller

import (
	"encoding/json"
	middleware "github.com/SilentPlaces/basicauth.git/internal/middleware/auth"
	userService "github.com/SilentPlaces/basicauth.git/internal/services/users"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type (
	UserController interface {
		GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params)
	}

	userController struct {
		userService userService.UserService
	}
)

func NewUserController(userService userService.UserService) UserController {
	return &userController{userService: userService}
}

func (u *userController) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Extract user id from the request context.
	userID, ok := r.Context().Value(middleware.UserContextKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	response, err := u.userService.GetUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var UserControllerProviderSet = wire.NewSet(NewUserController)
