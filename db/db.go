package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init(dsn string) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 10; i++ {
		if err = DB.Ping(); err == nil {
			log.Println("Database connection established")
			return
		}
		log.Printf("DB not ready (attempt %d/10): %v", i, err)
		time.Sleep(1 * time.Second)
	}
	log.Fatal("Database unreachable:", err)

	log.Println("Database connection established")
}
