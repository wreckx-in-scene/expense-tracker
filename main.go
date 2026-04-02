package main

import (
	"fmt"
	"log"
	"net/http"

	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/handlers"
	"www.github.com/wreckx-in-scene/expense-tracker/middleware"
)

func main() {
	connString := "postgres://postgres:Amogh%40123@localhost:5432/expense-tracker?sslmode=disable"

	err := db.Connect(connString)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	//routes
	http.HandleFunc("/register", handlers.Register)                                             //register route
	http.HandleFunc("/login", handlers.Login)                                                   //login route
	http.HandleFunc("/expenses", middleware.Auth(func(w http.ResponseWriter, r *http.Request) { //create expense route
		if r.Method == http.MethodPost {
			handlers.CreateExpense(w, r)
		} else if r.Method == http.MethodGet { //get expense route
			handlers.GetExpenses(w, r)
		}
	}))

	http.HandleFunc("/expenses/", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch { //update expense route
			handlers.UpdateExpense(w, r)
		} else if r.Method == http.MethodDelete { //delete expense route
			handlers.DeleteExpense(w, r)
		}
	}))

	//analytics routes
	http.HandleFunc("/analytics/summary", middleware.Auth(handlers.GetSummary))
	http.HandleFunc("/analytics/categories", middleware.Auth(handlers.GetCategory))

	fmt.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
