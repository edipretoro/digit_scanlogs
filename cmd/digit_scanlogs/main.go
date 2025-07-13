package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dbURL := os.Getenv("DIGIT_SCAN_DSN")
	if dbURL == "" {
		log.Fatal("DIGIT_SCAN_DSN environment variable is not set")
	}
	db, err := sql.Open("sqlite", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	scanDirectory := os.Getenv("DIGIT_SCAN_DIR")
	if scanDirectory == "" {
		log.Fatal("DIGIT_SCAN_DIR environment variable is not set")
	}
	log.Printf("Scanning directory: %s", scanDirectory)
	err = processingScanDirectory(scanDirectory)
	if err != nil {
		log.Fatalf("Error processing scan directory: %v", err)
	} else {
		log.Println("Scan directory processed successfully")

	}
}
