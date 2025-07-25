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
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	mockDB "github.com/sanjayj369/retrospect-backend/db/mock"
	db "github.com/sanjayj369/retrospect-backend/db/sqlc"
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateChallengeAPI(t *testing.T) {
	challenge := randomChallenge()
	userID := uuid.UUID(challenge.UserID.Bytes).String()
	endDate := challenge.EndDate.Time.Format("2006-01-02")

	testCases := []struct {
		name          string
		body          createChallengeRequest
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: createChallengeRequest{
				Title:       challenge.Title,
				UserID:      userID,
				Description: challenge.Description.String,
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				// Parse the date from the request body and create the expected pgtype.Date
				expectedEndDate, err := time.Parse("2006-01-02", endDate)
				require.NoError(t, err)
				expectedEndDatePg := pgtype.Date{Time: expectedEndDate, Valid: true}

				arg := db.CreateChallengeParams{
					Title:       challenge.Title,
					UserID:      challenge.UserID,
					Description: challenge.Description,
					EndDate:     expectedEndDatePg,
				}

				store.EXPECT().
					CreateChallenge(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(challenge, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchChallenge(t, recorder.Body, challenge)
			},
		},
		{
			name: "Bad Request - Invalid JSON",
			body: createChallengeRequest{}, // This will be overridden with invalid JSON
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateChallenge(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request - Missing Title",
			body: createChallengeRequest{
				UserID:      userID,
				Description: challenge.Description.String,
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateChallenge(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Bad Request - Invalid User ID",
			body: createChallengeRequest{
				Title:       challenge.Title,
				UserID:      "invalid-uuid",
				Description: challenge.Description.String,
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateChallenge(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "OK - Without End Date",
			body: createChallengeRequest{
				Title:       challenge.Title,
				UserID:      userID,
				Description: challenge.Description.String,
				EndDate:     "invalid-date", // This will make EndDate invalid
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				expectedArg := db.CreateChallengeParams{
					Title:       challenge.Title,
					UserID:      challenge.UserID,
					Description: challenge.Description,
					EndDate:     pgtype.Date{Valid: false}, // Should be invalid
				}

				challengeWithoutEndDate := challenge
				challengeWithoutEndDate.EndDate = pgtype.Date{Valid: false}

				store.EXPECT().
					CreateChallenge(gomock.Any(), gomock.Eq(expectedArg)).
					Times(1).
					Return(challengeWithoutEndDate, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Internal Server Error",
			body: createChallengeRequest{
				Title:       challenge.Title,
				UserID:      userID,
				Description: challenge.Description.String,
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					CreateChallenge(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Challenge{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
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

			url := "/challenges"
			req, err := http.NewRequest(http.MethodPost, url, bodyReader)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetChallengeAPI(t *testing.T) {
	challenge := randomChallenge()
	validUUID := uuid.UUID(challenge.ID.Bytes).String()

	testCases := []struct {
		name          string
		challengeID   string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			challengeID: validUUID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challenge.ID)).
					Times(1).
					Return(challenge, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChallenge(t, recorder.Body, challenge)
			},
		},
		{
			name:        "Bad Request - Invalid UUID",
			challengeID: "invalid-uuid",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Not Found",
			challengeID: validUUID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challenge.ID)).
					Times(1).
					Return(db.Challenge{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Internal Server Error",
			challengeID: validUUID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challenge.ID)).
					Times(1).
					Return(db.Challenge{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/challenges/%s", tc.challengeID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateChallengeAPI(t *testing.T) {
	challenge := randomChallenge()
	validUUID := uuid.UUID(challenge.ID.Bytes).String()
	endDate := challenge.EndDate.Time.Format("2006-01-02")
	challengeIDPGType := pgtype.UUID{Bytes: challenge.ID.Bytes, Valid: true}
	testCases := []struct {
		name          string
		challengeID   string
		body          updateChallengeRequest
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			challengeID: validUUID,
			body: updateChallengeRequest{
				ID:          validUUID,
				Title:       "Updated Title",
				Description: "Updated Description",
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				parsedEndDate, err := time.Parse("2006-01-02", endDate)
				require.NoError(t, err)
				arg := db.UpdateChallengeDetailsParams{
					ID:          challenge.ID,
					Title:       "Updated Title",
					Description: pgtype.Text{String: "Updated Description", Valid: true},
					EndDate:     pgtype.Date{Time: parsedEndDate, Valid: true},
				}

				updatedChallenge := challenge
				updatedChallenge.Title = "Updated Title"
				updatedChallenge.Description = pgtype.Text{String: "Updated Description", Valid: true}
				updatedChallenge.EndDate = pgtype.Date{Time: parsedEndDate, Valid: true}

				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeIDPGType)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					UpdateChallengeDetails(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedChallenge, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:        "Bad Request - Invalid JSON",
			challengeID: validUUID,
			body:        updateChallengeRequest{}, // This will be overridden with invalid JSON
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateChallengeDetails(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Bad Request - Missing Title",
			challengeID: validUUID,
			body: updateChallengeRequest{
				ID:          validUUID,
				Description: "Updated Description",
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					UpdateChallengeDetails(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Internal Server Error",
			challengeID: validUUID,
			body: updateChallengeRequest{
				ID:          validUUID,
				Title:       "Updated Title",
				Description: "Updated Description",
				EndDate:     endDate,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeIDPGType)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					UpdateChallengeDetails(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Challenge{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
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

			url := fmt.Sprintf("/challenges/%s", tc.challengeID)
			req, err := http.NewRequest(http.MethodPatch, url, bodyReader)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestDeleteChallengeAPI(t *testing.T) {
	challenge := randomChallenge()
	validUUID := uuid.UUID(challenge.ID.Bytes).String()
	challengeIDPGType := pgtype.UUID{Bytes: challenge.ID.Bytes, Valid: true}

	testCases := []struct {
		name          string
		challengeID   string
		buildStub     func(store *mockDB.MockStore)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK",
			challengeID: validUUID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeIDPGType)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					DeleteChallenge(gomock.Any(), gomock.Eq(challenge.ID)).
					Times(1).
					Return(challenge, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChallenge(t, recorder.Body, challenge)
			},
		},
		{
			name:        "Bad Request - Invalid UUID",
			challengeID: "invalid-uuid",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					DeleteChallenge(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Internal Server Error",
			challengeID: validUUID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeIDPGType)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					DeleteChallenge(gomock.Any(), gomock.Eq(challenge.ID)).
					Times(1).
					Return(db.Challenge{}, sql.ErrConnDone)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/challenges/%s", tc.challengeID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestListChallengesAPI(t *testing.T) {
	user := randomUser()
	userID := uuid.UUID(user.ID.Bytes).String()
	n := 5
	challenges := make([]db.Challenge, n)
	for i := 0; i < n; i++ {
		challenges[i] = randomChallenge()
		challenges[i].UserID = user.ID // Set same user ID for all challenges
	}

	testCases := []struct {
		name          string
		userID        string
		pageSize      int32
		pageIdx       int32
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			userID:   userID,
			pageSize: 5,
			pageIdx:  0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.ListChallengesByUserParams{
					UserID: user.ID,
					Limit:  5,
					Offset: 0,
				}

				store.EXPECT().
					ListChallengesByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(challenges, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChallenges(t, recorder.Body, challenges)
			},
		},
		{
			name:     "OK - Default Page Size",
			userID:   userID,
			pageSize: 0, // Should default to 10
			pageIdx:  0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.ListChallengesByUserParams{
					UserID: user.ID,
					Limit:  10, // Default page size
					Offset: 0,
				}

				store.EXPECT().
					ListChallengesByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(challenges, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "OK - With Pagination",
			userID:   userID,
			pageSize: 2,
			pageIdx:  1,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.ListChallengesByUserParams{
					UserID: user.ID,
					Limit:  2,
					Offset: 2, // pageIdx * pageSize
				}

				store.EXPECT().
					ListChallengesByUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(challenges[:2], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:     "Bad Request - Invalid User ID",
			userID:   "invalid-uuid",
			pageSize: 5,
			pageIdx:  0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListChallengesByUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:     "Internal Server Error",
			userID:   userID,
			pageSize: 5,
			pageIdx:  0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, user.ID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				store.EXPECT().
					ListChallengesByUser(gomock.Any(), gomock.Any()).
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Use the corrected route /users/:id/challenges
			url := fmt.Sprintf("/users/%s/challenges?page_size=%d&page_idx=%d", tc.userID, tc.pageSize, tc.pageIdx)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// randomChallenge generates a random challenge for testing
func randomChallenge() db.Challenge {
	return db.Challenge{
		ID:     util.GetUUIDPGType(),
		Title:  util.GetRandomString(10),
		UserID: util.GetUUIDPGType(),
		Description: pgtype.Text{
			String: util.GetRandomString(50),
			Valid:  true,
		},
		StartDate: util.GetRandomEndDate(30),
		EndDate:   util.GetRandomEndDate(30),
		Active: pgtype.Bool{
			Bool:  true,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  getRandomTimestamp(),
			Valid: true,
		},
	}
}

// requireBodyMatchChallenge checks that the response body matches the expected challenge
func requireBodyMatchChallenge(t *testing.T, body *bytes.Buffer, challenge db.Challenge) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotChallenge db.Challenge
	err = json.Unmarshal(data, &gotChallenge)
	require.NoError(t, err)

	require.Equal(t, challenge.ID, gotChallenge.ID)
	require.Equal(t, challenge.Title, gotChallenge.Title)
	require.Equal(t, challenge.UserID, gotChallenge.UserID)
	require.Equal(t, challenge.Description, gotChallenge.Description)
	require.Equal(t, challenge.StartDate, gotChallenge.StartDate)
	require.Equal(t, challenge.EndDate, gotChallenge.EndDate)
	require.Equal(t, challenge.Active, gotChallenge.Active)
}

// requireBodyMatchChallenges checks that the response body matches the expected challenges list
func requireBodyMatchChallenges(t *testing.T, body *bytes.Buffer, challenges []db.Challenge) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotChallenges []db.Challenge
	err = json.Unmarshal(data, &gotChallenges)
	require.NoError(t, err)

	require.Equal(t, len(challenges), len(gotChallenges))
	for i := range challenges {
		require.Equal(t, challenges[i].ID, gotChallenges[i].ID)
		require.Equal(t, challenges[i].Title, gotChallenges[i].Title)
		require.Equal(t, challenges[i].UserID, gotChallenges[i].UserID)
	}
}
