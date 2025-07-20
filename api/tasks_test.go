package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	mockDB "github.com/sanjayj369/retrospect-backend/db/mock"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateTaskAPI(t *testing.T) {
	task := randomTask()
	taskDayID := uuid.UUID(task.TaskDayID.Bytes).String()

	testCases := []struct {
		name          string
		body          createTaskRequest
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createTaskRequest{
				TaskDayID:   taskDayID,
				Title:       task.Title,
				Description: task.Description.String,
				Duration:    60, // 60 minutes
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.CreateTaskParams{
					TaskDayID:   task.TaskDayID,
					Title:       task.Title,
					Description: task.Description,
					Duration:    util.MinutesToPGInterval(60),
				}

				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name: "Bad Request - Invalid JSON",
			body: createTaskRequest{}, // This will be overridden with invalid JSON
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request - Missing Title",
			body: createTaskRequest{
				TaskDayID:   taskDayID,
				Description: task.Description.String,
				Duration:    60,
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request - Invalid Task Day ID",
			body: createTaskRequest{
				TaskDayID:   "invalid-uuid",
				Title:       task.Title,
				Description: task.Description.String,
				Duration:    60,
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request - Missing Duration",
			body: createTaskRequest{
				TaskDayID:   taskDayID,
				Title:       task.Title,
				Description: task.Description.String,
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "OK - Without Description",
			body: createTaskRequest{
				TaskDayID: taskDayID,
				Title:     task.Title,
				Duration:  60,
			},
			buildStub: func(store *mockDB.MockStore) {
				expectedArg := db.CreateTaskParams{
					TaskDayID:   task.TaskDayID,
					Title:       task.Title,
					Description: pgtype.Text{String: "", Valid: false}, // Empty description
					Duration:    util.MinutesToPGInterval(60),
				}

				taskWithoutDescription := task
				taskWithoutDescription.Description = pgtype.Text{String: "", Valid: false}

				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Eq(expectedArg)).
					Times(1).
					Return(taskWithoutDescription, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: createTaskRequest{
				TaskDayID:   taskDayID,
				Title:       task.Title,
				Description: task.Description.String,
				Duration:    60,
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Prepare request body
			var bodyReader *bytes.Reader
			if tc.name == "Bad Request - Invalid JSON" {
				bodyReader = bytes.NewReader([]byte("invalid json"))
			} else {
				bodyBytes, err := json.Marshal(tc.body)
				require.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			url := "/tasks"
			req, err := http.NewRequest(http.MethodPost, url, bodyReader)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetTaskAPI(t *testing.T) {
	task := randomTask()
	validUUID := uuid.UUID(task.ID.Bytes).String()

	testCases := []struct {
		name          string
		taskID        string
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			taskID: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTask(t, recorder.Body, task)
			},
		},
		{
			name:   "Bad Request - Invalid UUID",
			taskID: "invalid-uuid",
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Not Found",
			taskID: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(db.Task{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			taskID: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tasks/%s", tc.taskID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateTaskAPI(t *testing.T) {
	task := randomTask()
	validUUID := uuid.UUID(task.ID.Bytes).String()

	testCases := []struct {
		name          string
		taskID        string
		body          updateTaskBodyRequest
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK - Complete Task",
			taskID: validUUID,
			body: updateTaskBodyRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
				Completed:   true,
				Duration:    90,
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateTaskParams{
					ID:          task.ID,
					Title:       "Updated Title",
					Description: pgtype.Text{String: "Updated Description", Valid: true},
					Completed:   pgtype.Bool{Bool: true, Valid: true},
					Duration:    util.MinutesToPGInterval(90),
				}

				updatedTask := task
				updatedTask.Title = "Updated Title"
				updatedTask.Description = pgtype.Text{String: "Updated Description", Valid: true}
				updatedTask.Completed = pgtype.Bool{Bool: true, Valid: true}

				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedTask, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "OK - Mark as Incomplete",
			taskID: validUUID,
			body: updateTaskBodyRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
				Completed:   false,
				Duration:    90,
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateTaskParams{
					ID:          task.ID,
					Title:       "Updated Title",
					Description: pgtype.Text{String: "Updated Description", Valid: true},
					Completed:   pgtype.Bool{Bool: false, Valid: true},
					Duration:    util.MinutesToPGInterval(90),
				}

				updatedTask := task
				updatedTask.Title = "Updated Title"
				updatedTask.Description = pgtype.Text{String: "Updated Description", Valid: true}
				updatedTask.Completed = pgtype.Bool{Bool: false, Valid: true}

				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedTask, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "Bad Request - Invalid JSON",
			taskID: validUUID,
			body:   updateTaskBodyRequest{}, // This will be overridden with invalid JSON
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Bad Request - Missing Title",
			taskID: validUUID,
			body: updateTaskBodyRequest{
				Description: "Updated Description",
				Completed:   true,
				Duration:    90,
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			taskID: validUUID,
			body: updateTaskBodyRequest{
				Title:       "Updated Title",
				Description: "Updated Description",
				Completed:   true,
				Duration:    90,
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateTask(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// Prepare request body
			var bodyReader *bytes.Reader
			if tc.name == "Bad Request - Invalid JSON" {
				bodyReader = bytes.NewReader([]byte("invalid json"))
			} else {
				bodyBytes, err := json.Marshal(tc.body)
				require.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			url := fmt.Sprintf("/tasks/%s", tc.taskID)
			req, err := http.NewRequest(http.MethodPatch, url, bodyReader)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteTaskAPI(t *testing.T) {
	task := randomTask()
	validUUID := uuid.UUID(task.ID.Bytes).String()

	testCases := []struct {
		name          string
		taskID        string
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			taskID: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(task, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNoContent, recorder.Code)
			},
		},
		{
			name:   "Bad Request - Invalid UUID",
			taskID: "invalid-uuid",
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "Internal Server Error",
			taskID: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					DeleteTask(gomock.Any(), gomock.Eq(task.ID)).
					Times(1).
					Return(db.Task{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tasks/%s", tc.taskID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListTasksAPI(t *testing.T) {
	taskDay := randomTaskDay()
	taskDayID := uuid.UUID(taskDay.ID.Bytes).String()
	n := 3
	tasks := make([]db.Task, n)
	for i := 0; i < n; i++ {
		tasks[i] = randomTask()
		tasks[i].TaskDayID = taskDay.ID // Set same task day ID for all tasks
	}

	testCases := []struct {
		name          string
		taskDayID     string
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			taskDayID: taskDayID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListTasksByTaskDayId(gomock.Any(), gomock.Eq(taskDay.ID)).
					Times(1).
					Return(tasks, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTasks(t, recorder.Body, tasks)
			},
		},
		{
			name:      "Bad Request - Invalid Task Day ID",
			taskDayID: "invalid-uuid",
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListTasksByTaskDayId(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "Internal Server Error",
			taskDayID: taskDayID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListTasksByTaskDayId(gomock.Any(), gomock.Eq(taskDay.ID)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/tasks?task_day_id=%s", tc.taskDayID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// randomTask generates a random task for testing
func randomTask() db.Task {
	return db.Task{
		ID:        util.GetUUIDPGType(),
		TaskDayID: util.GetUUIDPGType(),
		Title:     util.GetRandomString(15),
		Description: pgtype.Text{
			String: util.GetRandomString(50),
			Valid:  true,
		},
		Duration: util.MinutesToPGInterval(int(getRandomDuration())),
		Completed: pgtype.Bool{
			Bool:  rand.Intn(2) == 1, // Random true/false
			Valid: true,
		},
	}
}

// randomTaskDay generates a random task day for testing
func randomTaskDay() db.TaskDay {
	return db.TaskDay{
		ID:     util.GetUUIDPGType(),
		UserID: util.GetUUIDPGType(),
		Date: pgtype.Date{
			Time:  getRandomDateUTC(),
			Valid: true,
		},
		Count: pgtype.Int4{
			Int32: rand.Int31n(10) + 1, // 1-10 tasks
			Valid: true,
		},
		TotalDuration:     util.MinutesToPGInterval(int(getRandomDuration())),
		CompletedDuration: util.MinutesToPGInterval(int(getRandomDuration())),
	}
}

// requireBodyMatchTask checks that the response body matches the expected task
func requireBodyMatchTask(t *testing.T, body *bytes.Buffer, task db.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTask db.Task
	err = json.Unmarshal(data, &gotTask)
	require.NoError(t, err)

	require.Equal(t, task.ID, gotTask.ID)
	require.Equal(t, task.TaskDayID, gotTask.TaskDayID)
	require.Equal(t, task.Title, gotTask.Title)
	require.Equal(t, task.Description, gotTask.Description)
	// Note: We don't check Duration and Completed fields in detail as they might vary
}

// requireBodyMatchTasks checks that the response body matches the expected tasks list
func requireBodyMatchTasks(t *testing.T, body *bytes.Buffer, tasks []db.Task) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTasks []db.Task
	err = json.Unmarshal(data, &gotTasks)
	require.NoError(t, err)

	require.Equal(t, len(tasks), len(gotTasks))
	for i := range tasks {
		require.Equal(t, tasks[i].ID, gotTasks[i].ID)
		require.Equal(t, tasks[i].TaskDayID, gotTasks[i].TaskDayID)
		require.Equal(t, tasks[i].Title, gotTasks[i].Title)
	}
}

// getRandomDuration generates a random duration in minutes for testing
func getRandomDuration() int64 {
	return rand.Int63n(480) + 15 // 15 to 495 minutes (8 hours max)
}

// getRandomDateUTC generates a random date in UTC for testing (reusing from challenge_test.go)
func getRandomDateUTC() time.Time {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)

	// Generate random day between min and max
	days := int(max.Sub(min).Hours() / 24)
	randomDays := rand.Intn(days)

	return min.AddDate(0, 0, randomDays)
}
