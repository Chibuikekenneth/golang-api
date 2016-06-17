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

type ErrorMsg struct {
	Error string `json:"error"`
}

type SerializedTodo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

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

type NewTodo struct {
	Text string `json:"text"`
}
