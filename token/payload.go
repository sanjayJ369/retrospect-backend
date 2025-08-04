package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken      = errors.New("token is invalid")
	PurposeLogin         = "login"
	PurposeVerifyEmail   = "verify_email"
	PurposeResetPassword = "reset_password"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Purpose   string    `json:"purpose"`
}

func NewPayLoad(userId uuid.UUID, duration time.Duration, purpose string) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := Payload{
		ID:        tokenId,
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Purpose:   purpose,
	}

	return &payload, nil
}

// GetExpirationTime implements the Claims interface
func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

// GetIssuedAt implements the Claims interface
func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

// GetNotBefore implements the Claims interface
func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

// GetIssuer implements the Claims interface
func (p *Payload) GetIssuer() (string, error) {
	return "", nil
}

// GetSubject implements the Claims interface
func (p *Payload) GetSubject() (string, error) {
	return p.UserId.String(), nil
}

// GetAudience implements the Claims interface
func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

// Valid checks if the token is valid and not expired
func (p *Payload) Valid() error {
	if time.Until(p.ExpiredAt) <= 0 {
		return ErrInvalidToken
	}
	return nil
}
