package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/handlers"
	"www.github.com/wreckx-in-scene/expense-tracker/middleware"
)

func main() {
	godotenv.Load()
	connString := "postgres://postgres:Amogh%40123@localhost:5432/expense-tracker?sslmode=disable"

	err := db.Connect(connString)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	//routes
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)

	http.HandleFunc("/expenses", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateExpense(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetExpenses(w, r)
		}
	}))

	http.HandleFunc("/expenses/", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			handlers.UpdateExpense(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.DeleteExpense(w, r)
		}
	}))

	//analytics routes
	http.HandleFunc("/analytics/summary", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetSummary(w, r)
		}
	}))

	http.HandleFunc("/analytics/categories", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetCategory(w, r)
		}
	}))

	//ai routes
	http.HandleFunc("/ai/categorize", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.Categorize(w, r)
		}
	}))

	http.HandleFunc("/ai/insights", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetInsights(w, r)
		}
	}))

	http.HandleFunc("/ai/chat", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.Chat(w, r)
		}
	}))

	fmt.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
