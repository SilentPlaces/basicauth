package controller

import (
	"encoding/json"
	"github.com/SilentPlaces/basicauth.git/internal/dto/auth/login"
	authdto "github.com/SilentPlaces/basicauth.git/internal/dto/auth/refresh_token"
	mapper "github.com/SilentPlaces/basicauth.git/internal/mappers/users"
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
		RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
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
		helpers.SendErrorResponse(w, http.StatusUnauthorized, "User ID not found ")
		return
	}
	response, err := u.userService.GetUser(userID)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	helpers.WriteJSON(w, http.StatusOK, response)
}

func (u *userController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//Get data
	requestData := login.LoginRequestDTO{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	err := validation.ValidateEmail(requestData.Email)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	//resend_verification auth
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
	helpers.WriteJSON(w, http.StatusOK, &login.LoginResponseDTO{
		User:         userData,
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	})
}

func (u *userController) RefreshToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var requestData authdto.RefreshTokenReqDto
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusBadRequest, "Bad request")
		return
	}
	tokens, err := u.authService.RefreshToken(requestData.RefreshToken)
	if err != nil {
		helpers.SendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
	mappedResponse := mapper.MapTokenToRefreshTokenResDTO(tokens)
	helpers.WriteJSON(w, http.StatusOK, &mappedResponse)
}

var UserControllerProviderSet = wire.NewSet(NewUserController)
