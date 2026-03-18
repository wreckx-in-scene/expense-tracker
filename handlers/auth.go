package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"www.github.com/wreckx-in-scene/expense-tracker/db"
)

type RegisterReq struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	MonthlyIncome float64 `json:"monthly_income"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//reading data and validating
	var req RegisterReq
	json.NewDecoder(r.Body).Decode(&req)
	if req.Email == "" || req.Password == "" || req.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "name , email and password are required"})

		return
	}

	//hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Something went wrong"})

		return
	}

	//saving to database
	_, err = db.DB.Exec(context.Background(),
		"INSERT INTO users(name , email , password , monthly_income) VALUES ($1 , $2 , $3 , $4)",
		req.Name, req.Email, string(hashedPassword), req.MonthlyIncome,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not create user"})
		return
	}

	//all success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
