package http

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// ApplyMiddleware apply http middleware to h as input Handle
func ApplyMiddleware(h httprouter.Handle, middlewares ...func(httprouter.Handle) httprouter.Handle) httprouter.Handle {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

// SendErrorResponse log and send http error with given message and code
func SendErrorResponse(w http.ResponseWriter, status int, message string) {
	log.Printf("Status Code : %d, Error: %s", status, message)
	http.Error(w, message, status)
}
