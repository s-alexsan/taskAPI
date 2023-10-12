package main

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type task struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var tasks []task

func validUser(c *gin.Context) bool {

	authHeader := c.GetHeader("Authorization")

	splitHeader := strings.SplitN(authHeader, " ", 2)

	if len(splitHeader) != 2 || splitHeader[0] != "Basic" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Header Authorization"})
		return false
	}

	if splitHeader[1] == "" {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Empty authorization"})
		return false
	}

	decodedCreds, err := base64.StdEncoding.DecodeString(splitHeader[1])
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": authHeader})
		return false
	}

	creds := strings.Split(string(decodedCreds), ":")

	if len(creds) != 2 {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credencial"})
		return false
	}

	if creds[0] == "admin" && creds[1] == "secretpassword" {
		return true
	}

	c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "User invalid"})
	return false

}

func getTasks(c *gin.Context) {

	valid := validUser(c)

	if valid {
		c.IndentedJSON(http.StatusOK, tasks)
	}
}

func putTask(c *gin.Context) {
	if validUser(c) {
		var putTask task
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		}

		for index, a := range tasks {
			if a.Id == id {
				if err := c.BindJSON(&putTask); err != nil {
					return
				} else {
					tasks[index] = putTask
					c.IndentedJSON(http.StatusOK, putTask)
				}
			}
		}
	}
}

func getTask(c *gin.Context) {
	if validUser(c) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		}

		for _, a := range tasks {
			if a.Id == id {
				c.IndentedJSON(http.StatusOK, a)
				return
			}
		}

		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
	}
}

func postTasks(c *gin.Context) {
	if validUser(c) {
		var newTask task

		if err := c.BindJSON(&newTask); err != nil {
			return
		}

		tasks = append(tasks, newTask)
		c.IndentedJSON(http.StatusCreated, newTask)
	}
}

func deleteTask(c *gin.Context) {
	if validUser(c) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)

		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		}

		for index, a := range tasks {
			if a.Id == id {
				tasks = append(tasks[:index], tasks[index+1:]...)
				c.IndentedJSON(http.StatusOK, gin.H{"message": "task " + string(id) + " deleted"})
			}
		}
	}
}

func main() {
	router := gin.Default()

	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTask)
	router.POST("/tasks", postTasks)
	router.PUT("tasks/:id", putTask)
	router.DELETE("tasks/:id", deleteTask)

	router.Run("localhost:8080")
}
