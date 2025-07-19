package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
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
