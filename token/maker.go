package token

import (
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(userId uuid.UUID, duration time.Duration) (string, *Payload, error)
	VerifyToken(tkn string) (*Payload, error)
}
