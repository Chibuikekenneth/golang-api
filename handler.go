package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json; charset=utf-8"
)

// ErrorMsg represents a JSON error body
type ErrorMsg struct {
	Error string `json:"error"`
}

// NewTodo is a JSON struct to intake new todo requests
type NewTodo struct {
	Text *string `json:"text"`
}

// IsInvalid checks if the NewTodo is invalid
func (n *NewTodo) IsInvalid() bool {
	return n.Text == nil
}

// EditTodoRequest represents a user's request to toggle the
// completed property of a todo
type EditTodoRequest struct {
	Completed *bool `json:"completed"`
}

// IsInvalid checks the EditTodoRequest for validity and returns false if
// Completed is nil
func (e *EditTodoRequest) IsInvalid() bool {
	return e.Completed == nil
}

// SerializedTodo is the serializable JSON shape of a Todo
type SerializedTodo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// AllTodos returns all todos in the datastore
type AllTodos struct {
	mapper TodoMapper
}

func (a *AllTodos) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(contentType, applicationJSON)

	encoder := json.NewEncoder(rw)
	todos, err := a.mapper.GetAllTodos()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{fmt.Sprintf("%s", err)})
		return
	}

	var serializedTodos []SerializedTodo
	for _, t := range todos {
		serializedTodos = append(serializedTodos, SerializedTodo{
			ID:        t.ID,
			Text:      t.Text,
			Completed: t.Completed,
		})
	}

	encoder.Encode(serializedTodos)
}

// CreateTodo is a handler for creating new todos
type CreateTodo struct {
	mapper TodoMapper
}

func (c *CreateTodo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(contentType, applicationJSON)

	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(r.Body)

	var newTodo NewTodo
	err := decoder.Decode(&newTodo)
	if err != nil || newTodo.IsInvalid() {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{"invalid request"})
		return
	}

	todo, err := c.mapper.AddTodo(newTodo)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{fmt.Sprintf("%s", err)})
		return
	}

	rw.WriteHeader(http.StatusCreated)
	encoder.Encode(SerializedTodo{
		ID:        todo.ID,
		Text:      todo.Text,
		Completed: todo.Completed,
	})
}

// UpdateTodo is the handler for updating todos in the database
type UpdateTodo struct {
	mapper TodoMapper
}

func (u *UpdateTodo) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set(contentType, applicationJSON)

	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(r.Body)

	vars := mux.Vars(r)
	todoIDStr := vars["ID"]

	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{"invalid request"})
		return
	}

	var editReq EditTodoRequest
	err = decoder.Decode(&editReq)
	if err != nil || editReq.IsInvalid() {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{"invalid request"})
		return
	}

	updatedTodo, err := u.mapper.UpdateTodo(EditedTodo{
		ID:        todoID,
		Completed: *editReq.Completed,
	})
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		encoder.Encode(&ErrorMsg{fmt.Sprintf("%s", err)})
		return
	}

	encoder.Encode(SerializedTodo{
		ID:        updatedTodo.ID,
		Text:      updatedTodo.Text,
		Completed: updatedTodo.Completed,
	})
}
