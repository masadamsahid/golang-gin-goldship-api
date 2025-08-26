package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func ConnectDB() {

	dsn := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=%s`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Println("Failed to establish connection to DB", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Println(err)
	}

	log.Println("Database connection successful.")
}

func StopDB() {
	DB.Close()
	log.Println("Success closing connection to DB")
}
