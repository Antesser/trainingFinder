package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrSessionNotFound = errors.New("session not found")
var ErrSessionExpired = errors.New("session is expired")

type Session struct {
	Token     uuid.UUID
	UserID    string
	Active    bool
	CreatedAt time.Time
	ExpiresAt time.Time
	RoleID    string
}

func (s *Session) IsActive() bool {
	return s.Active && s.ExpiresAt.After(time.Now().UTC())
} // время всегда в ЮТЭСЭ, ВООБЩЕ ВСЕГДА, ЧЁ БЫ НЕ ПРОИСХОДИЛО, ФРОНТ ПРИМЕТ В ЮТЭСЭ, ПОН?
