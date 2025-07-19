package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

// createRandomTask creates a new task with random data
// It first creates a task day to link the task to
func createRandomTask(t testing.TB) Task {
	t.Helper()
	taskDay := createRandomTaskDay(t)

	arg := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       util.GetRandomString(10),
		Description: pgtype.Text{String: util.GetRandomString(20), Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(time.Hour * 2), Valid: true}, // 2 hours
	}

	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, task)
	require.Equal(t, arg.TaskDayID, task.TaskDayID)
	require.Equal(t, arg.Title, task.Title)
	require.Equal(t, arg.Description, task.Description)
	require.Equal(t, arg.Duration, task.Duration)
	require.NotEmpty(t, task.ID)
	require.False(t, task.Completed.Bool) // Should default to false

	return task
}

func TestCreateTask(t *testing.T) {
	createRandomTask(t)
}

func TestGetTask(t *testing.T) {
	task1 := createRandomTask(t)
	task2, err := testQueries.GetTask(context.Background(), task1.ID)
	require.NoError(t, err)
	require.Equal(t, task1, task2)
}

func TestDeleteTask(t *testing.T) {
	task1 := createRandomTask(t)
	task2, err := testQueries.DeleteTask(context.Background(), task1.ID)
	require.NoError(t, err)
	require.Equal(t, task1, task2)

	// Verify the task was deleted
	task3, err := testQueries.GetTask(context.Background(), task1.ID)
	require.Error(t, err)
	require.Empty(t, task3)
	require.Equal(t, err, pgx.ErrNoRows)
}

func TestUpdateTask(t *testing.T) {
	task := createRandomTask(t)

	arg := UpdateTaskParams{
		ID:          task.ID,
		Title:       util.GetRandomString(15),
		Description: pgtype.Text{String: util.GetRandomString(30), Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(time.Hour * 3), Valid: true}, // 3 hours
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}

	_, err := testQueries.UpdateTask(context.Background(), arg)
	require.NoError(t, err)

	// Verify the task was updated
	updatedTask, err := testQueries.GetTask(context.Background(), task.ID)
	require.NoError(t, err)
	require.Equal(t, arg.Title, updatedTask.Title)
	require.Equal(t, arg.Description, updatedTask.Description)
	require.Equal(t, arg.Duration, updatedTask.Duration)
	require.Equal(t, arg.Completed, updatedTask.Completed)
	require.Equal(t, task.ID, updatedTask.ID)
	require.Equal(t, task.TaskDayID, updatedTask.TaskDayID)
}

func TestListTasks(t *testing.T) {
	// Create multiple tasks
	count := 10
	for i := 0; i < count; i++ {
		createRandomTask(t)
	}

	arg := ListTasksParams{
		Limit:  5,
		Offset: 2,
	}
	tasks, err := testQueries.ListTasks(context.Background(), arg)
	require.NoError(t, err)
	require.LessOrEqual(t, len(tasks), int(arg.Limit))

	for _, task := range tasks {
		require.NotEmpty(t, task)
		require.NotEmpty(t, task.ID)
		require.NotEmpty(t, task.TaskDayID)
		require.NotEmpty(t, task.Title)
	}
}

func TestListTasksEmpty(t *testing.T) {
	// Test with large offset to get empty results
	arg := ListTasksParams{
		Limit:  5,
		Offset: 1000,
	}
	tasks, err := testQueries.ListTasks(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, tasks)
}

func TestUpdateTaskCompleted(t *testing.T) {
	task := createRandomTask(t)

	// Update only the completed status
	arg := UpdateTaskParams{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Duration:    task.Duration,
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}

	_, err := testQueries.UpdateTask(context.Background(), arg)
	require.NoError(t, err)

	// Verify only completed status changed
	updatedTask, err := testQueries.GetTask(context.Background(), task.ID)
	require.NoError(t, err)
	require.True(t, updatedTask.Completed.Bool)
	require.True(t, updatedTask.Completed.Valid)
	require.Equal(t, task.Title, updatedTask.Title)
	require.Equal(t, task.Description, updatedTask.Description)
	require.Equal(t, task.Duration, updatedTask.Duration)
}
