package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDatabase menginisialisasi koneksi database MySQL
func InitDatabase() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Format DSN untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	DB = db
	log.Println("Database connected successfully!")
	return nil
}

// CloseDatabase menutup koneksi database
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
