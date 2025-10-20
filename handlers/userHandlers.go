package handlers

import (
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

type UserCreateResponse struct {
	UserID int64 `json:"user_id"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user UserCreateRequest
	if !tools.ValidateJSON(w, r, &user) {
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	result, err := db.DB.Exec("INSERT INTO users (username, email, password, user_level_id) VALUES (?, ?, ?, ?)", user.Username, user.Email, hashedPassword, 1)
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
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := UserCreateResponse{UserID: id}
	tools.JSONResponse(w, http.StatusCreated, response)
}

// GET /users/{id}
type UserResponse struct {
	UserID    int        `json:"user_id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	row := db.DB.QueryRow("SELECT user_id, username, email, created_at FROM users WHERE user_id = ?", id)

	var response UserResponse
	err = row.Scan(&response.UserID, &response.Username, &response.Email, &response.CreatedAt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	tools.JSONResponse(w, http.StatusOK, response)
}

// POST /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials LoginRequest
	if !tools.ValidateJSON(w, r, &credentials) {
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	var storedHashedPassword string
	var userID int64
	db.DB.QueryRow("SELECT password, user_id FROM users WHERE username = ?", credentials.Username).Scan(&storedHashedPassword, &userID)

	err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := tools.GenerateToken(userID, credentials.Username)
	if err != nil {
		http.Error(w, "could not generate token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{Token: token}
	tools.JSONResponse(w, http.StatusOK, response)
}
