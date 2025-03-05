package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) sendFormatFilesSelection(chatID int64, formats []string) {

	var buttons [][]tgbotapi.KeyboardButton
	for _, format := range formats {
		buttons = append(buttons, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(format)))
	}
	msg := tgbotapi.NewMessage(chatID, "Выберите формат, в который хотите конвертировать:")
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(buttons...)
	b.API.Send(msg)
}
func validateFormat(text string, formats []string) bool {
	for _, format := range formats {
		if text == format {
			return true
		}
	}
	return false
}
func (b *Bot) addFiles(chatID int64, fileID string, fileName string, session *UserSession) {

	session.Files = append(session.Files, fileID)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Файл %s загружен", fileName))
	button := tgbotapi.NewKeyboardButton("Конвертировать")
	replyMarkup := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(button))
	msg.ReplyMarkup = replyMarkup
	b.API.Send(msg)
}
func (b *Bot) uploadConverterFiles(chatID int64, session *UserSession) {
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Конвертирую %d файлов d %s...", len(session.Files)))
	b.API.Send(msg)
}
