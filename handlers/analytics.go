package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"www.github.com/wreckx-in-scene/expense-tracker/db"
)

// summary struct
type SummaryResponse struct {
	TotalExpenses     float64 `json:"total_expenses"`
	TotalIncome       float64 `json:"total_income"`
	Savings           float64 `json:"savings"`
	ThisMonthExpenses float64 `json:"this_month_expenses"`
	ThisMonthIncome   float64 `json:"this_month_income"`
}

// category struct
type CategoryResponse struct {
	Category   string  `json:"category"`
	Total      float64 `json:"total"`
	Percentage float64 `json:"percentage"`
}

// handler to get summary of all expenses
func GetSummary(w http.ResponseWriter, r *http.Request) {
	//get user id from context
	userID := r.Context().Value("user_id").(int)

	var summary SummaryResponse
	err := db.DB.QueryRow(context.Background(), "SELECT COALESCE(SUM(amount) , 0) FROM expenses WHERE user_id = $1", userID).Scan(&summary.TotalExpenses)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user id"})
		return
	}

	err = db.DB.QueryRow(context.Background(), "SELECT monthly_income FROM users WHERE id = $1", userID).Scan(&summary.TotalIncome)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch income"})
		return
	}

	summary.Savings = summary.TotalIncome - summary.TotalExpenses

	//---this month's expenses
	err = db.DB.QueryRow(context.Background(), "SELECT COALESCE(SUM(amount) , 0) FROM expenses WHERE user_id = $1 AND DATE_TRUNC('month', date) = DATE_TRUNC('month' , CURRENT_DATE)", userID).Scan(&summary.ThisMonthExpenses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch this month expenses"})
		return
	}

	//---this month's income
	err = db.DB.QueryRow(context.Background(), "SELECT COALESCE(SUM(amount) , 0) FROM income WHERE user_id = $1 AND DATE_TRUNC('month' , date) = DATE_TRUNC('month' , CURRENT_DATE)", userID).Scan(&summary.ThisMonthIncome)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch the information"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(summary)

}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	//get user_id from context first
	userID := r.Context().Value("user_id").(int)

	//get total first for percentage calucation
	var total float64
	err := db.DB.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(amount::float8), 0) FROM expenses WHERE user_id = $1",
		userID).Scan(&total)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch total"})
		return
	}

	rows, err := db.DB.Query(context.Background(),
		"SELECT category, SUM(amount) as total FROM expenses WHERE user_id = $1 GROUP BY category", userID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch categories"})
		return
	}

	defer rows.Close()

	var categories []CategoryResponse
	for rows.Next() {
		var cat CategoryResponse
		rows.Scan(&cat.Category, &cat.Total)
		if total > 0 {
			cat.Percentage = (cat.Total / total) * 100
		}

		categories = append(categories, cat)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)

}
