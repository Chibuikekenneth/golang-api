package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("a DATABASE_URL must be present in the environment to run.")
	}
	db := NewDB(url)

	todoMapper := &DBTodoMapper{db}

	r := mux.NewRouter()
	r.Handle("/todos", &AllTodos{todoMapper}).
		Methods(http.MethodGet)
	r.Handle("/todos", &CreateTodo{todoMapper}).
		Methods(http.MethodPost)
	r.Handle("/todos/{ID:[0-9]+}", &UpdateTodo{todoMapper}).
		Methods(http.MethodPatch)

	allowHeaders := handlers.AllowedHeaders([]string{
		"Content-Type",
		"Accept",
	})
	allowMethods := handlers.AllowedMethods([]string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
	})
	handler := handlers.CORS(allowHeaders, allowMethods)(r)

	log.Fatal(http.ListenAndServe(":8000", handler))
}
