package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/dylanbatar/github.com/handlers"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading env %s", err)
	}

	fmt.Println(os.Getenv("BOT_TOKEN"))
}

func main() {
	fmt.Println("Bot running")
}
