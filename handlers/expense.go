package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/models"
)

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	//fetch the user from context
	userID := r.Context().Value("user_id").(int)

	type ExpenseReq struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Date        string  `json:"date"`
	}

	//reading the json body
	var req ExpenseReq
	json.NewDecoder(r.Body).Decode(&req)

	//validating
	if req.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "amount is required"})
		return
	}

	//save to databse if req is valid
	_, err := db.DB.Exec(context.Background(), "INSERT INTO expenses (user_id , amount , description , category , date) VALUES ($1 , $2 , $3 , $4 , $5)", userID, req.Amount, req.Description, req.Category, req.Date)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not save expense"})
		return
	}

	//success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense created successfully"})
}

func GetExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	//using expense model
	var expenses []models.Expense

	rows, err := db.DB.Query(context.Background(), "SELECT id , amount , description , category , date FROM expenses WHERE user_id = $1", userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch expenses"})
		return
	}

	defer rows.Close()

	//looping through rows

	for rows.Next() {
		var exp models.Expense
		rows.Scan(&exp.ID, &exp.Amount, &exp.Description, &exp.Category, &exp.Date)
		expenses = append(expenses, exp)
	}

	json.NewEncoder(w).Encode(expenses)

}
