package main

import (
	"net/http"
	"fmt"
	_ "github.com/lib/pq"
	"project-start/internal/db"
	"project-start/internal/handler"
)

func main() {
	connStr := "postgres://dev:dev@localhost:5432/app?sslmode=disable"
	db.Connect(connStr)
	db.Migrate()

	http.HandleFunc("/health", handler.Health)      // health check
	http.HandleFunc("/users", handler.CreateUser)   // добавление пользователя
	http.HandleFunc("/user", handler.GetUser) // теперь GET /user?id=1
	http.HandleFunc("/user/", handler.DeleteUser)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}