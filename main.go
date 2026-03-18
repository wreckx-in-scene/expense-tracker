package main

import (
	"fmt"
	"log"
	"net/http"

	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/handlers"
)

func main() {
	connString := "postgres://postgres:Amogh%40123@localhost:5432/expense-tracker?sslmode=disable"

	err := db.Connect(connString)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	//routes
	http.HandleFunc("/register", handlers.Register) //register route
	http.HandleFunc("/login", handlers.Login)       //login route

	fmt.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
