package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	mockDB "github.com/sanjayj369/retrospect-backend/db/mock"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUserAPI(t *testing.T) {
	createUserReq := createUserRequest{
		Name:     util.GetRandomString(6),
		Email:    util.GetRandomString(6) + "@example.com",
		Password: util.GetRandomString(6),
	}

	createUserParams := db.CreateUserParams{
		Name:           createUserReq.Name,
		Email:          createUserReq.Email,
		HashedPassword: "",
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
					CreateUser(gomock.Any(), EqCreateUserParams(createUserParams, createUserReq.Password)).
					Times(1).
					Return(randomUser(), nil)
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
					CreateUser(gomock.Any(), EqCreateUserParams(createUserParams, createUserReq.Password)).
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users"

			req, err := http.NewRequest(http.MethodPost, url, tc.userDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetUserAPI(t *testing.T) {
	user := randomUser()
	validUUID := uuid.UUID(user.ID.Bytes).String()
	testcases := []struct {
		name          string
		userId        string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userId: validUUID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
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
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.userId)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
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

type eqUserMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqUserMatcher) Matches(x any) bool {
	arg, err := x.(db.CreateUserParams)
	if !err {
		return false
	}

	if err := util.CheckPassword(e.password, arg.HashedPassword); err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUserMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.arg)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqUserMatcher{
		arg:      arg,
		password: password,
	}
}

func TestLoginUserAPI(t *testing.T) {
	user := randomUser()
	password := util.GetRandomString(8)
	hashedPassword, err := util.HashedPassword(password)
	require.NoError(t, err)

	user.HashedPassword = hashedPassword

	loginUserReq := UserLoginRequest{
		Name:     user.Name,
		Password: password,
	}

	marshalledReq, err := json.Marshal(loginUserReq)
	require.NoError(t, err)
	validLoginDetails := bytes.NewReader(marshalledReq)

	invalidLoginReq := UserLoginRequest{
		Name:     "",    // Invalid: required field empty
		Password: "123", // Invalid: too short
	}
	marshalledInvalidReq, err := json.Marshal(invalidLoginReq)
	require.NoError(t, err)
	invalidLoginDetails := bytes.NewReader(marshalledInvalidReq)

	wrongPasswordReq := UserLoginRequest{
		Name:     user.Name,
		Password: "wrongpassword123",
	}
	marshalledWrongPasswordReq, err := json.Marshal(wrongPasswordReq)
	require.NoError(t, err)
	wrongPasswordDetails := bytes.NewReader(marshalledWrongPasswordReq)

	testCases := []struct {
		name          string
		loginDetails  *bytes.Reader
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:         "OK",
			loginDetails: validLoginDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUserByName(gomock.Any(), gomock.Eq(user.Name)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				// Check that response contains access token and user data
				var response UserLoginResponse
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.AccessToken)
				require.Equal(t, user.Name, response.User.Name)
				require.Equal(t, user.Email, response.User.Email)
			},
		},
		{
			name:         "BadRequest",
			loginDetails: invalidLoginDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUserByName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:         "UserNotFound",
			loginDetails: validLoginDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUserByName(gomock.Any(), gomock.Eq(user.Name)).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:         "WrongPassword",
			loginDetails: wrongPasswordDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUserByName(gomock.Any(), gomock.Eq(user.Name)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:         "InternalServerError",
			loginDetails: validLoginDetails,
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetUserByName(gomock.Any(), gomock.Eq(user.Name)).
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
			tc.loginDetails.Seek(0, 0)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/users/login"

			req, err := http.NewRequest(http.MethodPost, url, tc.loginDetails)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
