package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// ErrorResponse defines the standard error structure
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// methodNotAllowedHandler handles disallowed HTTP methods
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed")
	log.Printf("[405] %s %s - Method Not Allowed", r.Method, r.URL.Path)
}

// notFoundHandler handles undefined routes
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeErrorResponse(w, http.StatusNotFound, "route_not_found")
	log.Printf("[404] %s %s - Route Not Found", r.Method, r.URL.Path)
}

// writeErrorResponse writes a standardized JSON error response
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

// loggingMiddleware logs all incoming requests and their response status
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		startTime := time.Now()

		next.ServeHTTP(lrw, r)

		duration := time.Since(startTime)
		status := "success"
		if lrw.statusCode >= 400 {
			status = "error"
		}

		log.Printf("[%d] %s %s - %s - Status: %s - Duration: %s - Response: %s",
			lrw.statusCode, r.Method, r.URL.Path, http.StatusText(lrw.statusCode), status, duration, lrw.body.String())
	})
}

// loggingResponseWriter is a wrapper to capture status codes and response body
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

// RegisterRoutes defines all API endpoints
func RegisterRoutes(router *mux.Router) {
	router.Use(loggingMiddleware)

	// Register Auth and User Routes
	RegisterAuthRoutes(router)
	RegisterUserRoutes(router)

	// Custom 404 Not Found Handler
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
}
