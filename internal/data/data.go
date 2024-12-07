package data

import (
	"database/sql"
	_ "embed"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3" // Blank import to initialize sqlite3 driver
)

type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"customer_id"`
	Date       string `json:"date"`
}

type Dataset struct {
	Customers []Customer `json:"customers"`
	Orders    []Order    `json:"orders"`
}

//go:embed dataset.json
var dataset []byte

func getStartingData() (*Dataset, error) {
	var result Dataset
	if err := json.Unmarshal(dataset, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func CreateDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		return nil, err
	}
	// defer db.Close() // TODO: has to move

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS customers (
		id INTEGER PRIMARY KEY,
		name TEXT
	); CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY,
		customer_id INTEGER,
		order_date TEXT
	); DELETE FROM customers; DELETE FROM orders;`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	startingDataset, err := getStartingData()
	if err != nil {
		return nil, err
	}

	insertStmt := `INSERT INTO customers (id, name) VALUES (?, ?)`
	for _, customer := range startingDataset.Customers {
		_, err := db.Exec(insertStmt, customer.ID, customer.Name)
		if err != nil {
			return nil, err
		}
	}

	insertOrdersStmt := "INSERT INTO orders (id, customer_id, order_date) VALUES (?, ?, ?)"
	for _, order := range startingDataset.Orders {
		_, err := db.Exec(insertOrdersStmt, order.ID, order.CustomerID, order.Date)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}