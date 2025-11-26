package main

import (
	"altoai_mvp/logic"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using environment variables")
	}
}

func main() {
	logic.LogicGpt()
}	
