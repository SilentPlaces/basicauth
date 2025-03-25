package controller

import (
	"encoding/json"
	login_dto "github.com/SilentPlaces/basicauth.git/internal/dto/login"
	middleware "github.com/SilentPlaces/basicauth.git/internal/middleware/auth"
	authService "github.com/SilentPlaces/basicauth.git/internal/services/auth"
	userService "github.com/SilentPlaces/basicauth.git/internal/services/users"
	validation "github.com/SilentPlaces/basicauth.git/internal/validation/user"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/http"
	"github.com/google/wire"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type (
	UserController interface {
		GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
		Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	}

	userController struct {
		userService userService.UserService
		authService authService.AuthService
	}
)

func NewUserController(userService userService.UserService, authService authService.AuthService) UserController {
	return &userController{
		userService: userService,
		authService: authService,
	}
}

// GetUser is controller to get user data
func (u *userController) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

func (u *userController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//Get data
	requestData := login_dto.LoginRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	err := validation.ValidateEmail(requestData.Email)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	//verify login
	userData, err := u.userService.VerifyLogin(requestData.Email, requestData.Password)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Wrong Email or Password")
		return
	}
	//generate token and refresh token
	token, err := u.authService.GenerateToken(userData.ID)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//send final
	helpers.WriteJSON(w, http.StatusOK, &login_dto.LoginResponseDTO{
		User:         userData,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

var UserControllerProviderSet = wire.NewSet(NewUserController)
