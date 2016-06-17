package main

import "database/sql"

type Todo struct {
	ID   int
	Text string
}

type TodoMapper interface {
	GetAllTodos() ([]Todo, error)
	AddTodo(todo NewTodo) (Todo, error)
}

type DBTodoMapper struct {
	db *sql.DB
}

func (t *DBTodoMapper) GetAllTodos() ([]Todo, error) {
	rows, err := t.db.Query(`
		SELECT id, text FROM todo
	`)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	allTodos := make([]Todo, 0)
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
