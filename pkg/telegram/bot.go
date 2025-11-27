package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"time"
	"trinity/internal/models"
	"trinity/internal/repository_redis"
	"trinity/pkg/code_gen"
)

type Bot struct {
	*tgbotapi.BotAPI
	db  models.Model
	rdb repository_redis.RedisDB
}

func NewBot(db models.Model, redisDB repository_redis.RedisDB) *Bot {
	return &Bot{
		db:  db,
		rdb: redisDB,
	}
}

func (bot *Bot) ConnectBot() error {
	token := os.Getenv("TG_BOT_TOKEN")
	if token == "" {
		return errors.New("Укажи TG_BOT_TOKEN в переменных окружения")
	}
	var err error
	bot.BotAPI, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = false
	log.Printf("Бот авторизован: @%s", bot.Self.UserName)
	return nil
}

// Обработка обновлений
func (bot *Bot) HandleUpdates() {
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
			code, _ := code_gen.GenerateRandomNumber()
			username := update.Message.From.UserName

			user, err := bot.db.GetUserByTelegram(username)
			if err == nil {
				_ = bot.rdb.SetData("tg_"+user.Login, []byte(code), time.Minute*10)
			}
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
