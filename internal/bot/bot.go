package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Bot struct {
	API *tgbotapi.BotAPI
}

func NewBot() (*Bot, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Ошибка загрузки env")
	}
	token := os.Getenv("TELEGRAM_BOT_API")
	if token == "" {
		log.Fatal("Токен бота не найден")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = true
	log.Printf("Бот запущен", bot.Self.UserName)
	return &Bot{API: bot}, nil
}
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.API.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет выбери формат файла и загрузи их")
			b.API.Send(msg)
		}
	}
}
