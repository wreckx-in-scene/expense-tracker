package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	fmt.Println("REQ : ", req) //debug log

	//validating
	if req.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "amount is required"})
		return
	}

	//save to databse if req is valid
	_, err := db.DB.Exec(context.Background(), "INSERT INTO expenses (user_id , amount , description , category , date) VALUES ($1 , $2 , $3 , $4 , $5)", userID, req.Amount, req.Description, req.Category, req.Date)

	if err != nil {
		fmt.Println("DB error: ", err)
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
	defer rows.Close()

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

// update expense route
func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user id"})
		return
	}

	type UpdateReq struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Date        string  `json:"date"`
	}

	var req UpdateReq
	json.NewDecoder(r.Body).Decode(&req)

	//validating
	if req.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Amount is required"})
		return
	}

	//query
	_, err = db.DB.Exec(context.Background(), "UPDATE expenses SET amount = $1 , description = $2 , category = $3 , date = $4 WHERE id = $5 AND user_id = $6", req.Amount, req.Description, req.Category, req.Date, id, userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not update expense"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense updated successfully"})
}

// handler function for delete Expense
func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	//getting user id from context
	userID := r.Context().Value("user_id").(int)

	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user id"})
		return
	}

	//setting up query
	_, err = db.DB.Exec(context.Background(), "DELETE FROM expenses WHERE id = $1 AND user_id = $2", id, userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not delete expense"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Expense Deleted"})
}
