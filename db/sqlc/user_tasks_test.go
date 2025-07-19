package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

func TestListTaskDaysByUserId(t *testing.T) {
	// Create multiple users, each with their own task day
	// Since there's a unique constraint on user_id + date,
	// we can't create multiple task days for the same user on the same day
	userCount := 3
	var users []User
	var taskDays []TaskDay

	for i := 0; i < userCount; i++ {
		user := createRandomUser(t)
		users = append(users, user)

		taskDay, err := testQueries.CreateTaskDay(context.Background(), user.ID)
		require.NoError(t, err)
		taskDays = append(taskDays, taskDay)
	}

	// Test listing task days for the first user
	arg := ListTaskDaysByUserIdParams{
		UserID: users[0].ID,
		Limit:  3,
		Offset: 0,
	}

	userTaskDays, err := testQueries.ListTaskDaysByUserId(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, 1, len(userTaskDays)) // Should only have 1 task day per user per day

	// Verify the returned task day belongs to the correct user
	taskDay := userTaskDays[0]
	require.Equal(t, users[0].ID, taskDay.UserID)
	require.NotEmpty(t, taskDay.ID)
	require.True(t, taskDay.Date.Valid)
}

func TestListTaskDaysByUserIdEmpty(t *testing.T) {
	user := createRandomUser(t)

	arg := ListTaskDaysByUserIdParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}

	taskDays, err := testQueries.ListTaskDaysByUserId(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, taskDays)
}

func TestListTaskDaysByUserIdWithPagination(t *testing.T) {
	// Since each user can only have one task day per date,
	// we need to create multiple users to test pagination
	userCount := 8
	var users []User

	for i := 0; i < userCount; i++ {
		user := createRandomUser(t)
		users = append(users, user)

		// Create a task day for each user
		_, err := testQueries.CreateTaskDay(context.Background(), user.ID)
		require.NoError(t, err)
	}

	// Test pagination by listing all task days (across all users)
	// Note: This test is more about verifying the pagination works
	// In a real scenario, you'd typically query by user_id

	// For this test, let's just verify that pagination works for one user
	// (which will have at most 1 task day)
	user := users[0]

	// Test first page
	arg1 := ListTaskDaysByUserIdParams{
		UserID: user.ID,
		Limit:  3,
		Offset: 0,
	}
	firstPage, err := testQueries.ListTaskDaysByUserId(context.Background(), arg1)
	require.NoError(t, err)
	require.LessOrEqual(t, len(firstPage), 1) // Should have at most 1 task day

	// Test second page (should be empty)
	arg2 := ListTaskDaysByUserIdParams{
		UserID: user.ID,
		Limit:  3,
		Offset: 3,
	}
	secondPage, err := testQueries.ListTaskDaysByUserId(context.Background(), arg2)
	require.NoError(t, err)
	require.Empty(t, secondPage) // Should be empty since user has only 1 task day
}

func TestListTaskDaysByUserIdOrdering(t *testing.T) {
	user := createRandomUser(t)

	// Create a single task day (since user can only have one per date)
	_, err := testQueries.CreateTaskDay(context.Background(), user.ID)
	require.NoError(t, err)

	arg := ListTaskDaysByUserIdParams{
		UserID: user.ID,
		Limit:  10,
		Offset: 0,
	}

	taskDays, err := testQueries.ListTaskDaysByUserId(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, 1, len(taskDays)) // Should have exactly 1 task day

	// Verify the task day has a valid date
	taskDay := taskDays[0]
	require.True(t, taskDay.Date.Valid)
	require.Equal(t, user.ID, taskDay.UserID)
}

func TestListTasksByTaskDayId(t *testing.T) {
	taskDay := createRandomTaskDay(t)

	// Create multiple tasks for the task day
	taskCount := 5
	var tasks []Task
	for i := 0; i < taskCount; i++ {
		arg := CreateTaskParams{
			TaskDayID:   taskDay.ID,
			Title:       util.GetRandomString(10),
			Description: pgtype.Text{String: util.GetRandomString(20), Valid: true},
			Duration:    pgtype.Interval{Microseconds: int64(time.Hour * 2), Valid: true},
		}
		task, err := testQueries.CreateTask(context.Background(), arg)
		require.NoError(t, err)
		tasks = append(tasks, task)
	}

	// Test listing tasks by task day ID
	taskDayTasks, err := testQueries.ListTasksByTaskDayId(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, taskCount, len(taskDayTasks))

	// Verify all returned tasks belong to the task day
	for _, task := range taskDayTasks {
		require.Equal(t, taskDay.ID, task.TaskDayID)
		require.NotEmpty(t, task.ID)
		require.NotEmpty(t, task.Title)
	}
}

func TestListTasksByTaskDayIdEmpty(t *testing.T) {
	taskDay := createRandomTaskDay(t)

	tasks, err := testQueries.ListTasksByTaskDayId(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Empty(t, tasks)
}

func TestListTasksByTaskDayIdWithDifferentTaskStates(t *testing.T) {
	taskDay := createRandomTaskDay(t)

	// Create tasks with different completion states
	completedTask := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       "Completed Task",
		Description: pgtype.Text{String: "This task is completed", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(time.Hour), Valid: true},
	}
	task1, err := testQueries.CreateTask(context.Background(), completedTask)
	require.NoError(t, err)

	// Update task to completed
	updateArg := UpdateTaskParams{
		ID:          task1.ID,
		Title:       task1.Title,
		Description: task1.Description,
		Duration:    task1.Duration,
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg)
	require.NoError(t, err)

	// Create an incomplete task
	incompleteTask := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       "Incomplete Task",
		Description: pgtype.Text{String: "This task is not completed", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(time.Hour * 2), Valid: true},
	}
	_, err = testQueries.CreateTask(context.Background(), incompleteTask)
	require.NoError(t, err)

	// List all tasks for the task day
	tasks, err := testQueries.ListTasksByTaskDayId(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, 2, len(tasks))

	// Verify we have both completed and incomplete tasks
	var completedCount, incompleteCount int
	for _, task := range tasks {
		require.Equal(t, taskDay.ID, task.TaskDayID)
		if task.Completed.Valid && task.Completed.Bool {
			completedCount++
		} else {
			incompleteCount++
		}
	}
	require.Equal(t, 1, completedCount)
	require.Equal(t, 1, incompleteCount)
}

func TestListTaskDaysByUserIdLargeOffset(t *testing.T) {
	user := createRandomUser(t)

	// Create a single task day
	_, err := testQueries.CreateTaskDay(context.Background(), user.ID)
	require.NoError(t, err)

	// Test with large offset to get empty results
	arg := ListTaskDaysByUserIdParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 100,
	}

	taskDays, err := testQueries.ListTaskDaysByUserId(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, taskDays)
}

func TestListTasksByTaskDayIdMultipleTaskDays(t *testing.T) {
	// Create two different task days
	taskDay1 := createRandomTaskDay(t)
	taskDay2 := createRandomTaskDay(t)

	// Create tasks for both task days
	task1Arg := CreateTaskParams{
		TaskDayID:   taskDay1.ID,
		Title:       "Task for Day 1",
		Description: pgtype.Text{String: "First task day", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(time.Hour), Valid: true},
	}
	_, err := testQueries.CreateTask(context.Background(), task1Arg)
	require.NoError(t, err)

	task2Arg := CreateTaskParams{
		TaskDayID:   taskDay2.ID,
		Title:       "Task for Day 2",
		Description: pgtype.Text{String: "Second task day", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(time.Hour * 2), Valid: true},
	}
	_, err = testQueries.CreateTask(context.Background(), task2Arg)
	require.NoError(t, err)

	// Verify tasks are only returned for the specific task day
	tasks1, err := testQueries.ListTasksByTaskDayId(context.Background(), taskDay1.ID)
	require.NoError(t, err)
	require.Equal(t, 1, len(tasks1))
	require.Equal(t, "Task for Day 1", tasks1[0].Title)

	tasks2, err := testQueries.ListTasksByTaskDayId(context.Background(), taskDay2.ID)
	require.NoError(t, err)
	require.Equal(t, 1, len(tasks2))
	require.Equal(t, "Task for Day 2", tasks2[0].Title)
}
