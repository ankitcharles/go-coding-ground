package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"restapi-go/internal/models"
	"restapi-go/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

func isValidSortOrder(order string) bool {
	order = strings.ToUpper(order)
	return order == "ASC" || order == "DESC"
}
func isValidSortField(field string) bool {
	validFileds := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"class":      true,
		"subject":    true,
		"email":      true,
	}
	return validFileds[field]
}
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
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
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers where 1=1"
		var args []interface{}

		query, args = addFilter(r, query, args)

		query = addSorting(r, query)

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

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"]
	if len(sortParams) > 0 {
		query += " ORDER BY "
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], strings.ToUpper(parts[1])
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ", "
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func addFilter(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}
func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

// PUT /teachers/{id}
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//updatedTeacher.ID = id

	dbName := "school"
	db, err := sqlconnect.ConnectDb(dbName)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("Connected to database:", dbName)

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class,subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database query error", http.StatusInternalServerError)
		}
		return
	}
	updatedTeacher.ID = existingTeacher.ID
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		http.Error(w, "Database update error while updating teacher", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
	fmt.Println("Updated teacher:", updatedTeacher)
}

// PATCH /teachers/{id}
func PatchSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//updatedTeacher.ID = id

	dbName := "school"
	db, err := sqlconnect.ConnectDb(dbName)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("Connected to database:", dbName)

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class,subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database query error", http.StatusInternalServerError)
		}
		return
	}

	// Use reflect to build the update query dynamically
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			jsonTag := field.Tag.Get("json")
			tagName := strings.Split(jsonTag, ",")[0]

			if tagName == k {
				fieldVal := teacherVal.Field(i)
				if fieldVal.CanSet() {
					val := reflect.ValueOf(v)
					if val.Type().ConvertibleTo(fieldVal.Type()) {
						fieldVal.Set(val.Convert(fieldVal.Type()))
					}
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)
	if err != nil {
		http.Error(w, "Database update error while updating teacher", http.StatusInternalServerError)
		return
	}
	// Update the teacher in the database
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
	fmt.Println("Updated teacher:", existingTeacher)
}

// patch /teachers
func PatchMultipleTeachersHandler(w http.ResponseWriter, r *http.Request) {
	dbName := "school"
	db, err := sqlconnect.ConnectDb(dbName)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var updates []map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println("Error starting transaction:", err)
		http.Error(w, "Database transaction error", http.StatusInternalServerError)
		return
	}
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			tx.Rollback()
			http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		}

		var teacherFromDB models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class,subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDB.ID, &teacherFromDB.FirstName, &teacherFromDB.LastName, &teacherFromDB.Email, &teacherFromDB.Class, &teacherFromDB.Subject)
		if err != nil {
			if err == sql.ErrNoRows {
				tx.Rollback()
				http.Error(w, "Teacher not found", http.StatusNotFound)
			} else {
				tx.Rollback()
				http.Error(w, "Database query error", http.StatusInternalServerError)
			}
		}
		//Apply updates using reflection in DB
		teacherVal := reflect.ValueOf(&teacherFromDB).Elem()
		teacherType := teacherVal.Type()
		for k, v := range update {
			if k == "id" {
				continue // skip the eky
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+" ,omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							tx.Rollback()
							http.Error(w, "Invalid value type", http.StatusBadRequest)
							return
						}
					}
					break
				}
			}
		}
		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDB.FirstName, teacherFromDB.LastName, teacherFromDB.Email, teacherFromDB.Class, teacherFromDB.Subject, id)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Database update error while updating teacher", http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction:", err)
		http.Error(w, "Database transaction error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	fmt.Println("Updated teachers:", updates)
}

// DELETE /teachers/{id}
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	dbName := "school"
	db, err := sqlconnect.ConnectDb(dbName)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("Connected to database:", dbName)
	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database delete error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error getting rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	//w.WriteHeader(http.StatusNoContent)
	fmt.Println("Deleted teacher with ID:", id)

	// Optionally, you can return a response body indicating successful deletion
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		Id     int
	}{
		Status: "success",
		Id:     id,
	}
	json.NewEncoder(w).Encode(response)
}
