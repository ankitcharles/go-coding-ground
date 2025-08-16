package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client supports gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// If not, just call the next handler
			next.ServeHTTP(w, r)
			return
		}

		// Set the Content-Encoding header to gzip
		w.Header().Set("Content-Encoding", "gzip")

		// Create a gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap the ResponseWriter to use the gzip writer
		gzw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}

		// Call the next handler with the wrapped ResponseWriter
		next.ServeHTTP(gzw, r)
		fmt.Println("Compression Middleware Completed")
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
