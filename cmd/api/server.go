package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	mw "restapi-go/internal/api/middlewares"
	"strconv"
	"strings"
	"sync"
)

type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

type Teacher struct {
	ID        int
	FirstName string
	LastName  string
	Class     string
	Subject   string
}

var (
	teachers = make(map[int]Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "10",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "11",
		Subject:   "Science",
	}
	nextId++
	teachers[nextId] = Teacher{
		ID:        nextId,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "12",
		Subject:   "English",
	}
}

func getTeacherHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimSuffix(pathParts, "/")
	fmt.Println("ID:", id)
	if id == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teacherList := make([]Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting ID to integer:", err)
		return
	}
	teacher, ok := teachers[newId]
	if !ok {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(teacher)
}
func addTeacherHandler(w http.ResponseWriter, r *http.Request) {

}

func rootHandler(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintf(w, "Hello World")
	w.Write([]byte("Hello Root Route"))
	fmt.Println("Root route accessed")
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//w.Write([]byte("Http Get method on Teachers route only"))
		getTeacherHandler(w, r)
		fmt.Println("Method Get called on Teachers route")
	case http.MethodPost:
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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func execsHandler(w http.ResponseWriter, r *http.Request) {
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
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func main() {
	port := ":3002"
	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/teachers/", teachersHandler)
	mux.HandleFunc("/students/", studentsHandler)
	mux.HandleFunc("/execs/", execsHandler)
	fmt.Println("Starting server on port", port)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	//rl := mw.NewRateLimiter(5, time.Minute)

	//create custom server with TLS config
	server := &http.Server{
		Addr: port,
		// Handler: rl.RateLimiters(mw.Compression(mw.ResponseTime(mw.SecurityHeaders(mw.Cors(mux))))),
		Handler:   mw.Cors(mux),
		TLSConfig: tlsConfig,
	}
	// Start the server with TLS
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		fmt.Println("Error starting server", err)
	}

}
