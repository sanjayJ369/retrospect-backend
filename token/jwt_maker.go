package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	minSecretKeySize = 32
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be atleast %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey: secretKey}, nil
}

func (j *JWTMaker) CreateToken(userId uuid.UUID, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayLoad(userId, duration)
	if err != nil {
		return "", nil, err
	}

	jwttkn := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tkn, err := jwttkn.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tkn, payload, nil
}

func (j *JWTMaker) VerifyToken(tkn string) (*Payload, error) {
	keyfunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}
	token, err := jwt.ParseWithClaims(tkn, &Payload{}, keyfunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
