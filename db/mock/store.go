// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sanjayj369/retrospect-backend/db/sqlc (interfaces: Store)
//
// Generated by this command:
//
//	mockgen -package mockDB -destination ./db/mock/store.go github.com/sanjayj369/retrospect-backend/db/sqlc Store
//

// Package mockDB is a generated GoMock package.
package mockDB

import (
	context "context"
	reflect "reflect"

	pgtype "github.com/jackc/pgx/v5/pgtype"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	gomock "go.uber.org/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
	isgomock struct{}
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateChallenge mocks base method.
func (m *MockStore) CreateChallenge(ctx context.Context, arg db.CreateChallengeParams) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChallenge", ctx, arg)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChallenge indicates an expected call of CreateChallenge.
func (mr *MockStoreMockRecorder) CreateChallenge(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChallenge", reflect.TypeOf((*MockStore)(nil).CreateChallenge), ctx, arg)
}

// CreateChallengeEntry mocks base method.
func (m *MockStore) CreateChallengeEntry(ctx context.Context, challengeID pgtype.UUID) (db.ChallengeEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateChallengeEntry", ctx, challengeID)
	ret0, _ := ret[0].(db.ChallengeEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateChallengeEntry indicates an expected call of CreateChallengeEntry.
func (mr *MockStoreMockRecorder) CreateChallengeEntry(ctx, challengeID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateChallengeEntry", reflect.TypeOf((*MockStore)(nil).CreateChallengeEntry), ctx, challengeID)
}

// CreateSession mocks base method.
func (m *MockStore) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, arg)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockStoreMockRecorder) CreateSession(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockStore)(nil).CreateSession), ctx, arg)
}

// CreateTask mocks base method.
func (m *MockStore) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, arg)
	ret0, _ := ret[0].(db.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockStoreMockRecorder) CreateTask(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockStore)(nil).CreateTask), ctx, arg)
}

// CreateTaskDay mocks base method.
func (m *MockStore) CreateTaskDay(ctx context.Context, userID pgtype.UUID) (db.TaskDay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTaskDay", ctx, userID)
	ret0, _ := ret[0].(db.TaskDay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTaskDay indicates an expected call of CreateTaskDay.
func (mr *MockStoreMockRecorder) CreateTaskDay(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTaskDay", reflect.TypeOf((*MockStore)(nil).CreateTaskDay), ctx, userID)
}

// CreateTaskDaysForUsersInTimezone mocks base method.
func (m *MockStore) CreateTaskDaysForUsersInTimezone(ctx context.Context, timezone pgtype.Interval) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTaskDaysForUsersInTimezone", ctx, timezone)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTaskDaysForUsersInTimezone indicates an expected call of CreateTaskDaysForUsersInTimezone.
func (mr *MockStoreMockRecorder) CreateTaskDaysForUsersInTimezone(ctx, timezone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTaskDaysForUsersInTimezone", reflect.TypeOf((*MockStore)(nil).CreateTaskDaysForUsersInTimezone), ctx, timezone)
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, arg)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), ctx, arg)
}

// DeleteChallenge mocks base method.
func (m *MockStore) DeleteChallenge(ctx context.Context, id pgtype.UUID) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChallenge", ctx, id)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteChallenge indicates an expected call of DeleteChallenge.
func (mr *MockStoreMockRecorder) DeleteChallenge(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChallenge", reflect.TypeOf((*MockStore)(nil).DeleteChallenge), ctx, id)
}

// DeleteChallengeEntry mocks base method.
func (m *MockStore) DeleteChallengeEntry(ctx context.Context, id pgtype.UUID) (db.ChallengeEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteChallengeEntry", ctx, id)
	ret0, _ := ret[0].(db.ChallengeEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteChallengeEntry indicates an expected call of DeleteChallengeEntry.
func (mr *MockStoreMockRecorder) DeleteChallengeEntry(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteChallengeEntry", reflect.TypeOf((*MockStore)(nil).DeleteChallengeEntry), ctx, id)
}

// DeleteTask mocks base method.
func (m *MockStore) DeleteTask(ctx context.Context, id pgtype.UUID) (db.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", ctx, id)
	ret0, _ := ret[0].(db.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockStoreMockRecorder) DeleteTask(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockStore)(nil).DeleteTask), ctx, id)
}

// DeleteTaskDay mocks base method.
func (m *MockStore) DeleteTaskDay(ctx context.Context, id pgtype.UUID) (db.TaskDay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTaskDay", ctx, id)
	ret0, _ := ret[0].(db.TaskDay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTaskDay indicates an expected call of DeleteTaskDay.
func (mr *MockStoreMockRecorder) DeleteTaskDay(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTaskDay", reflect.TypeOf((*MockStore)(nil).DeleteTaskDay), ctx, id)
}

// DeleteUser mocks base method.
func (m *MockStore) DeleteUser(ctx context.Context, id pgtype.UUID) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", ctx, id)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStoreMockRecorder) DeleteUser(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStore)(nil).DeleteUser), ctx, id)
}

// GetChallenge mocks base method.
func (m *MockStore) GetChallenge(ctx context.Context, id pgtype.UUID) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChallenge", ctx, id)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChallenge indicates an expected call of GetChallenge.
func (mr *MockStoreMockRecorder) GetChallenge(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChallenge", reflect.TypeOf((*MockStore)(nil).GetChallenge), ctx, id)
}

// GetChallengeEntry mocks base method.
func (m *MockStore) GetChallengeEntry(ctx context.Context, id pgtype.UUID) (db.ChallengeEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChallengeEntry", ctx, id)
	ret0, _ := ret[0].(db.ChallengeEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChallengeEntry indicates an expected call of GetChallengeEntry.
func (mr *MockStoreMockRecorder) GetChallengeEntry(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChallengeEntry", reflect.TypeOf((*MockStore)(nil).GetChallengeEntry), ctx, id)
}

// GetSessions mocks base method.
func (m *MockStore) GetSessions(ctx context.Context, id pgtype.UUID) (db.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSessions", ctx, id)
	ret0, _ := ret[0].(db.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSessions indicates an expected call of GetSessions.
func (mr *MockStoreMockRecorder) GetSessions(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSessions", reflect.TypeOf((*MockStore)(nil).GetSessions), ctx, id)
}

// GetTask mocks base method.
func (m *MockStore) GetTask(ctx context.Context, id pgtype.UUID) (db.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", ctx, id)
	ret0, _ := ret[0].(db.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockStoreMockRecorder) GetTask(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockStore)(nil).GetTask), ctx, id)
}

// GetTaskDay mocks base method.
func (m *MockStore) GetTaskDay(ctx context.Context, id pgtype.UUID) (db.TaskDay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskDay", ctx, id)
	ret0, _ := ret[0].(db.TaskDay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskDay indicates an expected call of GetTaskDay.
func (mr *MockStoreMockRecorder) GetTaskDay(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskDay", reflect.TypeOf((*MockStore)(nil).GetTaskDay), ctx, id)
}

// GetTimezonesWhereDayIsStarting mocks base method.
func (m *MockStore) GetTimezonesWhereDayIsStarting(ctx context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimezonesWhereDayIsStarting", ctx)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTimezonesWhereDayIsStarting indicates an expected call of GetTimezonesWhereDayIsStarting.
func (mr *MockStoreMockRecorder) GetTimezonesWhereDayIsStarting(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimezonesWhereDayIsStarting", reflect.TypeOf((*MockStore)(nil).GetTimezonesWhereDayIsStarting), ctx)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(ctx context.Context, id pgtype.UUID) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, id)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), ctx, id)
}

// GetUserByName mocks base method.
func (m *MockStore) GetUserByName(ctx context.Context, name string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByName", ctx, name)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByName indicates an expected call of GetUserByName.
func (mr *MockStoreMockRecorder) GetUserByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByName", reflect.TypeOf((*MockStore)(nil).GetUserByName), ctx, name)
}

// ListChallengeEntries mocks base method.
func (m *MockStore) ListChallengeEntries(ctx context.Context, arg db.ListChallengeEntriesParams) ([]db.ChallengeEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChallengeEntries", ctx, arg)
	ret0, _ := ret[0].([]db.ChallengeEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChallengeEntries indicates an expected call of ListChallengeEntries.
func (mr *MockStoreMockRecorder) ListChallengeEntries(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChallengeEntries", reflect.TypeOf((*MockStore)(nil).ListChallengeEntries), ctx, arg)
}

// ListChallengeEntriesByChallengeId mocks base method.
func (m *MockStore) ListChallengeEntriesByChallengeId(ctx context.Context, arg db.ListChallengeEntriesByChallengeIdParams) ([]db.ChallengeEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChallengeEntriesByChallengeId", ctx, arg)
	ret0, _ := ret[0].([]db.ChallengeEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChallengeEntriesByChallengeId indicates an expected call of ListChallengeEntriesByChallengeId.
func (mr *MockStoreMockRecorder) ListChallengeEntriesByChallengeId(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChallengeEntriesByChallengeId", reflect.TypeOf((*MockStore)(nil).ListChallengeEntriesByChallengeId), ctx, arg)
}

// ListChallenges mocks base method.
func (m *MockStore) ListChallenges(ctx context.Context, arg db.ListChallengesParams) ([]db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChallenges", ctx, arg)
	ret0, _ := ret[0].([]db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChallenges indicates an expected call of ListChallenges.
func (mr *MockStoreMockRecorder) ListChallenges(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChallenges", reflect.TypeOf((*MockStore)(nil).ListChallenges), ctx, arg)
}

// ListChallengesByUser mocks base method.
func (m *MockStore) ListChallengesByUser(ctx context.Context, arg db.ListChallengesByUserParams) ([]db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChallengesByUser", ctx, arg)
	ret0, _ := ret[0].([]db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChallengesByUser indicates an expected call of ListChallengesByUser.
func (mr *MockStoreMockRecorder) ListChallengesByUser(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChallengesByUser", reflect.TypeOf((*MockStore)(nil).ListChallengesByUser), ctx, arg)
}

// ListTaskDays mocks base method.
func (m *MockStore) ListTaskDays(ctx context.Context, arg db.ListTaskDaysParams) ([]db.TaskDay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTaskDays", ctx, arg)
	ret0, _ := ret[0].([]db.TaskDay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTaskDays indicates an expected call of ListTaskDays.
func (mr *MockStoreMockRecorder) ListTaskDays(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaskDays", reflect.TypeOf((*MockStore)(nil).ListTaskDays), ctx, arg)
}

// ListTaskDaysByUserId mocks base method.
func (m *MockStore) ListTaskDaysByUserId(ctx context.Context, arg db.ListTaskDaysByUserIdParams) ([]db.TaskDay, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTaskDaysByUserId", ctx, arg)
	ret0, _ := ret[0].([]db.TaskDay)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTaskDaysByUserId indicates an expected call of ListTaskDaysByUserId.
func (mr *MockStoreMockRecorder) ListTaskDaysByUserId(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTaskDaysByUserId", reflect.TypeOf((*MockStore)(nil).ListTaskDaysByUserId), ctx, arg)
}

// ListTasks mocks base method.
func (m *MockStore) ListTasks(ctx context.Context, arg db.ListTasksParams) ([]db.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasks", ctx, arg)
	ret0, _ := ret[0].([]db.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks.
func (mr *MockStoreMockRecorder) ListTasks(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockStore)(nil).ListTasks), ctx, arg)
}

// ListTasksByTaskDayId mocks base method.
func (m *MockStore) ListTasksByTaskDayId(ctx context.Context, taskDayID pgtype.UUID) ([]db.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasksByTaskDayId", ctx, taskDayID)
	ret0, _ := ret[0].([]db.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasksByTaskDayId indicates an expected call of ListTasksByTaskDayId.
func (mr *MockStoreMockRecorder) ListTasksByTaskDayId(ctx, taskDayID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasksByTaskDayId", reflect.TypeOf((*MockStore)(nil).ListTasksByTaskDayId), ctx, taskDayID)
}

// ListUsers mocks base method.
func (m *MockStore) ListUsers(ctx context.Context, arg db.ListUsersParams) ([]db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", ctx, arg)
	ret0, _ := ret[0].([]db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockStoreMockRecorder) ListUsers(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockStore)(nil).ListUsers), ctx, arg)
}

// UpdateChallengeActiveStatus mocks base method.
func (m *MockStore) UpdateChallengeActiveStatus(ctx context.Context, arg db.UpdateChallengeActiveStatusParams) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChallengeActiveStatus", ctx, arg)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChallengeActiveStatus indicates an expected call of UpdateChallengeActiveStatus.
func (mr *MockStoreMockRecorder) UpdateChallengeActiveStatus(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChallengeActiveStatus", reflect.TypeOf((*MockStore)(nil).UpdateChallengeActiveStatus), ctx, arg)
}

// UpdateChallengeDescription mocks base method.
func (m *MockStore) UpdateChallengeDescription(ctx context.Context, arg db.UpdateChallengeDescriptionParams) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChallengeDescription", ctx, arg)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChallengeDescription indicates an expected call of UpdateChallengeDescription.
func (mr *MockStoreMockRecorder) UpdateChallengeDescription(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChallengeDescription", reflect.TypeOf((*MockStore)(nil).UpdateChallengeDescription), ctx, arg)
}

// UpdateChallengeDetails mocks base method.
func (m *MockStore) UpdateChallengeDetails(ctx context.Context, arg db.UpdateChallengeDetailsParams) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChallengeDetails", ctx, arg)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChallengeDetails indicates an expected call of UpdateChallengeDetails.
func (mr *MockStoreMockRecorder) UpdateChallengeDetails(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChallengeDetails", reflect.TypeOf((*MockStore)(nil).UpdateChallengeDetails), ctx, arg)
}

// UpdateChallengeEndDate mocks base method.
func (m *MockStore) UpdateChallengeEndDate(ctx context.Context, arg db.UpdateChallengeEndDateParams) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChallengeEndDate", ctx, arg)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChallengeEndDate indicates an expected call of UpdateChallengeEndDate.
func (mr *MockStoreMockRecorder) UpdateChallengeEndDate(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChallengeEndDate", reflect.TypeOf((*MockStore)(nil).UpdateChallengeEndDate), ctx, arg)
}

// UpdateChallengeEntry mocks base method.
func (m *MockStore) UpdateChallengeEntry(ctx context.Context, arg db.UpdateChallengeEntryParams) (db.ChallengeEntry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChallengeEntry", ctx, arg)
	ret0, _ := ret[0].(db.ChallengeEntry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChallengeEntry indicates an expected call of UpdateChallengeEntry.
func (mr *MockStoreMockRecorder) UpdateChallengeEntry(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChallengeEntry", reflect.TypeOf((*MockStore)(nil).UpdateChallengeEntry), ctx, arg)
}

// UpdateChallengeTitle mocks base method.
func (m *MockStore) UpdateChallengeTitle(ctx context.Context, arg db.UpdateChallengeTitleParams) (db.Challenge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChallengeTitle", ctx, arg)
	ret0, _ := ret[0].(db.Challenge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateChallengeTitle indicates an expected call of UpdateChallengeTitle.
func (mr *MockStoreMockRecorder) UpdateChallengeTitle(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChallengeTitle", reflect.TypeOf((*MockStore)(nil).UpdateChallengeTitle), ctx, arg)
}

// UpdateTask mocks base method.
func (m *MockStore) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", ctx, arg)
	ret0, _ := ret[0].(db.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockStoreMockRecorder) UpdateTask(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockStore)(nil).UpdateTask), ctx, arg)
}

// UpdateUserEmail mocks base method.
func (m *MockStore) UpdateUserEmail(ctx context.Context, arg db.UpdateUserEmailParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserEmail", ctx, arg)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserEmail indicates an expected call of UpdateUserEmail.
func (mr *MockStoreMockRecorder) UpdateUserEmail(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserEmail", reflect.TypeOf((*MockStore)(nil).UpdateUserEmail), ctx, arg)
}

// UpdateUserName mocks base method.
func (m *MockStore) UpdateUserName(ctx context.Context, arg db.UpdateUserNameParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserName", ctx, arg)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserName indicates an expected call of UpdateUserName.
func (mr *MockStoreMockRecorder) UpdateUserName(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserName", reflect.TypeOf((*MockStore)(nil).UpdateUserName), ctx, arg)
}

// UpdateUserTimezone mocks base method.
func (m *MockStore) UpdateUserTimezone(ctx context.Context, arg db.UpdateUserTimezoneParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserTimezone", ctx, arg)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserTimezone indicates an expected call of UpdateUserTimezone.
func (mr *MockStoreMockRecorder) UpdateUserTimezone(ctx, arg any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserTimezone", reflect.TypeOf((*MockStore)(nil).UpdateUserTimezone), ctx, arg)
}
