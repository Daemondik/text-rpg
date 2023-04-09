package main

import (
	"github.com/joho/godotenv"
	"log"
	"wer/cmd/game/api/telegram"
)

func main() {
	tg := telegram_api.Api{}
	tg.Start()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
