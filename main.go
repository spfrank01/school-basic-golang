package main

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func getTodosHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	stmt, err := db.Prepare("SELECT id, title, status FROM todos;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var todos []Todo
	for rows.Next() {
		t2 := Todo{}

		if err := rows.Scan(&t2.ID, &t2.Title, &t2.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, t2)
	}
	c.JSON(http.StatusOK, todos)
}

func main() {
	r := gin.Default()

	r.GET("/api/todos", getTodosHandler)

	err := r.Run(":1234")
	if err != nil {
		error.Error(err)
	}
}
