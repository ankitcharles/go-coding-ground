package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Response Time Middleware Invoked")
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		duration := time.Since(start)
		rw.Header().Set("X-Response-Time", duration.String())
		if r.Method == http.MethodGet {
			next.ServeHTTP(rw, r)
		} else {
			rw.statusCode = http.StatusMethodNotAllowed
		}
		// Log the response time
		duration = time.Since(start)
		log.Printf("Method: %s, Path: %s, Status: %d, Duration: %s", r.Method, r.URL.Path, rw.statusCode, duration)
		// Set the response time in the header
		fmt.Println("Response Time Middleware Completed")

	})
}

// custom response writer to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
