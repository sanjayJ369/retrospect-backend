package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

// createRandomTaskDay creates a new user and then
// creates a task day for that user
func createRandomTaskDay(t testing.TB) TaskDay {
	t.Helper()
	user := createRandomUser(t)

	taskDay, err := testQueries.CreateTaskDay(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, taskDay)
	require.Equal(t, user.ID, taskDay.UserID)
	require.NotEmpty(t, taskDay.ID)

	return taskDay
}

func TestCreateTaskDay(t *testing.T) {
	createRandomTaskDay(t)
}

func TestGetTaskDay(t *testing.T) {
	taskDay1 := createRandomTaskDay(t)
	taskDay2, err := testQueries.GetTaskDay(context.Background(), taskDay1.ID)
	require.NoError(t, err)
	require.Equal(t, taskDay1, taskDay2)
}

func TestDeleteTaskDay(t *testing.T) {
	taskDay1 := createRandomTaskDay(t)
	taskDay2, err := testQueries.DeleteTaskDay(context.Background(), taskDay1.ID)
	require.NoError(t, err)
	require.Equal(t, taskDay1, taskDay2)

	// Verify the task day was deleted
	taskDay3, err := testQueries.GetTaskDay(context.Background(), taskDay1.ID)
	require.Error(t, err)
	require.Empty(t, taskDay3)
	require.Equal(t, err, pgx.ErrNoRows)
}

func TestListTaskDays(t *testing.T) {
	// Create multiple task days
	count := 5
	for i := 0; i < count; i++ {
		createRandomTaskDay(t)
	}

	arg := ListTaskDaysParams{
		Limit:  3,
		Offset: 1,
	}
	taskDays, err := testQueries.ListTaskDays(context.Background(), arg)
	require.NoError(t, err)
	require.LessOrEqual(t, len(taskDays), int(arg.Limit))

	for _, taskDay := range taskDays {
		require.NotEmpty(t, taskDay)
		require.NotEmpty(t, taskDay.ID)
		require.NotEmpty(t, taskDay.UserID)
	}
}

func TestListTaskDaysEmpty(t *testing.T) {
	// Test with large offset to get empty results
	arg := ListTaskDaysParams{
		Limit:  5,
		Offset: 1000,
	}
	taskDays, err := testQueries.ListTaskDays(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, taskDays)
}

func TestGetTaskDayByDateAndUser(t *testing.T) {
	// Create a task day
	taskDay1 := createRandomTaskDay(t)

	// Get the task day by date and user_id
	arg := GetTaskDayByDateAndUserParams{
		Date:   taskDay1.Date,
		UserID: taskDay1.UserID,
	}
	taskDay2, err := testQueries.GetTaskDayByDateAndUser(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, taskDay1, taskDay2)
}

func TestGetTaskDayByDateAndUserNotFound(t *testing.T) {
	user := createRandomUser(t)

	// Try to get a task day that doesn't exist
	arg := GetTaskDayByDateAndUserParams{
		Date:   pgtype.Date{Time: time.Now().AddDate(0, 0, -30), Valid: true}, // 30 days ago
		UserID: user.ID,
	}
	taskDay, err := testQueries.GetTaskDayByDateAndUser(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, taskDay)
	require.Equal(t, err, pgx.ErrNoRows)
}

func TestGetTaskDayByDateAndUserDifferentUser(t *testing.T) {
	// Create a task day for user1
	taskDay1 := createRandomTaskDay(t)

	// Create a different user
	user2 := createRandomUser(t)

	// Try to get the task day with user2's ID (should not find it)
	arg := GetTaskDayByDateAndUserParams{
		Date:   taskDay1.Date,
		UserID: user2.ID,
	}
	taskDay, err := testQueries.GetTaskDayByDateAndUser(context.Background(), arg)
	require.Error(t, err)
	require.Empty(t, taskDay)
	require.Equal(t, err, pgx.ErrNoRows)
}
