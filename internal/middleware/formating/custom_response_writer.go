package middleware

import (
	"bytes"
	"net/http"
)

// CustomResponseWriter that implements http.ResponseWriter to wrap it and use in middlewares
type CustomResponseWriter struct {
	http.ResponseWriter
	bodyBuffer bytes.Buffer
	statusCode int
}

func (cw *CustomResponseWriter) Write(b []byte) (int, error) {
	return cw.bodyBuffer.Write(b)
}

func (cw *CustomResponseWriter) WriteHeader(statusCode int) {
	cw.statusCode = statusCode
}

func (cw *CustomResponseWriter) Header() http.Header {
	return cw.ResponseWriter.Header()
}
