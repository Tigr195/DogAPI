package main

import (
	"Dogs/internal/app"
	"flag"
	"github.com/joho/godotenv"
	"log"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	count := flag.Int("l", 1, "number of dogs")
	flag.Parse()

	if err := app.New().Run(*count); err != nil {
		log.Fatalf("Critical error: %v", err)
	}
}
