package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

func TestListChallengesByUser(t *testing.T) {
	user := createRandomUser(t)

	// Create multiple challenges for the user
	challengeCount := 5
	var challenges []Challenge
	for i := 0; i < challengeCount; i++ {
		arg := CreateChallengeParams{
			Title:       "Test Challenge " + string(rune('A'+i)),
			UserID:      user.ID,
			Description: pgtype.Text{String: "Test description", Valid: true},
			EndDate:     util.GetRandomEndDate(30),
		}
		challenge, err := testQueries.CreateChallenge(context.Background(), arg)
		require.NoError(t, err)
		challenges = append(challenges, challenge)
	}

	// Test listing challenges
	arg := ListChallengesByUserParams{
		UserID: user.ID,
		Limit:  3,
		Offset: 0,
	}

	userChallenges, err := testQueries.ListChallengesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.LessOrEqual(t, len(userChallenges), int(arg.Limit))

	// Verify all returned challenges belong to the user
	for _, challenge := range userChallenges {
		require.Equal(t, user.ID, challenge.UserID)
		require.NotEmpty(t, challenge.ID)
		require.NotEmpty(t, challenge.Title)
	}
}

func TestListChallengesByUserEmpty(t *testing.T) {
	user := createRandomUser(t)

	arg := ListChallengesByUserParams{
		UserID: user.ID,
		Limit:  5,
		Offset: 0,
	}

	challenges, err := testQueries.ListChallengesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, challenges)
}

func TestListChallengesByUserWithPagination(t *testing.T) {
	user := createRandomUser(t)

	// Create multiple challenges for pagination testing
	challengeCount := 10
	for i := 0; i < challengeCount; i++ {
		arg := CreateChallengeParams{
			Title:       "Challenge " + string(rune('A'+i)),
			UserID:      user.ID,
			Description: pgtype.Text{String: "Description", Valid: true},
			EndDate:     util.GetRandomEndDate(30),
		}
		_, err := testQueries.CreateChallenge(context.Background(), arg)
		require.NoError(t, err)
	}

	// Test first page
	arg1 := ListChallengesByUserParams{
		UserID: user.ID,
		Limit:  3,
		Offset: 0,
	}
	firstPage, err := testQueries.ListChallengesByUser(context.Background(), arg1)
	require.NoError(t, err)
	require.LessOrEqual(t, len(firstPage), 3)

	// Test second page
	arg2 := ListChallengesByUserParams{
		UserID: user.ID,
		Limit:  3,
		Offset: 3,
	}
	secondPage, err := testQueries.ListChallengesByUser(context.Background(), arg2)
	require.NoError(t, err)
	require.LessOrEqual(t, len(secondPage), 3)

	// Verify pages don't contain the same challenges
	if len(firstPage) > 0 && len(secondPage) > 0 {
		require.NotEqual(t, firstPage[0].ID, secondPage[0].ID)
	}
}

func TestListChallengeEntriesByChallengeId(t *testing.T) {
	challenge := createRandomChallenge(t)

	// Create multiple entries for the challenge
	entryCount := 5
	var entries []ChallengeEntry
	for i := 0; i < entryCount; i++ {
		entry, err := testQueries.CreateChallengeEntry(context.Background(), challenge.ID)
		require.NoError(t, err)
		entries = append(entries, entry)
	}

	// Test listing entries
	arg := ListChallengeEntriesByChallengeIdParams{
		ChallengeID: challenge.ID,
		Limit:       3,
		Offset:      0,
	}

	challengeEntries, err := testQueries.ListChallengeEntriesByChallengeId(context.Background(), arg)
	require.NoError(t, err)
	require.LessOrEqual(t, len(challengeEntries), int(arg.Limit))

	// Verify all returned entries belong to the challenge
	for _, entry := range challengeEntries {
		require.Equal(t, challenge.ID, entry.ChallengeID)
		require.NotEmpty(t, entry.ID)
		require.True(t, entry.Date.Valid)
	}
}

func TestListChallengeEntriesByChallengeIdEmpty(t *testing.T) {
	challenge := createRandomChallenge(t)

	arg := ListChallengeEntriesByChallengeIdParams{
		ChallengeID: challenge.ID,
		Limit:       5,
		Offset:      0,
	}

	entries, err := testQueries.ListChallengeEntriesByChallengeId(context.Background(), arg)
	require.NoError(t, err)
	require.Empty(t, entries)
}

func TestListChallengeEntriesByChallengeIdWithPagination(t *testing.T) {
	challenge := createRandomChallenge(t)

	// Create multiple entries for pagination testing
	entryCount := 8
	for i := 0; i < entryCount; i++ {
		_, err := testQueries.CreateChallengeEntry(context.Background(), challenge.ID)
		require.NoError(t, err)
	}

	// Test first page
	arg1 := ListChallengeEntriesByChallengeIdParams{
		ChallengeID: challenge.ID,
		Limit:       3,
		Offset:      0,
	}
	firstPage, err := testQueries.ListChallengeEntriesByChallengeId(context.Background(), arg1)
	require.NoError(t, err)
	require.LessOrEqual(t, len(firstPage), 3)

	// Test second page
	arg2 := ListChallengeEntriesByChallengeIdParams{
		ChallengeID: challenge.ID,
		Limit:       3,
		Offset:      3,
	}
	secondPage, err := testQueries.ListChallengeEntriesByChallengeId(context.Background(), arg2)
	require.NoError(t, err)
	require.LessOrEqual(t, len(secondPage), 3)

	// Verify pages don't contain the same entries (if both have entries)
	if len(firstPage) > 0 && len(secondPage) > 0 {
		require.NotEqual(t, firstPage[0].ID, secondPage[0].ID)
	}
}

func TestListChallengeEntriesByChallengeIdOrdering(t *testing.T) {
	challenge := createRandomChallenge(t)

	// Create multiple entries
	entryCount := 5
	for i := 0; i < entryCount; i++ {
		_, err := testQueries.CreateChallengeEntry(context.Background(), challenge.ID)
		require.NoError(t, err)
	}

	arg := ListChallengeEntriesByChallengeIdParams{
		ChallengeID: challenge.ID,
		Limit:       10,
		Offset:      0,
	}

	entries, err := testQueries.ListChallengeEntriesByChallengeId(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)

	// Verify entries are ordered by date (ascending)
	for i := 1; i < len(entries); i++ {
		// Both dates should be valid for comparison
		require.True(t, entries[i-1].Date.Valid)
		require.True(t, entries[i].Date.Valid)
		// Current entry date should be >= previous entry date
		require.True(t, entries[i].Date.Time.After(entries[i-1].Date.Time) || entries[i].Date.Time.Equal(entries[i-1].Date.Time))
	}
}

func TestListChallengesByUserOrdering(t *testing.T) {
	user := createRandomUser(t)

	// Create challenges with different start dates
	challengeCount := 3
	for i := 0; i < challengeCount; i++ {
		arg := CreateChallengeParams{
			Title:       "Challenge " + string(rune('A'+i)),
			UserID:      user.ID,
			Description: pgtype.Text{String: "Description", Valid: true},
			EndDate:     util.GetRandomEndDate(30),
		}
		_, err := testQueries.CreateChallenge(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListChallengesByUserParams{
		UserID: user.ID,
		Limit:  10,
		Offset: 0,
	}

	challenges, err := testQueries.ListChallengesByUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, challenges)

	// Verify challenges are ordered by start_date (ascending)
	for i := 1; i < len(challenges); i++ {
		// Both start dates should be valid for comparison
		require.True(t, challenges[i-1].StartDate.Valid)
		require.True(t, challenges[i].StartDate.Valid)
		// Current challenge start date should be >= previous challenge start date
		require.True(t, challenges[i].StartDate.Time.After(challenges[i-1].StartDate.Time) || challenges[i].StartDate.Time.Equal(challenges[i-1].StartDate.Time))
	}
}
