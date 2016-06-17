package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// NewDB returns a new database connection
func NewDB(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	url := "postgres://hhgffzzd:aZ28VW-KxEUfRuzxzPbut7HeREh3DTNu@pellefant.db.elephantsql.com:5432/hhgffzzd"
	db := NewDB(url)

	todoMapper := &DBTodoMapper{db}

	r := mux.NewRouter()
	r.Handle("/todos", &AllTodos{todoMapper}).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
