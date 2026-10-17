package telegram

import (
	"context"
	"strings"

	telegramModel "github.com/Antesser/trainingFinder/internal/model/telegram"
	"github.com/Antesser/trainingFinder/internal/redis"
	red "github.com/Antesser/trainingFinder/internal/redis"
)

func (bt *Bot) HandleState(ctx context.Context, chatID int64, text string) (string, error) {
	s, err := redis.LoadSession(ctx, chatID)
	if err != nil {
		return "", err
	}

	switch s.State {
	case telegramModel.StateIdle:
		s.Login = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingLogin
		s.State = telegramModel.StateAwaitingPass
		return "Теперь введите пароль:", red.SaveSession(ctx, chatID, s)

	case telegramModel.StateAwaitingLogin:
		s.Login = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingPass
		return "Теперь введите пароль:", red.SaveSession(ctx, chatID, s)

	case telegramModel.StateAwaitingPass:
		s.Password = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingRole
		return "Введите роль (например, admin или user):", red.SaveSession(ctx, chatID, s)

	case telegramModel.StateAwaitingRole:
		s.Role = strings.TrimSpace(text)
		if _, err := bt.service.SignUp(ctx, s.Login, s.Password, s.Role); err != nil {
			return "Ошибка авторизации: " + err.Error(), nil
		}
		if err := resetSession(ctx, chatID); err != nil {
			return "", err
		}
		return "Авторизация успешна!", nil

	default:
		_ = resetSession(ctx, chatID)
		return "Начнём заново. Введите логин:", nil
	}
}
