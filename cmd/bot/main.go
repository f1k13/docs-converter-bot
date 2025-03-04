package main

import (
	"github.com/f1k13/go-file-converter-bot/internal/bot"
	"log"
)

func main() {
	b, err := bot.NewBot()
	if err != nil {
		log.Fatal(err)
	}
	b.Start()
}
