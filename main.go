package main

import (
	"fmt"
	"go-api/db"
	"go-api/handlers"
	"go-api/middleware"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	mux := http.NewServeMux()

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", user, pass, host, port, name)
	db.Init(dsn)

	mux.Handle("GET /users/{id}", middleware.Chain(
		http.HandlerFunc(handlers.GetUser),
	))
	mux.Handle("POST /users",
		middleware.Chain(
			http.HandlerFunc(handlers.CreateUser),
			middleware.TestMiddleware,
		),
	)

	// Login reitti
	mux.Handle("POST /login", middleware.Chain(
		http.HandlerFunc(handlers.LoginUser),
	))

	// Media item reitit
	mux.Handle("GET /mediaitems", middleware.Chain(
		http.HandlerFunc(handlers.GetMediaItems),
	))
	mux.Handle("GET /mediaitems/bytitle", middleware.Chain(
		http.HandlerFunc(handlers.GetMediaItemsByTitle),
	))
	mux.Handle("GET /files/{filename}", middleware.Chain(
		http.HandlerFunc(handlers.ServeFile),
	))
	mux.Handle("POST /mediaitems", middleware.Chain(
		http.HandlerFunc(handlers.CreateMediaItem),
		middleware.Auth,
	))
	mux.Handle("DELETE /mediaitems", middleware.Chain(
		http.HandlerFunc(handlers.DeleteMediaItem),
		middleware.Auth,
	))
	mux.Handle("PUT /mediaitems", middleware.Chain(
		http.HandlerFunc(handlers.UpdateMediaItem),
		middleware.Auth,
	))

	http.ListenAndServe(":8080", mux)
}
