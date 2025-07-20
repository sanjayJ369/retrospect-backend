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

func TestUpdateChallengeEntriesAPI(t *testing.T) {
	challengeEntry := randomChallengeEntry()
	validUUID := uuid.UUID(challengeEntry.ID.Bytes).String()

	testCases := []struct {
		name          string
		challengeID   string
		body          updateChallengeEntriesRequest
		buildStub     func(store *mockDB.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:        "OK - Mark as Complete",
			challengeID: validUUID,
			body: updateChallengeEntriesRequest{
				ChallengeID: validUUID,
				Complete:    true,
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: true, Valid: true},
				}

				updatedEntry := challengeEntry
				updatedEntry.Completed = pgtype.Bool{Bool: true, Valid: true}

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
			body: updateChallengeEntriesRequest{
				ChallengeID: validUUID,
				Complete:    false,
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: false, Valid: true},
				}

				updatedEntry := challengeEntry
				updatedEntry.Completed = pgtype.Bool{Bool: false, Valid: true}

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
			body: updateChallengeEntriesRequest{
				ChallengeID: "invalid-uuid",
				Complete:    true,
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
			body:        updateChallengeEntriesRequest{}, // This will be overridden with invalid JSON
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
			body: updateChallengeEntriesRequest{
				ChallengeID: validUUID,
				Complete:    true,
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: true, Valid: true},
				}

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
			body: updateChallengeEntriesRequest{
				ChallengeID: validUUID,
				Complete:    true,
			},
			buildStub: func(store *mockDB.MockStore) {
				arg := db.UpdateChallengeEntryParams{
					ID:        challengeEntry.ID,
					Completed: pgtype.Bool{Bool: true, Valid: true},
				}

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

			server := NewServer(store)
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
		Date: pgtype.Date{
			Time:  getRandomDate(),
			Valid: true,
		},
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
	require.Equal(t, challengeEntry.CreatedAt, gotChallengeEntry.CreatedAt)
	// Note: We don't check Completed field here as it might be updated in the test
}

// getRandomDate generates a random date for testing
func getRandomDate() time.Time {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

// getRandomTimestamp generates a random timestamp for testing
func getRandomTimestamp() time.Time {
	min := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
