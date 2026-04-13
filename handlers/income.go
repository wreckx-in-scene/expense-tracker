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

func CreateIncome(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	type IncomeReq struct {
		Amount float64 `json:"amount"`
		Source string  `json:"source"`
		Date   string  `json:"date"`
	}

	var req IncomeReq
	json.NewDecoder(r.Body).Decode(&req)

	if req.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "amount is required"})
		return
	}

	_, err := db.DB.Exec(context.Background(),
		"INSERT INTO income (user_id, amount, source, date) VALUES ($1, $2, $3, $4)",
		userID, req.Amount, req.Source, req.Date)
	if err != nil {
		fmt.Println("Error is :", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not save income"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Income created successfully"})
}

func GetIncomes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var incomes []models.Income

	rows, err := db.DB.Query(context.Background(),
		"SELECT id, user_id, amount, source, date FROM income WHERE user_id = $1", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch incomes"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var inc models.Income
		rows.Scan(&inc.ID, &inc.UserID, &inc.Amount, &inc.Source, &inc.Date)
		incomes = append(incomes, inc)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incomes)
}

func UpdateIncome(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid income id"})
		return
	}

	type UpdateIncomeReq struct {
		Amount float64 `json:"amount"`
		Source string  `json:"source"`
		Date   string  `json:"date"`
	}

	var req UpdateIncomeReq
	json.NewDecoder(r.Body).Decode(&req)

	if req.Amount == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "amount is required"})
		return
	}

	_, err = db.DB.Exec(context.Background(),
		"UPDATE income SET amount=$1, source=$2, date=$3 WHERE id=$4 AND user_id=$5",
		req.Amount, req.Source, req.Date, id, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not update income"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Income updated successfully"})
}

func DeleteIncome(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid income id"})
		return
	}

	_, err = db.DB.Exec(context.Background(),
		"DELETE FROM income WHERE id=$1 AND user_id=$2", id, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not delete income"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Income deleted successfully"})
}
