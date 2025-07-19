package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

func TestTaskDayAggregatesOnTaskInsertion(t *testing.T) {
	// Create a task day
	taskDay := createRandomTaskDay(t)

	// Verify initial state
	initialTaskDay, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), initialTaskDay.Count.Int32)
	require.Equal(t, int64(0), initialTaskDay.TotalDuration.Microseconds)
	require.Equal(t, int64(0), initialTaskDay.CompletedDuration.Microseconds)

	// Create first task (not completed)
	task1Duration := time.Hour * 2
	arg1 := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       util.GetRandomString(10),
		Description: pgtype.Text{String: "First task", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(task1Duration), Valid: true},
	}
	task1, err := testQueries.CreateTask(context.Background(), arg1)
	require.NoError(t, err)

	// Verify aggregates after first task insertion
	updatedTaskDay, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), updatedTaskDay.Count.Int32)
	require.Equal(t, int64(task1Duration), updatedTaskDay.TotalDuration.Microseconds)
	require.Equal(t, int64(0), updatedTaskDay.CompletedDuration.Microseconds) // Not completed yet

	// Create second task (completed)
	task2Duration := time.Hour * 3
	arg2 := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       util.GetRandomString(10),
		Description: pgtype.Text{String: "Second task", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(task2Duration), Valid: true},
	}
	task2, err := testQueries.CreateTask(context.Background(), arg2)
	require.NoError(t, err)

	// Mark second task as completed
	updateArg := UpdateTaskParams{
		ID:          task2.ID,
		Title:       task2.Title,
		Description: task2.Description,
		Duration:    task2.Duration,
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg)
	require.NoError(t, err)

	// Verify aggregates after second task insertion and completion
	finalTaskDay, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(2), finalTaskDay.Count.Int32)
	require.Equal(t, int64(task1Duration+task2Duration), finalTaskDay.TotalDuration.Microseconds)
	require.Equal(t, int64(task2Duration), finalTaskDay.CompletedDuration.Microseconds)

	// Clean up
	_, err = testQueries.DeleteTask(context.Background(), task1.ID)
	require.NoError(t, err)
	_, err = testQueries.DeleteTask(context.Background(), task2.ID)
	require.NoError(t, err)
}

func TestTaskDayAggregatesOnTaskUpdate(t *testing.T) {
	// Create a task day and task
	taskDay := createRandomTaskDay(t)

	originalDuration := time.Hour * 2
	arg := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       util.GetRandomString(10),
		Description: pgtype.Text{String: "Test task", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(originalDuration), Valid: true},
	}
	task, err := testQueries.CreateTask(context.Background(), arg)
	require.NoError(t, err)

	// Verify initial state after task creation
	taskDay1, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), taskDay1.Count.Int32)
	require.Equal(t, int64(originalDuration), taskDay1.TotalDuration.Microseconds)
	require.Equal(t, int64(0), taskDay1.CompletedDuration.Microseconds)

	// Update task to completed (no duration change)
	updateArg1 := UpdateTaskParams{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Duration:    task.Duration,
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg1)
	require.NoError(t, err)

	// Verify completed_duration is updated
	taskDay2, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), taskDay2.Count.Int32)
	require.Equal(t, int64(originalDuration), taskDay2.TotalDuration.Microseconds)
	require.Equal(t, int64(originalDuration), taskDay2.CompletedDuration.Microseconds)

	// Update task duration while keeping it completed
	newDuration := time.Hour * 4
	updateArg2 := UpdateTaskParams{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Duration:    pgtype.Interval{Microseconds: int64(newDuration), Valid: true},
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg2)
	require.NoError(t, err)

	// Verify both total and completed durations are updated
	taskDay3, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), taskDay3.Count.Int32)
	require.Equal(t, int64(newDuration), taskDay3.TotalDuration.Microseconds)
	require.Equal(t, int64(newDuration), taskDay3.CompletedDuration.Microseconds)

	// Update task to not completed
	updateArg3 := UpdateTaskParams{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Duration:    pgtype.Interval{Microseconds: int64(newDuration), Valid: true},
		Completed:   pgtype.Bool{Bool: false, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg3)
	require.NoError(t, err)

	// Verify completed_duration is reset to 0
	taskDay4, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), taskDay4.Count.Int32)
	require.Equal(t, int64(newDuration), taskDay4.TotalDuration.Microseconds)
	require.Equal(t, int64(0), taskDay4.CompletedDuration.Microseconds)

	// Clean up
	_, err = testQueries.DeleteTask(context.Background(), task.ID)
	require.NoError(t, err)
}

func TestTaskDayAggregatesOnTaskDeletion(t *testing.T) {
	// Create a task day
	taskDay := createRandomTaskDay(t)

	// Create multiple tasks with different states
	task1Duration := time.Hour * 1
	task2Duration := time.Hour * 2
	task3Duration := time.Hour * 3

	// Task 1: Not completed
	arg1 := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       "Task 1",
		Description: pgtype.Text{String: "Not completed", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(task1Duration), Valid: true},
	}
	task1, err := testQueries.CreateTask(context.Background(), arg1)
	require.NoError(t, err)

	// Task 2: Completed
	arg2 := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       "Task 2",
		Description: pgtype.Text{String: "Completed", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(task2Duration), Valid: true},
	}
	task2, err := testQueries.CreateTask(context.Background(), arg2)
	require.NoError(t, err)

	// Mark task2 as completed
	updateArg := UpdateTaskParams{
		ID:          task2.ID,
		Title:       task2.Title,
		Description: task2.Description,
		Duration:    task2.Duration,
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg)
	require.NoError(t, err)

	// Task 3: Completed
	arg3 := CreateTaskParams{
		TaskDayID:   taskDay.ID,
		Title:       "Task 3",
		Description: pgtype.Text{String: "Also completed", Valid: true},
		Duration:    pgtype.Interval{Microseconds: int64(task3Duration), Valid: true},
	}
	task3, err := testQueries.CreateTask(context.Background(), arg3)
	require.NoError(t, err)

	// Mark task3 as completed
	updateArg3 := UpdateTaskParams{
		ID:          task3.ID,
		Title:       task3.Title,
		Description: task3.Description,
		Duration:    task3.Duration,
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateArg3)
	require.NoError(t, err)

	// Verify initial aggregates
	taskDayBefore, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(3), taskDayBefore.Count.Int32)
	require.Equal(t, int64(task1Duration+task2Duration+task3Duration), taskDayBefore.TotalDuration.Microseconds)
	require.Equal(t, int64(task2Duration+task3Duration), taskDayBefore.CompletedDuration.Microseconds)

	// Delete the incomplete task (task1)
	_, err = testQueries.DeleteTask(context.Background(), task1.ID)
	require.NoError(t, err)

	// Verify aggregates after deleting incomplete task
	taskDayAfter1, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(2), taskDayAfter1.Count.Int32)
	require.Equal(t, int64(task2Duration+task3Duration), taskDayAfter1.TotalDuration.Microseconds)
	require.Equal(t, int64(task2Duration+task3Duration), taskDayAfter1.CompletedDuration.Microseconds)

	// Delete a completed task (task2)
	_, err = testQueries.DeleteTask(context.Background(), task2.ID)
	require.NoError(t, err)

	// Verify aggregates after deleting completed task
	taskDayAfter2, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(1), taskDayAfter2.Count.Int32)
	require.Equal(t, int64(task3Duration), taskDayAfter2.TotalDuration.Microseconds)
	require.Equal(t, int64(task3Duration), taskDayAfter2.CompletedDuration.Microseconds)

	// Delete the last task
	_, err = testQueries.DeleteTask(context.Background(), task3.ID)
	require.NoError(t, err)

	// Verify aggregates are reset to zero
	taskDayFinal, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), taskDayFinal.Count.Int32)
	require.Equal(t, int64(0), taskDayFinal.TotalDuration.Microseconds)
	require.Equal(t, int64(0), taskDayFinal.CompletedDuration.Microseconds)
}

func TestTaskDayAggregatesComplexScenario(t *testing.T) {
	// Create a task day
	taskDay := createRandomTaskDay(t)

	// Test a complex scenario with multiple operations
	durations := []time.Duration{
		time.Minute * 30,
		time.Hour * 1,
		time.Hour * 2,
		time.Minute * 45,
	}

	var tasks []Task
	var totalDuration time.Duration

	// Create multiple tasks
	for _, duration := range durations {
		arg := CreateTaskParams{
			TaskDayID:   taskDay.ID,
			Title:       util.GetRandomString(10),
			Description: pgtype.Text{String: "Task description", Valid: true},
			Duration:    pgtype.Interval{Microseconds: int64(duration), Valid: true},
		}
		task, err := testQueries.CreateTask(context.Background(), arg)
		require.NoError(t, err)
		tasks = append(tasks, task)
		totalDuration += duration
	}

	// Verify total_duration after all insertions
	taskDayAfterInsert, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(len(tasks)), taskDayAfterInsert.Count.Int32)
	require.Equal(t, int64(totalDuration), taskDayAfterInsert.TotalDuration.Microseconds)
	require.Equal(t, int64(0), taskDayAfterInsert.CompletedDuration.Microseconds)

	// Complete some tasks (first and third)
	var completedDuration time.Duration
	for _, i := range []int{0, 2} {
		updateArg := UpdateTaskParams{
			ID:          tasks[i].ID,
			Title:       tasks[i].Title,
			Description: tasks[i].Description,
			Duration:    tasks[i].Duration,
			Completed:   pgtype.Bool{Bool: true, Valid: true},
		}
		_, err = testQueries.UpdateTask(context.Background(), updateArg)
		require.NoError(t, err)
		completedDuration += durations[i]
	}

	// Verify aggregates after completing some tasks
	taskDayAfterCompletion, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(len(tasks)), taskDayAfterCompletion.Count.Int32)
	require.Equal(t, int64(totalDuration), taskDayAfterCompletion.TotalDuration.Microseconds)
	require.Equal(t, int64(completedDuration), taskDayAfterCompletion.CompletedDuration.Microseconds)

	// Change duration of a completed task
	newDuration := time.Hour * 5
	originalDuration := durations[0]
	updateDurationArg := UpdateTaskParams{
		ID:          tasks[0].ID,
		Title:       tasks[0].Title,
		Description: tasks[0].Description,
		Duration:    pgtype.Interval{Microseconds: int64(newDuration), Valid: true},
		Completed:   pgtype.Bool{Bool: true, Valid: true},
	}
	_, err = testQueries.UpdateTask(context.Background(), updateDurationArg)
	require.NoError(t, err)

	// Verify aggregates after duration change
	expectedTotalDuration := totalDuration - originalDuration + newDuration
	expectedCompletedDuration := completedDuration - originalDuration + newDuration

	taskDayAfterDurationChange, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(len(tasks)), taskDayAfterDurationChange.Count.Int32)
	require.Equal(t, int64(expectedTotalDuration), taskDayAfterDurationChange.TotalDuration.Microseconds)
	require.Equal(t, int64(expectedCompletedDuration), taskDayAfterDurationChange.CompletedDuration.Microseconds)

	// Clean up all tasks
	for _, task := range tasks {
		_, err = testQueries.DeleteTask(context.Background(), task.ID)
		require.NoError(t, err)
	}

	// Verify final state
	taskDayFinal, err := testQueries.GetTaskDay(context.Background(), taskDay.ID)
	require.NoError(t, err)
	require.Equal(t, int32(0), taskDayFinal.Count.Int32)
	require.Equal(t, int64(0), taskDayFinal.TotalDuration.Microseconds)
	require.Equal(t, int64(0), taskDayFinal.CompletedDuration.Microseconds)
}
