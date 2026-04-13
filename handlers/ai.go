package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/utils"
)

type CategorizeReq struct {
	Description string `json:"description"`
}

func Categorize(w http.ResponseWriter, r *http.Request) {
	var req CategorizeReq
	json.NewDecoder(r.Body).Decode(&req)

	if req.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "description is required"})
		return
	}

	prompt := "Categorize this expense: '" + req.Description + "'.Reply with just the category name , nothing else. Categories: Food , Transport , Entertainment , Shopping , Health , Bills , Other"

	category, err := utils.CallGemini(prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not categorize expense"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"category": category})
}

func GetInsights(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	rows, err := db.DB.Query(context.Background(),
		"SELECT category , amount , description FROM expenses WHERE user_id = $1", userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch expenses"})
		return
	}

	defer rows.Close()

	//building a string summary of all expenses
	expenseList := ""
	for rows.Next() {
		var category, description string
		var amount float64
		rows.Scan(&category, &amount, &description)
		expenseList += category + " - " + description + " - ₹" + fmt.Sprintf("%.2f", amount) + "\n"
	}

	prompt := "You are a financial advisor . Analyse these expenses and give 3 short insights:\n" + expenseList

	insights, err := utils.CallGemini(prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not get insights"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"insights": insights})
}

func Chat(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	// read user question
	type ChatReq struct {
		Message string `json:"message"`
	}
	var req ChatReq
	json.NewDecoder(r.Body).Decode(&req)
	if req.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "message is required"})
		return
	}

	// fetch expenses
	rows, err := db.DB.Query(context.Background(),
		"SELECT category, amount, description FROM expenses WHERE user_id = $1", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch expenses"})
		return
	}
	defer rows.Close()

	// build expense list
	expenseList := ""
	for rows.Next() {
		var category, description string
		var amount float64
		rows.Scan(&category, &amount, &description)
		expenseList += category + " - " + description + " - ₹" + fmt.Sprintf("%.2f", amount) + "\n"
	}

	// fetch totals
	var totalIncome float64
	var totalExpenses float64

	db.DB.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(amount), 0) FROM income WHERE user_id = $1", userID).Scan(&totalIncome)

	db.DB.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(amount::float8), 0) FROM expenses WHERE user_id = $1", userID).Scan(&totalExpenses)

	totalSavings := totalIncome - totalExpenses

	// build prompt
	prompt := fmt.Sprintf(`You are a personal finance assistant. Here is the user's financial data:

Expenses:
%s

Total Income: ₹%.2f
Total Expenses: ₹%.2f
Total Savings: ₹%.2f

User question: %s

Answer concisely based on the data above.`, expenseList, totalIncome, totalExpenses, totalSavings, req.Message)

	response, err := utils.CallGemini(prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not get response"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}
