package token

import "time"

type Maker interface {
	CreateToken(userId string, duration time.Duration) (string, error)
	VerifyToken(tkn string) (*Payload, error)
}
