package main

import "database/sql"

// Todo is the domain object representing todos
type Todo struct {
	ID        int
	Text      string
	Completed bool
}

// EditedTodo represents a todo object with a new completed state
type EditedTodo struct {
	ID        int
	Completed bool
}

// TodoMapper is an interface describing how to read Todos from a datastore
type TodoMapper interface {
	GetAllTodos() ([]Todo, error)
	AddTodo(todo NewTodo) (*Todo, error)
	UpdateTodo(todo EditedTodo) (*Todo, error)
}

// DBTodoMapper is an implementation of TodoMapper using a SQL database
type DBTodoMapper struct {
	db *sql.DB
}

// GetAllTodos returns all the todos in the database
func (t *DBTodoMapper) GetAllTodos() ([]Todo, error) {
	rows, err := t.db.Query(`
		SELECT id, text, completed FROM todo
	`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var allTodos []Todo
	for rows.Next() {
		var id int
		var text string
		var completed bool
		err = rows.Scan(&id, &text, &completed)
		if err != nil {
			return nil, err
		}

		allTodos = append(allTodos, Todo{ID: id, Text: text, Completed: completed})
	}

	return allTodos, nil
}

// AddTodo takes a NewTodo object, inserts it into the DB, and returns a Todo
func (t *DBTodoMapper) AddTodo(todo NewTodo) (*Todo, error) {
	var id int
	var text string
	var completed bool
	err := t.db.QueryRow(`
		INSERT INTO todo(text)
		VALUES($1)
		RETURNING id, text, completed
	`, todo.Text).Scan(&id, &text, &completed)
	if err != nil {
		return nil, err
	}

	return &Todo{ID: id, Text: text, Completed: completed}, nil
}

func (t *DBTodoMapper) UpdateTodo(editedTodo EditedTodo) (*Todo, error) {
	var id int
	var text string
	var completed bool
	err := t.db.QueryRow(`
		UPDATE todo SET completed = $1 WHERE id = $2
		RETURNING id, text, completed
	`, editedTodo.Completed, editedTodo.ID).Scan(&id, &text, &completed)
	if err != nil {
		return nil, err
	}
	return &Todo{ID: id, Completed: completed, Text: text}, nil
}
