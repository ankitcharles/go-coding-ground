package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"restapi-go/internal/models"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("Http Post method on Execs route only"))
		fmt.Println("Method Post called on Execs route")
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			fmt.Println("Error parsing form:", err)
			return
		}
		fmt.Println("Form data:", r.Form)
		// prepare a response
		response := make(map[string]interface{})
		for key, values := range r.Form {
			response[key] = values[0]
		}
		fmt.Println("Response:", response)
		// raw data read
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading body", http.StatusBadRequest)
			fmt.Println("Error reading body:", err)
			return
		}
		defer r.Body.Close()
		fmt.Println("Raw data:", string(body))

		var user models.User
		// parse the body as JSON
		err = json.Unmarshal(body, &user)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			fmt.Println("Error parsing JSON:", err)
			return
		}
		fmt.Println("Parsed User:", user)
	case http.MethodPut:
		w.Write([]byte("Http Put method on Teachers route only"))
		fmt.Println("Method Put called on Teachers route")
	case http.MethodDelete:
		w.Write([]byte("Http Delete method on Teachers route only"))
		fmt.Println("Method Delete called on Teachers route")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
