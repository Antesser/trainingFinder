package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrSessionNotFound error = fmt.Errorf("session not found")
var ErrSessionExpired error = fmt.Errorf("session is expired")

type Session struct {
	Token     uuid.UUID
	UserID    string
	Active    bool
	CreatedAt time.Time
	ExpiresAt time.Time
}

func (s *Session) IsActive() bool {
	return s.Active && s.ExpiresAt.After(time.Now().UTC())
} // время всегда в ЮТЭСЭ, ВООБЩЕ ВСЕГДА, ЧЁ БЫ НЕ ПРОИСХОДИЛО, ФРОНТ ПРИМЕТ В ЮТЭСЭ, ПОН?
