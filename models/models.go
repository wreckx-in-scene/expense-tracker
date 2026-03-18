package models

import "time"

type User struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	Password      string    `json:"password" db:"password"`
	MonthlyIncome float64   `json:"monthly_income" db:"monthly_income"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Expense struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Date        time.Time `json:"date" db:"date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Income struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Source    string    `json:"source" db:"source"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Budget struct {
	ID          int     `json:"id" db:"id"`
	UserID      int     `json:"user_id" db:"user_id"`
	Category    string  `json:"category" db:"category"`
	LimitAmount float64 `json:"limit_amount" db:"limit_amount"`
}
