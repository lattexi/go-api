package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-api/db"
	"go-api/tools"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// POST /users
type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID    int        `json:"user_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	_, err = db.DB.Exec("INSERT INTO users (username, email, password, user_level_id) VALUES (?, ?, ?, ?)", user.Username, user.Email, hashedPassword, 1)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				http.Error(w, "username or email already exists", http.StatusConflict)
				return
			}
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Accepted")
	w.WriteHeader(http.StatusAccepted)
}

// GET /users/{id}
func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow("SELECT user_id, username, email, created_at FROM users WHERE user_id = ?", id)

	var user UserResponse
	err = row.Scan(&user.UserID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// POST /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials LoginRequest
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	var storedHashedPassword string
	var userID int64
	err = db.DB.QueryRow("SELECT password, user_id FROM users WHERE username = ?", credentials.Username).Scan(&storedHashedPassword, &userID)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := tools.GenerateToken(userID, credentials.Username)

	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	fmt.Println("Login successful, token generated")

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"token": "%s"}`, token)
}
