package telegram

import (
	"context"
	"strings"

	roleModel "github.com/Antesser/trainingFinder/internal/model/roles"
	telegramModel "github.com/Antesser/trainingFinder/internal/model/telegram"
)

func (bt *Bot) HandleState(ctx context.Context, chatID int64, text string) (string, error) {
	s, err := bt.cache.LoadSession(ctx, chatID)
	if err != nil {
		return "", err
	}

	switch s.State {
	case telegramModel.StateIdle:
		s.State = telegramModel.StateAwaitingLogin
		return "Введите логин:", bt.cache.SaveSession(ctx, chatID, s)

	case telegramModel.StateAwaitingLogin:
		s.Login = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingPass
		return "Теперь введите пароль:", bt.cache.SaveSession(ctx, chatID, s)

	case telegramModel.StateAwaitingPass:
		s.Password = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingRole
		return "Введите роль coach или user:", bt.cache.SaveSession(ctx, chatID, s)

	case telegramModel.StateAwaitingRole:
		role := strings.TrimSpace(text)
		roleID, ok := roleModel.GetRoleID(role)
		if !ok {
			return "Неизвестная роль, попробуйте ещё раз:", bt.cache.SaveSession(ctx, chatID, s)
		}
		s.Role = roleID

		if _, err := bt.service.SignUp(ctx, s.Login, s.Password, s.Role); err != nil {
			return "Ошибка авторизации: " + err.Error(), nil
		}
		if err := bt.cache.DelSession(ctx, chatID); err != nil {
			return "", err
		}
		return "Авторизация успешна!", nil

	default:
		_ = bt.cache.DelSession(ctx, chatID)
		return "Начнём заново. Введите логин:", nil
	}
}
