package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {

	password := GetRandomString(6)
	hashPassword, err := HashedPassword(password)
	require.NoError(t, err)

	err = CheckPassword(password, hashPassword)
	require.NoError(t, err)

	wrongPassword := GetRandomString(6)
	err = CheckPassword(wrongPassword, hashPassword)
	require.Error(t, err)
}
