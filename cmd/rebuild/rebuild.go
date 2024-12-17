package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3" // Blank import to initialize sqlite3 driver
)

func main() {
	statusCode := 1
	defer func() {
		os.Exit(statusCode)
	}()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		fmt.Println("Error opening database", err)
	}

	deleteQuery := `DROP TABLE customers; DROP TABLE orders; DROP TABLE products;`
	_, err = db.Exec(deleteQuery)
	if err != nil {
		fmt.Println("Error running query", err)
	}

	statusCode = 0
}
