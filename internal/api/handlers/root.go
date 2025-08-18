package handlers

import (
	"fmt"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "Hello World")
	w.Write([]byte("Hello Root Route"))
	fmt.Println("Root route accessed")
}
