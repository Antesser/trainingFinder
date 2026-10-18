package telegram

import (
	"context"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	service authService
	token   string
	cache   sessionCache
}

func New(token string, service authService, newCache sessionCache) (*Bot, error) {
	log.Printf("token len=%d, prefix=%q, suffix=%q", len(token), token[:min(10, len(token))], token[max(0, len(token)-5):])
	return &Bot{service: service, token: token, cache: newCache}, nil
}

func (bt *Bot) Run(ctx context.Context) error {
	api, err := bot.New(bt.token, bot.WithDefaultHandler(bt.defaultHandler))
	if err != nil {

		return err
	}
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
		err := bt.cache.DelSession(ctx, chatID)
		if err != nil {
			log.Printf("state error %v", err)
			return
		}
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
	//b.SendMessage(ctx, &bot.SendMessageParams{
	//	ChatID: chatID,
	//	Text:   "Введите логин",
	//})
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
