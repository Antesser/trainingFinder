package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	roleModel "github.com/Antesser/trainingFinder/internal/model/roles"
	telegramModel "github.com/Antesser/trainingFinder/internal/model/telegram"
)

func (bt *Bot) HandleState(ctx context.Context, chatID int64, text string) (string, error) {
	s, err := bt.cache.LoadSession(ctx, chatID)
	if err != nil {
		return "", fmt.Errorf("load session for chat %d: %w", chatID, err)
	}

	switch s.State {
	case telegramModel.StateIdle:
		s.State = telegramModel.StateAwaitingLogin
		if err := bt.cache.SaveSession(ctx, chatID, s); err != nil {
			return "", fmt.Errorf("save session for chat %d: %w", chatID, err)
		}
		return "Введите логин:", nil

	case telegramModel.StateAwaitingLogin:
		s.Login = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingPass
		if err := bt.cache.SaveSession(ctx, chatID, s); err != nil {
			return "", fmt.Errorf("save session for chat %d: %w", chatID, err)
		}
		return "Теперь введите пароль:", nil

	case telegramModel.StateAwaitingPass:
		s.Password = strings.TrimSpace(text)
		s.State = telegramModel.StateAwaitingRole
		if err := bt.cache.SaveSession(ctx, chatID, s); err != nil {
			return "", fmt.Errorf("save session for chat %d: %w", chatID, err)
		}
		return "Введите роль coach или user:", nil

	case telegramModel.StateAwaitingRole:
		role := strings.TrimSpace(text)
		roleID, ok := roleModel.GetRoleID(role)
		if !ok {
			if err := bt.cache.SaveSession(ctx, chatID, s); err != nil {
				return "", fmt.Errorf("save session for chat %d: %w", chatID, err)
			}
			return "Неизвестная роль, попробуйте ещё раз:", nil
		}
		s.Role = roleID

		if _, err := bt.service.SignUp(ctx, s.Login, s.Password, s.Role); err != nil {
			log.Printf("sign up chat %d: %v", chatID, err)
			return "Ошибка регистрации, попробуйте позже", nil
		}
		if err := bt.cache.DelSession(ctx, chatID); err != nil {
			return "", fmt.Errorf("delete session for chat %d: %w", chatID, err)
		}
		return "Авторизация успешна!", nil

	default:
		if err := bt.cache.DelSession(ctx, chatID); err != nil {
			log.Printf("delete session for chat %d: %v", chatID, err)
		}
		return "Начнём заново. Введите логин:", nil
	}
}
