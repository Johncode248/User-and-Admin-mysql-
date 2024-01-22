package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var(
    db *sql.DB
)

func OpenDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "asek1:pass12@tcp(db-project:3306)/bigproject")

	if err != nil {
		return nil, err
	}

	// Warto dodać test połączenia, aby upewnić się, że połączenie jest prawidłowe.
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTableIfNotExists() error {

	var err error
	db, err = OpenDatabaseConnection()
	if err != nil {
		return err
	}

	// Utworzenie tabeli
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS project_table (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(45),
		surname VARCHAR(45),
		date_birth VARCHAR(100),
		email VARCHAR(100),
		password VARCHAR(120),
		updated_at VARCHAR(100)
	);
`)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	fmt.Println("Table created successfully or already exists.")
	return nil
}
