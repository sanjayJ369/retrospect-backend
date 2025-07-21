package token

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

func TestPasteoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.GetRandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	user, err := uuid.NewUUID()
	require.NoError(t, err)

	issuedAt := time.Now()
	expiredAt := time.Now().Add(time.Minute)

	token, err := maker.CreateToken(user, time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, user, payload.UserId)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Minute)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.GetRandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	user, err := uuid.NewUUID()
	require.NoError(t, err)
	token, err := maker.CreateToken(user, -time.Minute)
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Empty(t, payload)
}
