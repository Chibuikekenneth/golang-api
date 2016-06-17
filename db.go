package main

import "database/sql"

// Todo is the domain object representing todos
type Todo struct {
	ID   int
	Text string
}

// TodoMapper is an interface describing how to read Todos from a datastore
type TodoMapper interface {
	GetAllTodos() ([]Todo, error)
	AddTodo(todo NewTodo) (Todo, error)
}

// DBTodoMapper is an implementation of TodoMapper using a SQL database
type DBTodoMapper struct {
	db *sql.DB
}

// GetAllTodos returns all the todos in the database
func (t *DBTodoMapper) GetAllTodos() ([]Todo, error) {
	rows, err := t.db.Query(`
		SELECT id, text FROM todo
	`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var allTodos []Todo
	for rows.Next() {
		var id int
		var text string
		err = rows.Scan(&id, &text)
		if err != nil {
			return nil, err
		}

		allTodos = append(allTodos, Todo{ID: id, Text: text})
	}

	return allTodos, nil
}

// AddTodo takes a NewTodo object, inserts it into the DB, and returns a Todo
func (t *DBTodoMapper) AddTodo(todo NewTodo) (Todo, error) {
	var id int
	var text string
	err := t.db.QueryRow(`
		INSERT INTO todo(text)
		VALUES($1)
		RETURNING id, text
	`, todo.Text).Scan(&id, &text)
	if err != nil {
		return Todo{}, err
	}

	return Todo{ID: id, Text: text}, nil
}
