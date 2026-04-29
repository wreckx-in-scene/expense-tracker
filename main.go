package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/handlers"
	"www.github.com/wreckx-in-scene/expense-tracker/middleware"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	connString := "postgres://postgres:Amogh%40123@localhost:5432/expense-tracker?sslmode=disable"

	err := db.Connect(connString)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	go func() {
		for {
			time.Sleep(3 * time.Minute)
			middleware.CleanupClients()
		}
	}()

	// auth routes
	http.HandleFunc("/register", middleware.Logger(middleware.RateLimit(handlers.Register)))
	http.HandleFunc("/login", middleware.Logger(middleware.RateLimit(handlers.Login)))

	// expense routes
	http.HandleFunc("/expenses", middleware.Logger(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateExpense(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetExpenses(w, r)
		}
	})))

	http.HandleFunc("/expenses/", middleware.Logger(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			handlers.UpdateExpense(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.DeleteExpense(w, r)
		}
	})))

	// income routes
	http.HandleFunc("/incomes", middleware.Logger(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateIncome(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetIncomes(w, r)
		}
	})))

	http.HandleFunc("/incomes/", middleware.Logger(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			handlers.UpdateIncome(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.DeleteIncome(w, r)
		}
	})))

	// analytics routes
	http.HandleFunc("/analytics/summary", middleware.Logger(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetSummary(w, r)
		}
	})))

	http.HandleFunc("/analytics/categories", middleware.Logger(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetCategory(w, r)
		}
	})))

	// ai routes
	http.HandleFunc("/ai/categorize", middleware.Logger(middleware.RateLimit(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.Categorize(w, r)
		}
	}))))

	http.HandleFunc("/ai/insights", middleware.Logger(middleware.RateLimit(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetInsights(w, r)
		}
	}))))

	http.HandleFunc("/ai/chat", middleware.Logger(middleware.RateLimit(middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.Chat(w, r)
		}
	}))))

	http.HandleFunc("/transactions/recent", middleware.Logger(middleware.Auth(handlers.GetRecentTransactions)))

	fmt.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(http.DefaultServeMux)))
}
