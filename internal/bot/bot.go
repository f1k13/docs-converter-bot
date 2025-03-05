package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Bot struct {
	API *tgbotapi.BotAPI
}
type UserSession struct {
	SelectedFormat string
	Files          []string
}

var userSessions = make(map[int64]*UserSession)

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
	formats := []string{"PDF", "DOCX", "TXT", "JPG", "PNG", "MD"}
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
		chatID := update.Message.Chat.ID
		text := update.Message.Text
		if text == "/start" {
			b.sendFormatFilesSelection(chatID, formats)
		}
		if validateFormat(text, formats) {
			userSessions[chatID] = &UserSession{SelectedFormat: text}
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Вы выбрали %s. Теперь загрузите файлы", text))
			b.API.Send(msg)
			continue
		}
		fileID := update.Message.Document.FileID
		fileName := update.Message.Document.FileName
		if update.Message.Document != nil {
			session, exist := userSessions[chatID]
			if !exist || session.SelectedFormat == "" {
				msg := tgbotapi.NewMessage(chatID, "Сначала выберите формат")
				b.API.Send(msg)
				continue
			}
			b.addFiles(chatID, fileID, fileName, session)
			continue
		}
		if text == "конвертировать" {
			session, exist := userSessions[chatID]
			if !exist || len(session.Files) == 0 {
				msg := tgbotapi.NewMessage(chatID, "Вы не загрузили ни одного файла")
				b.API.Send(msg)
				continue
			}
			b.uploadConverterFiles(chatID, session)

			delete(userSessions, chatID)
			continue
		}
	}
}
