package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t testing.TB) User {
	arg := CreateUserParams{
		Email: util.GetRandomString(10),
		Name:  util.GetRandomString(10),
	}

	user1, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user1)
	require.Equal(t, arg.Email, user1.Email)
	require.Equal(t, arg.Name, user1.Name)

	return user1
}

func TestUserCreation(t *testing.T) {
	createRandomUser(t)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
	require.Equal(t, user1, user2)

	user3, err := testQueries.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.Empty(t, user3)
	require.Equal(t, err, pgx.ErrNoRows)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.Equal(t, user1, user2)
}

func TestListUsers(t *testing.T) {
	// TODO: update test for more through testing
	arg := ListUsersParams{
		Limit:  2,
		Offset: 2,
	}
	res, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int(arg.Limit), len(res))
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)
	arg := UpdateUserEmailParams{
		ID:    user.ID,
		Email: util.GetRandomString(10),
	}
	res, err := testQueries.UpdateUserEmail(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, arg.Email, res.Email)

	user1, err := testQueries.GetUser(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, arg.Email, user1.Email)
}

func TestUpdateUserName(t *testing.T) {
	user := createRandomUser(t)
	arg := UpdateUserNameParams{
		ID:   user.ID,
		Name: util.GetRandomString(10),
	}

	user1, err := testQueries.UpdateUserName(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, user1.Name, arg.Name)
}

func TestUpdateUserTimeZone(t *testing.T) {
	user := createRandomUser(t)
	arg := UpdateUserTimezoneParams{
		ID:       user.ID,
		Timezone: "GMT",
	}
	user1, err := testQueries.UpdateUserTimezone(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Timezone, user1.Timezone)
}
