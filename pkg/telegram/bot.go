package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"trinity/pkg/code_gen"
)

var bot *tgbotapi.BotAPI

func Init() error {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		return errors.New("Укажи TG_BOT_TOKEN в переменных окружения")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = false
	log.Printf("Бот авторизован: @%s", bot.Self.UserName)
	return nil
}

// Обработка обновлений
func HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Обрабатываем только команды
		if !update.Message.IsCommand() {
			continue
		}

		chatID := update.Message.Chat.ID
		command := update.Message.Command()

		switch command {
		case "start":
			msg := tgbotapi.NewMessage(chatID, "Привет! Используй /generateCode для получения кода")
			bot.Send(msg)

		case "generateCode":
			// Генерируем код (можешь выбрать любой вариант)
			code, _ := code_gen.GenerateRandomNumber() // Например: A7K9M2X1
			// code := generateHexCode(4)  // Например: 3f7a9c2e

			responseText := fmt.Sprintf("🎫 Твой код: <code>%s</code>", code)

			msg := tgbotapi.NewMessage(chatID, responseText)
			msg.ParseMode = "HTML"

			if _, err := bot.Send(msg); err != nil {
				log.Printf("Ошибка отправки: %v", err)
			}

		default:
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда. Попробуй /generateCode")
			bot.Send(msg)
		}
	}
}

// Если нужно отправить код конкретному пользователю
func SendCodeToUser(chatID int64) error {
	code, _ := code_gen.GenerateRandomNumber()
	text := fmt.Sprintf("🎫 Твой код: <code>%s</code>\n\nИспользуй его в течение 24 часов", code)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil
}
