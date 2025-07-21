package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	mockDB "github.com/sanjayj369/retrospect-backend/db/mock"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserAPI(t *testing.T) {
	user := randomUser()
	validUUID := uuid.UUID(user.ID.Bytes).String()
	testcases := []struct {
		name          string
		userId        string
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "Not Found",
			userId: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalServerError",
			userId: validUUID,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		}, {
			name:   "InvalidID",
			userId: "abcd",
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Any()).
					Times(0).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockDB.NewMockStore(ctrl)
			tc.buildStub(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.userId)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestCreateUserAPI(t *testing.T) {
	user := randomUser()
	createUserReq := db.CreateUserParams{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.HashedPassword,
	}

	marshalledReq, err := json.Marshal(createUserReq)
	require.NoError(t, err)
	validUserDetails := bytes.NewReader(marshalledReq)
	invalidUserDetails := bytes.NewReader([]byte("hello world :)"))

	testCases := []struct {
		name          string
		userDetails   *bytes.Reader
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			userDetails: validUserDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name:        "BadRequest",
			userDetails: invalidUserDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "InternalServerError",
			userDetails: validUserDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
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
			tc.userDetails.Seek(0, 0)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/users"

			req, err := http.NewRequest(http.MethodPost, url, tc.userDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// randomUser generates a random user for testing
func randomUser() db.User {
	return db.User{
		ID:             util.GetUUIDPGType(),
		Email:          util.GetRandomString(10) + "@example.com",
		Name:           util.GetRandomString(10),
		Timezone:       util.GetRandomTimezone(),
		HashedPassword: util.GetRandomString(32),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User

	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.Equal(t, user.ID, gotUser.ID)
	require.Equal(t, user.Name, gotUser.Name)
	require.Equal(t, user.Email, gotUser.Email)
}
