package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi-go/internal/models"
	"restapi-go/internal/repository/sqlconnect"
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
	dbName := "school"
	db, err := sqlconnect.ConnectDb(dbName)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	fmt.Println("Connected to database:", dbName)
	pathParts := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimSuffix(pathParts, "/")
	fmt.Println("ID:", id)

	if id == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers where 1=1"
		var args []interface{}
		if firstName != "" {
			query += " AND first_name = ?"
			args = append(args, firstName)
		}
		if lastName != "" {
			query += " AND last_name = ?"
			args = append(args, lastName)
		}
		rows, err := db.Query(query, args...)
		if err != nil { // Handle the error
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		teacherList := make([]models.Teacher, 0)
		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(w, "Database scan error", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}
		if err = rows.Err(); err != nil {
			http.Error(w, "Database rows error", http.StatusInternalServerError)
			return
		}
		//
		//teacherList := make([]models.Teacher, 0, len(teachers))
		// for _, teacher := range teachers {
		// 	if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
		// 		teacherList = append(teacherList, teacher)
		// 	}
		// }

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
		return
	}

	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error converting ID to integer:", err)
		return
	}
	// teacher, ok := teachers[newId]
	// if !ok {
	// 	http.Error(w, "Teacher not found", http.StatusNotFound)
	// 	return
	// }
	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", newId).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database query error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}
func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	// var newTeachers []models.Teacher
	// err := json.NewDecoder(r.Body).Decode(&newTeachers)
	// if err != nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	// addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	// for _, newTeacher := range newTeachers {
	// 	newTeacher.ID = nextId
	// 	teachers[nextId] = newTeacher
	// 	addedTeachers = append(addedTeachers, newTeacher)
	// 	nextId++
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// response := struct {
	// 	Status string           `json:"status"`
	// 	Count  int              `json:"count"`
	// 	Data   []models.Teacher `json:"data"`
	// }{
	// 	Status: "success",
	// 	Count:  len(addedTeachers),
	// 	Data:   addedTeachers,
	// }
	// json.NewEncoder(w).Encode(response)
	// fmt.Println("Added new teachers:", addedTeachers)
	dbName := "school"
	db, err := sqlconnect.ConnectDb(dbName)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("Connected to database:", dbName)

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	smt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?,?, ?, ?)")
	if err != nil {
		http.Error(w, "Database insertion failed", http.StatusInternalServerError)
		return
	}
	defer smt.Close()

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, newTeacher := range newTeachers {
		res, err := smt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Prepared statement execution failed", http.StatusInternalServerError)
			return
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting last id", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastId)
		addedTeachers = append(addedTeachers, newTeacher)
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
