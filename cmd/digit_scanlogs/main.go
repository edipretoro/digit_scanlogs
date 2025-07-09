package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
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
