package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	mw "restapi-go/internal/api/middlewares"
	"restapi-go/internal/api/routers"
	"restapi-go/internal/repository/sqlconnect"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := sqlconnect.ConnectDb("test_database")
	if err != nil {
		fmt.Println("Error connecting to database", err)
		return
	}
	defer db.Close()

	port := ":3002"
	cert := "cert.pem"
	key := "key.pem"

	fmt.Println("Starting server on port", port)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	//rl := mw.NewRateLimiter(5, time.Minute)

	//create custom server with TLS config
	router := routers.Router()
	mux := mw.Cors(router)
	server := &http.Server{
		Addr: port,
		// Handler: rl.RateLimiters(mw.Compression(mw.ResponseTime(mw.SecurityHeaders(mw.Cors(mux))))),
		Handler:   mux,
		TLSConfig: tlsConfig,
	}
	// Start the server with TLS
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		fmt.Println("Error starting server", err)
	}

}
