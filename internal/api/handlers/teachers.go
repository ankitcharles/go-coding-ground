package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restapi-go/internal/models"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "10",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "11",
		Subject:   "Science",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "12",
		Subject:   "English",
	}
	nextId++
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeacherHandler(w, r)
		fmt.Println("Method Get called on Teachers route")
	case http.MethodPost:
		addTeacherHandler(w, r)
		fmt.Println("Method Post called on Teachers route")
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

func getTeacherHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimSuffix(pathParts, "/")
	fmt.Println("ID:", id)
	if id == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teacherList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
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
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, newTeacher := range newTeachers {
		newTeacher.ID = nextId
		teachers[nextId] = newTeacher
		addedTeachers = append(addedTeachers, newTeacher)
		nextId++
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
	fmt.Println("Added new teachers:", addedTeachers)

}
