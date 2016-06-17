package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json; charset=utf-8"
)

// ErrorMsg represents a JSON error body
type ErrorMsg struct {
	Error string `json:"error"`
}

// SerializedTodo is the serializable JSON shape of a Todo
type SerializedTodo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// AllTodos returns all todos in the datastore
type AllTodos struct {
	mapper TodoMapper
}

func (a *AllTodos) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	todos, err := a.mapper.GetAllTodos()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{fmt.Sprintf("%s", err)})
		return
	}

	rw.Header().Set(contentType, applicationJSON)

	var serializedTodos []SerializedTodo
	for _, t := range todos {
		serializedTodos = append(serializedTodos, SerializedTodo{
			ID:   t.ID,
			Text: t.Text,
		})
	}

	encoder.Encode(serializedTodos)
}

// NewTodo is a JSON struct to intake new todo requests
type NewTodo struct {
	Text string `json:"text"`
}
