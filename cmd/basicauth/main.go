package main

import (
	"log"
	"net/http"

	middlewareAuth "github.com/SilentPlaces/basicauth.git/internal/middleware/auth"
	middlewareFormating "github.com/SilentPlaces/basicauth.git/internal/middleware/formating"
	helpers "github.com/SilentPlaces/basicauth.git/pkg/helper/http"
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Initialize controllers and services
	userController, err := InitializeUserController()
	if err != nil {
		log.Fatalf("failed to initialize user controller: %v", err)
	}

	registrationController, err := InitializeRegistrationController()
	if err != nil {
		log.Fatalf("failed to initialize registration controller: %v", err)
	}

	authService := InitializeAuthService()
	if authService == nil {
		log.Fatal("failed to initialize auth service: authService is nil")
	}

	// Initialize router
	router := httprouter.New()

	// Middleware for chaining
	authMiddleware := func(next httprouter.Handle) httprouter.Handle {
		return middlewareAuth.JWTAuthMiddleware(authService, next)
	}

	// Common middleware for response formatting
	responseFormattingMiddleware := middlewareFormating.ResponseFormattingMiddleware

	// Routes with applied middlewares
	//user routes
	router.GET("/user", helpers.ApplyMiddleware(userController.GetUser, authMiddleware, responseFormattingMiddleware))
	router.POST("/auth/login", helpers.ApplyMiddleware(userController.Login, responseFormattingMiddleware))
	router.POST("/auth/refresh-token", helpers.ApplyMiddleware(userController.RefreshToken, responseFormattingMiddleware))

	//Registration routes
	router.POST("/register/init", helpers.ApplyMiddleware(registrationController.SignUp, responseFormattingMiddleware))
	router.POST("/register/verify", helpers.ApplyMiddleware(registrationController.VerifyMail, responseFormattingMiddleware))
	router.POST("/register/resend-resend_verification", helpers.ApplyMiddleware(registrationController.ResendVerification, responseFormattingMiddleware))

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
