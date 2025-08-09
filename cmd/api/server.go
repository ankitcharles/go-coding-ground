package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello World")
	w.Write([]byte("Hello Root Route"))
	fmt.Println("Root route accessed")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodGet:
		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		userID := strings.TrimSuffix(path, "/")
		fmt.Println(userID)

		fmt.Println("Querry Parameters:", r.URL.Query())
		queryParamas := r.URL.Query()
		sortBy := queryParamas.Get("sortby")
		sortOrder := queryParamas.Get("sortorder")
		fmt.Printf("Sort By: %s, Sort Order: %s\n", sortBy, sortOrder)

		w.Write([]byte("Http Get method on Teachers route only"))
		fmt.Println("Method Get called on Teachers route")

	case http.MethodPost:
		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		userID := strings.TrimSuffix(path, "/")
		fmt.Println(userID)
		fmt.Println("Querry Parameters:", r.URL.Query())
		w.Write([]byte("Http Post method on Teachers route only"))
		fmt.Println("Method Post called on Teachers route")
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

		var user User
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
	}
	// if r.Method == http.MethodGet {
	// 	w.Write([]byte("Http Get method on Teachers route only"))
	// 	fmt.Println("Method Get called on Teachers route")
	// }
	w.Write([]byte("Hello Teachers Route"))
	fmt.Println("Teachers route accessed")
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case http.MethodPost:
		w.Write([]byte("Http Post method on Students route only"))
		fmt.Println("Method Post called on Students route")
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

		var user User
		// parse the body as JSON
		err = json.Unmarshal(body, &user)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			fmt.Println("Error parsing JSON:", err)
			return
		}
		fmt.Println("Parsed User:", user)
	case http.MethodPut:
		w.Write([]byte("Http Put method on Students route only"))
		fmt.Println("Method Put called on Students route")
	case http.MethodDelete:
		w.Write([]byte("Http Delete method on Students route only"))
		fmt.Println("Method Delete called on Students route")
	}
	// if r.Method == http.MethodGet {
	// 	w.Write([]byte("Http Get method on Teachers route only"))
	// 	fmt.Println("Method Get called on Teachers route")
	// }
	w.Write([]byte("Hello Students Route"))
	fmt.Println("Students route accessed")
}
func execsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
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

		var user User
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
	}
	// if r.Method == http.MethodGet {
	// 	w.Write([]byte("Http Get method on Teachers route only"))
	// 	fmt.Println("Method Get called on Teachers route")
	// }
	w.Write([]byte("Hello Execs Route"))
	fmt.Println("Excses route accessed")
}
func main() {
	port := ":3002"
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/teachers/", teachersHandler)
	http.HandleFunc("/students/", studentsHandler)
	http.HandleFunc("/execs/", execsHandler)
	fmt.Println("Starting server on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server", err)
	}

}
