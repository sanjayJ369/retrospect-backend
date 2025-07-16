package db

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/sanjayj369/retrospect-backend/util"
	"github.com/stretchr/testify/require"
)

// createRandomChallenge create a new user and inserts
// a challenge for the user
func createRandomChallenge(t testing.TB) Challenge {
	t.Helper()
	user := createRandomUser(t)
	arg := CreateChallengeParams{
		Title:  util.GetRandomString(10),
		UserID: user.ID,
		Description: sql.NullString{
			String: util.GetRandomString(100),
			Valid:  true,
		},
		EndDate: sql.NullTime{},
	}
	challenge, err := testQueries.CreateChallenge(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Title, challenge.Title)
	require.Equal(t, arg.UserID, challenge.UserID)
	require.Equal(t, arg.Description, challenge.Description)
	require.Equal(t, arg.EndDate, challenge.EndDate)
	return challenge
}

func TestCreateChallenge(t *testing.T) {
	createRandomChallenge(t)
}

func TestGetChallenge(t *testing.T) {
	challenge := createRandomChallenge(t)
	challenge1, err := testQueries.GetChallenge(context.Background(), challenge.ID)
	require.NoError(t, err)
	require.Equal(t, challenge1, challenge)
}

func TestDeleteChallenge(t *testing.T) {
	challenge := createRandomChallenge(t)
	challenge1, err := testQueries.DeleteChallenge(context.Background(), challenge.ID)
	require.NoError(t, err)
	require.Equal(t, challenge1, challenge)

	challenge2, err := testQueries.GetChallenge(context.Background(), challenge.ID)
	require.Error(t, err)
	require.Empty(t, challenge2)
}

func TestListChallenges(t *testing.T) {
	// TODO: update test for more through testing
	count := 10
	for i := 0; i < count; i++ {
		createRandomChallenge(t)
	}

	arg := ListChallengesParams{
		Limit:  2,
		Offset: 2,
	}
	res, err := testQueries.ListChallenges(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, int(arg.Limit), len(res))
}

func TestUpdateChallengeActiveStatus(t *testing.T) {
	challenge := createRandomChallenge(t)
	arg := UpdateChallengeActiveStatusParams{
		ID: challenge.ID,
		Active: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	}
	challenge1, err := testQueries.UpdateChallengeActiveStatus(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, challenge1.Active.Bool, arg.Active.Bool)
}

func TestUpdateChallengeDescription(t *testing.T) {
	challenge := createRandomChallenge(t)
	arg := UpdateChallengeDescriptionParams{
		ID: challenge.ID,
		Description: sql.NullString{
			String: util.GetRandomString(100),
			Valid:  true,
		},
	}
	challenge1, err := testQueries.UpdateChallengeDescription(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, challenge1.Description.String, arg.Description.String)
}

func TestUpdateChallengeDetails(t *testing.T) {
	challenge := createRandomChallenge(t)
	arg := UpdateChallengeDetailsParams{
		ID:    challenge.ID,
		Title: util.GetRandomString(10),
		Description: sql.NullString{
			String: util.GetRandomString(100),
			Valid:  true,
		},
		EndDate: getRandomEndDate(t, 100),
	}
	challenge1, err := testQueries.UpdateChallengeDetails(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Title, challenge1.Title)
	require.Equal(t, arg.Description.String, challenge1.Description.String)
	require.Equal(t, arg.EndDate.Time, arg.EndDate.Time)
}

func TestUpdateChallengeEndDate(t *testing.T) {
	challenge := createRandomChallenge(t)
	arg := UpdateChallengeEndDateParams{
		ID:      challenge.ID,
		EndDate: getRandomEndDate(t, 100),
	}
	challenge1, err := testQueries.UpdateChallengeEndDate(context.Background(), arg)
	require.NoError(t, err)
	expectedTimeInUTC := arg.EndDate.Time.UTC()
	require.Equal(t, expectedTimeInUTC, challenge1.EndDate.Time.UTC())
}

func TestUpdateChallengeTitle(t *testing.T) {
	challenge := createRandomChallenge(t)
	arg := UpdateChallengeTitleParams{
		ID:    challenge.ID,
		Title: util.GetRandomString(10),
	}
	challenge1, err := testQueries.UpdateChallengeTitle(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, arg.Title, challenge1.Title)
}

// getRandomEndDate returns a random date between current day
// and current day + max or returns nulltime randomly
func getRandomEndDate(t testing.TB, max int) sql.NullTime {
	t.Helper()
	p := rand.IntN(101)
	alpha := 30
	if p > alpha {
		nowInLocation := time.Now().In(time.Local)
		today := nowInLocation.Truncate(24 * time.Hour)
		daysToAdd := util.GetRandomInt(1, max)
		endDate := today.Add(time.Duration(daysToAdd) * 24 * time.Hour)
		return sql.NullTime{
			Time:  endDate,
			Valid: true,
		}
	}
	return sql.NullTime{
		Valid: false,
	}
}
