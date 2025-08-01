package token

import (
	"fmt"
	"time"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must have %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userId uuid.UUID, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayLoad(userId, duration)
	if err != nil {
		return "", nil, err
	}

	tkn, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}
	return tkn, payload, nil
}

func (maker *PasetoMaker) VerifyToken(tkn string) (*Payload, error) {
	var payload Payload
	err := maker.paseto.Decrypt(tkn, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
