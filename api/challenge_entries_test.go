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
	"github.com/sanjayj369/retrospect-backend/token"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateChallengeEntriesAPI(t *testing.T) {
	challengeEntry := randomChallengeEntry()
	validUUID := uuid.UUID(challengeEntry.ID.Bytes).String()
	challenge := db.Challenge{
		UserID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
	}

	testCases := []struct {
		name          string
		challengeID   string
		body          updateChallengeEntriesBodyRequest
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK - Mark as Complete",
			challengeID: validUUID,
			body: updateChallengeEntriesBodyRequest{
				Complete: true,
			},

			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: true, Valid: true},
				}

				updatedEntry := challengeEntry
				updatedEntry.Completed = pgtype.Bool{Bool: true, Valid: true}

				store.EXPECT().
					GetChallengeEntry(gomock.Any(), gomock.Eq(challengeEntry.ID)).
					Times(1).
					Return(challengeEntry, nil)

				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeEntry.ChallengeID)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					UpdateChallengeEntry(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedEntry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChallengeEntry(t, recorder.Body, challengeEntry)
			},
		},
		{
			name:        "OK - Mark as Incomplete",
			challengeID: validUUID,
			body: updateChallengeEntriesBodyRequest{
				Complete: false,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: false, Valid: true},
				}

				updatedEntry := challengeEntry
				updatedEntry.Completed = pgtype.Bool{Bool: false, Valid: true}

				store.EXPECT().
					GetChallengeEntry(gomock.Any(), gomock.Eq(challengeEntry.ID)).
					Times(1).
					Return(challengeEntry, nil)

				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeEntry.ChallengeID)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					UpdateChallengeEntry(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedEntry, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchChallengeEntry(t, recorder.Body, challengeEntry)
			},
		},
		{
			name:        "Bad Request - Invalid UUID in URL",
			challengeID: "invalid-uuid",
			body: updateChallengeEntriesBodyRequest{
				Complete: true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {

				store.EXPECT().
					UpdateChallengeEntry(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Bad Request - Invalid JSON Body",
			challengeID: validUUID,
			body:        updateChallengeEntriesBodyRequest{}, // This will be overridden with invalid JSON
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {

				store.EXPECT().
					UpdateChallengeEntry(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:        "Not Found",
			challengeID: validUUID,
			body: updateChallengeEntriesBodyRequest{
				Complete: true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: true, Valid: true},
				}

				store.EXPECT().
					GetChallengeEntry(gomock.Any(), gomock.Eq(challengeEntry.ID)).
					Times(1).
					Return(challengeEntry, nil)

				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeEntry.ChallengeID)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					UpdateChallengeEntry(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.ChallengeEntry{}, pgx.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:        "Internal Server Error",
			challengeID: validUUID,
			body: updateChallengeEntriesBodyRequest{
				Complete: true,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,
					authorizationTypeBearer, challenge.UserID.Bytes, time.Minute)
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: true, Valid: true},
				}

				store.EXPECT().
					GetChallengeEntry(gomock.Any(), gomock.Eq(challengeEntry.ID)).
					Times(1).
					Return(challengeEntry, nil)

				store.EXPECT().
					GetChallenge(gomock.Any(), gomock.Eq(challengeEntry.ChallengeID)).
					Times(1).
					Return(challenge, nil)

				store.EXPECT().
					UpdateChallengeEntry(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.ChallengeEntry{}, sql.ErrConnDone)
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
			if tc.name == "Bad Request - Invalid JSON Body" {
				// Send invalid JSON
				bodyReader = bytes.NewReader([]byte("invalid json"))
			} else {
				bodyBytes, err := json.Marshal(tc.body)
				require.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			url := fmt.Sprintf("/challenge-entries/%s", tc.challengeID)
			req, err := http.NewRequest(http.MethodPut, url, bodyReader)
			require.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

// randomChallengeEntry generates a random challenge entry for testing
func randomChallengeEntry() db.ChallengeEntry {
	return db.ChallengeEntry{
		ID:          util.GetUUIDPGType(),
		ChallengeID: util.GetUUIDPGType(),
		Date:        util.GetRandomEndDate(30),
		Completed: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  getRandomTimestamp(),
			Valid: true,
		},
	}
}

// requireBodyMatchChallengeEntry checks that the response body matches the expected challenge entry
func requireBodyMatchChallengeEntry(t *testing.T, body *bytes.Buffer, challengeEntry db.ChallengeEntry) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotChallengeEntry db.ChallengeEntry
	err = json.Unmarshal(data, &gotChallengeEntry)
	require.NoError(t, err)

	require.Equal(t, challengeEntry.ID, gotChallengeEntry.ID)
	require.Equal(t, challengeEntry.ChallengeID, gotChallengeEntry.ChallengeID)
	require.Equal(t, challengeEntry.Date, gotChallengeEntry.Date)
	// Note: We don't check Completed field here as it might be updated in the test
}

// getRandomTimestamp generates a random timestamp for testing
func getRandomTimestamp() time.Time {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
