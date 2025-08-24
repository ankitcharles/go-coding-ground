package routers

import (
	"net/http"
	"restapi-go/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeacherHandler)
	mux.HandleFunc("GET /teachers/", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /teachers/", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeacherHandler)
	mux.HandleFunc("/students/", handlers.StudentsHandler)
	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
