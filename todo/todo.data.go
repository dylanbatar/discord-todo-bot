package todo

import (
	"database/sql"
	"fmt"

	"github.com/dylanbatar/github.com/database"
)

func GetTodosByUser(userId string) ([]Todo, error) {
	result, err := database.DbConn.Query("select id, name, description, complete from todos where userId = ?", userId)

	if err != nil {
		fmt.Println("Query error", err)
		return nil, err
	}

	defer result.Close()

	todos := []Todo{}

	for result.Next() {
		var todo Todo
		result.Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Complete)

		todos = append(todos, todo)
	}

	return todos, nil
}

func GetTodoByUser(userId string, todoId string) (*Todo, error) {
	result := database.DbConn.QueryRow("select id, name, description, complete from todos where userId = ? and id = ?", userId, todoId)

	var todo Todo

	err := result.Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Complete)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func CompleteTodo(userId string, todoId string) (*Todo, error) {
	_, err := database.DbConn.Exec("update todos set complete = true where userId = ? and id = ?", userId, todoId)

	if err != nil {
		return nil, err
	}

	result, err := GetTodoByUser(userId, todoId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateTodo(name, description, userId string) (int, error) {
	result, err := database.DbConn.Exec("insert into todos (name, description, complete, userId) values (?, ?, ?, ?)",
		name, description, false, userId)

	if err != nil {
		return 0, err
	}

	insertId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(insertId), nil
}
