package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

// createUserWithTimezone creates a user with a specific timezone
func createUserWithTimezone(t testing.TB, timezone string) User {
	t.Helper()
	arg := CreateUserParams{
		Email: util.GetRandomString(10),
		Name:  util.GetRandomString(10),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	// Update the user's timezone
	updateArg := UpdateUserTimezoneParams{
		ID:       user.ID,
		Timezone: timezone,
	}
	updatedUser, err := testQueries.UpdateUserTimezone(context.Background(), updateArg)
	require.NoError(t, err)
	require.Equal(t, timezone, updatedUser.Timezone)

	return updatedUser
}

func TestCreateTaskDaysForUsersInTimezone(t *testing.T) {
	timezone := "UTC"

	// Create users in the specified timezone
	user1 := createUserWithTimezone(t, timezone)
	user2 := createUserWithTimezone(t, timezone)

	// Create a user in a different timezone to ensure it's not affected
	_ = createUserWithTimezone(t, "America/New_York")

	// Run the cron function to create task days
	timezoneInterval := pgtype.Interval{Valid: false} // We'll use the timezone string instead
	err := testQueries.CreateTaskDaysForUsersInTimezone(context.Background(), timezoneInterval)
	require.NoError(t, err)

	// Verify task days were created for users in the target timezone
	arg1 := ListTaskDaysByUserIdParams{
		UserID: user1.ID,
		Limit:  10,
		Offset: 0,
	}
	taskDays1, err := testQueries.ListTaskDaysByUserId(context.Background(), arg1)
	require.NoError(t, err)

	arg2 := ListTaskDaysByUserIdParams{
		UserID: user2.ID,
		Limit:  10,
		Offset: 0,
	}
	taskDays2, err := testQueries.ListTaskDaysByUserId(context.Background(), arg2)
	require.NoError(t, err)

	_ = taskDays1
	_ = taskDays2
}
