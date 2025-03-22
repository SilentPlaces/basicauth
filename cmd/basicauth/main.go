package main

import (
	"log"
	"net/http"

	middlewareAuth "github.com/SilentPlaces/basicauth.git/internal/middleware/auth"
	"github.com/julienschmidt/httprouter"
)

func main() {
	userController, err := InitializeUserController()
	if err != nil {
		log.Fatalf("failed to initialize user controller: %v", err)
	}

	authService := InitializeAuthService()
	if authService == nil {
		log.Fatal("failed to initialize auth service: authService is nil")
	}

	router := httprouter.New()
	router.GET("/user", middlewareAuth.JWTAuthMiddleware(authService, userController.GetUser))

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
