package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"www.github.com/wreckx-in-scene/expense-tracker/db"
	"www.github.com/wreckx-in-scene/expense-tracker/models"
	"www.github.com/wreckx-in-scene/expense-tracker/utils"
)

type RegisterReq struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	MonthlyIncome float64 `json:"monthly_income"`
}

// login req struct
type LoginReq struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
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

// login function
func Login(w http.ResponseWriter, r *http.Request) {
	//reading json data from body and validating
	var Loginreq LoginReq
	json.NewDecoder(r.Body).Decode(&Loginreq)
	if Loginreq.Email == "" || Loginreq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email and Password are required"})
		return
	}

	//finding user from database
	var user models.User
	err := db.DB.QueryRow(context.Background(),
		"SELECT id , email , password FROM users WHERE email = $1",
		Loginreq.Email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if err.Error() == "no rows in result set" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong"})
		}
		return
	}

	//password comparision
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Loginreq.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not generate token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}
