package main

import (
	"database/sql"
	"fmt"
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

func getTodosByIdHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	stmt, err := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1;")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	row := stmt.QueryRow(c.Param("id"))
	t := Todo{}
	if err := row.Scan(&t.ID, &t.Title, &t.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
	}

	c.JSON(http.StatusOK, t)
}

func postTodosHandler(c *gin.Context) {
	t := Todo{}
	if err := c.ShouldBind(&t); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	stmt, err := db.Prepare("INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if err := stmt.QueryRow(t.Title, t.Status).Scan(&t.ID); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, t)
}

func deleteTodosByIdHandler(c *gin.Context) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	stmt, err := db.Prepare("DELECT FROM todos WHERE id=$1;")
	if err != nil {
		fmt.Println("err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	if val, err := stmt.Query(c.Param("id")); err != nil {
		fmt.Println("val", val)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
	}

	c.JSON(http.StatusOK, "removed")
}

func main() {
	r := gin.Default()

	r.GET("/api/todos", getTodosHandler)
	r.GET("/api/todos/:id", getTodosByIdHandler)
	r.POST("/api/todos", postTodosHandler)
	r.DELETE("/api/todos/:id", deleteTodosByIdHandler)

	err := r.Run(":1234")
	if err != nil {
		error.Error(err)
	}
}
