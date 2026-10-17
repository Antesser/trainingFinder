package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Antesser/trainingFinder/internal/redis"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	service authService
	token   string
}

func New(token string, service authService) (*Bot, error) {
	log.Printf("token len=%d, prefix=%q, suffix=%q", len(token), token[:min(10, len(token))], token[max(0, len(token)-5):])
	return &Bot{service: service, token: token}, nil
}

func (bt *Bot) Run(ctx context.Context) error {
	api, err := bot.New(bt.token, bot.WithDefaultHandler(bt.defaultHandler))
	if err != nil {

		fmt.Println("err", err)
		return err
	}
	fmt.Println("api", api)
	api.Start(ctx)

	return nil
}

func (bt *Bot) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		return
	}

	chatID := update.Message.Chat.ID

	if text == "/reset" {
		_ = resetSession(ctx, chatID)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Сессия сброшена. Введите логин:",
		})
		reply, _ := bt.HandleState(ctx, chatID, "")
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   reply,
		})
		return
	}

	reply, err := bt.HandleState(ctx, chatID, text)
	if err != nil {
		log.Printf("state error for %d: %v", chatID, err)
		reply = "Внутренняя ошибка, попробуйте позже."
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   reply,
	})
}

func resetSession(ctx context.Context, chatID int64) error {
	return redis.Rdb.Del(ctx, redis.SessionKey(chatID)).Err()
}
