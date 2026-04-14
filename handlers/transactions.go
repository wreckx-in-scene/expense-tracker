package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"www.github.com/wreckx-in-scene/expense-tracker/db"
)

type Transaction struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}

func GetRecentTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	query := `SELET 'expense' as type , amount::float64 , description , category , date::text FROM expenses WHERE user_id = $1 UNION ALL
	SELECT 'income' as type , amount::float8 , source as description , 'Income' as category, date::text FROM income WHERE user_id = $1
	ORDER BY date DESC LIMIT 10`

	rows, err := db.DB.Query(context.Background(), query, userID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "could not fetch transaction"})
		return
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var t Transaction
		rows.Scan(&t.Type, &t.Amount, &t.Description, &t.Category, &t.Date)
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
