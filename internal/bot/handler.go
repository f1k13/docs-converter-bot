package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	for _, fileID := range session.Files {
		fileURL, err := b.downloadFile(fileID)
		if err != nil {
			b.API.Send(tgbotapi.NewMessage(chatID, "Ошибка загрузки файла"))
			continue
		}
		convertedFile, err := convertFile(userSessions[chatID].SelectedFormat, fileURL)
		if err != nil {
			b.API.Send(tgbotapi.NewMessage(chatID, "Ошибка конвертации"))
			continue
		}
		b.sendConvertedFiles(chatID, convertedFile)
	}
	delete(userSessions, chatID)
	
}
func (b *Bot) downloadFile(fileID string) (string, error) {
	file, err := b.API.GetFile(tgbotapi.FileConfig{fileID})
	if err != nil {
		return "", err
	}
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", b.API.Token, file.FilePath)
	log.Println("скачивание файла  по ссылке:", fileURL)
	res, err := http.Get(fileURL)
	if err != nil {
		log.Println("Ошибка скачивания файла:", err)
		return "", err
	}
	defer res.Body.Close()
	localPath := "./downloads" + filepath.Base(file.FilePath)
	outFile, err := os.Create(localPath)
	if err != nil {
		log.Println("Ошибка создания файла", err)
		return "", err
	}
	_, err = io.Copy(outFile, res.Body)
	if err != nil {
		log.Println("Ошибка записи файла", err)
		return "", err
	}
	return localPath, nil

}
func convertFile(format, fileURL string) (string, error) {
	switch format {
	case "PDF":
		return convertToPDF(fileURL)
	default:
		return "", fmt.Errorf("Неизвестный формат: %s", format)
	}
}

func convertToPDF(inputFile string) (string, error) {
	outputFile := strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".pdf"
	cmd := exec.Command("unoconv", "-f", "pdf", "-0", outputFile, inputFile)
	err := cmd.Run()
	if err != nil {
		return "", nil
	}
	return outputFile, nil
}
func (b *Bot) sendConvertedFiles(chatID int64, filePath string) {
	file := tgbotapi.NewDocumentUpload(chatID, filePath)
	b.API.Send(file)
}
