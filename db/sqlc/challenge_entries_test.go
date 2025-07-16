package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

// createRandomChallengeEntry creates a new challenge and then
// creates a challenge entry for that challenge
func createRandomChallengeEntry(t testing.TB) ChallengeEntry {
	t.Helper()
	challenge := createRandomChallenge(t)

	challengeEntry, err := testQueries.CreateChallengeEntry(context.Background(), challenge.ID)
	require.NoError(t, err)
	require.NotEmpty(t, challengeEntry)
	require.Equal(t, challenge.ID, challengeEntry.ChallengeID)
	require.NotEmpty(t, challengeEntry.ID)
	require.NotZero(t, challengeEntry.CreatedAt)

	return challengeEntry
}

func TestCreateChallengeEntry(t *testing.T) {
	createRandomChallengeEntry(t)
}

func TestGetChallengeEntry(t *testing.T) {
	challengeEntry1 := createRandomChallengeEntry(t)
	challengeEntry2, err := testQueries.GetChallengeEntry(context.Background(), challengeEntry1.ID)
	require.NoError(t, err)
	require.Equal(t, challengeEntry1, challengeEntry2)
}

func TestDeleteChallengeEntry(t *testing.T) {
	challengeEntry1 := createRandomChallengeEntry(t)
	challengeEntry2, err := testQueries.DeleteChallengeEntry(context.Background(), challengeEntry1.ID)
	require.NoError(t, err)
	require.Equal(t, challengeEntry1, challengeEntry2)

	// Verify the challenge entry was deleted
	challengeEntry3, err := testQueries.GetChallengeEntry(context.Background(), challengeEntry1.ID)
	require.Error(t, err)
	require.Empty(t, challengeEntry3)
	require.Equal(t, err, sql.ErrNoRows)
}

func TestListChallengeEntries(t *testing.T) {

	count := 5
	for i := 0; i < count; i++ {
		createRandomChallengeEntry(t)
	}

	arg := ListChallengeEntriesParams{
		Limit:  3,
		Offset: 1,
	}
	challengeEntries, err := testQueries.ListChallengeEntries(context.Background(), arg)
	require.NoError(t, err)
	require.LessOrEqual(t, len(challengeEntries), int(arg.Limit))

	for _, challengeEntry := range challengeEntries {
		require.NotEmpty(t, challengeEntry)
	}
}

func TestUpdateChallengeEntry(t *testing.T) {
	challengeEntry1 := createRandomChallengeEntry(t)

	// Test updating to completed
	arg := UpdateChallengeEntryParams{
		ID: challengeEntry1.ID,
		Completed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}

	err := testQueries.UpdateChallengeEntry(context.Background(), arg)
	require.NoError(t, err)

	challengeEntry2, err := testQueries.GetChallengeEntry(context.Background(), challengeEntry1.ID)
	require.NoError(t, err)
	require.True(t, challengeEntry2.Completed.Valid)
	require.True(t, challengeEntry2.Completed.Bool)
	require.Equal(t, challengeEntry1.ID, challengeEntry2.ID)
	require.Equal(t, challengeEntry1.ChallengeID, challengeEntry2.ChallengeID)
}

func TestUpdateChallengeEntryToNotCompleted(t *testing.T) {
	challengeEntry1 := createRandomChallengeEntry(t)

	arg := UpdateChallengeEntryParams{
		ID: challengeEntry1.ID,
		Completed: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
	}

	err := testQueries.UpdateChallengeEntry(context.Background(), arg)
	require.NoError(t, err)

	// Verify the update
	challengeEntry2, err := testQueries.GetChallengeEntry(context.Background(), challengeEntry1.ID)
	require.NoError(t, err)
	require.True(t, challengeEntry2.Completed.Valid)
	require.False(t, challengeEntry2.Completed.Bool)
}

func TestUpdateChallengeEntryToNull(t *testing.T) {
	challengeEntry1 := createRandomChallengeEntry(t)

	arg := UpdateChallengeEntryParams{
		ID: challengeEntry1.ID,
		Completed: sql.NullBool{
			Bool:  false,
			Valid: false,
		},
	}

	err := testQueries.UpdateChallengeEntry(context.Background(), arg)
	require.NoError(t, err)

	challengeEntry2, err := testQueries.GetChallengeEntry(context.Background(), challengeEntry1.ID)
	require.NoError(t, err)
	require.False(t, challengeEntry2.Completed.Valid) // Should be NULL
}
