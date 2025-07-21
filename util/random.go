package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// generates a random number between min and max including both
func GetRandomInt(min int, max int) int64 {
	return int64(min) + rand.Int63n(int64(max-min+1))
}

var alphabets = "abcdefghijklmnopqrst"

func GetRandomString(l int) string {
	var strbul strings.Builder
	for i := 0; i < l; i++ {
		rnum := GetRandomInt(0, len(alphabets)-1)
		rchar := alphabets[rnum]
		strbul.WriteByte(rchar)
	}
	return strbul.String()
}

func GetUUIDPGType() pgtype.UUID {
	newID := uuid.New()
	return pgtype.UUID{
		Bytes: newID,
		Valid: true,
	}
}

var timezones = []string{
	"UTC",
	"America/New_York",
	"America/Los_Angeles",
	"Europe/London",
	"Europe/Paris",
	"Asia/Tokyo",
	"Asia/Dubai",
	"Asia/Kolkata",
	"Australia/Sydney",
	"Pacific/Honolulu",
}

func GetRandomTimezone() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return timezones[r.Intn(len(timezones))]
}

// getRandomEndDate returns a random date between current day
// and current day + max or returns nulltime randomly
func GetRandomEndDate(max int) pgtype.Date {

	p := GetRandomInt(0, 101)
	alpha := 30
	if int(p) > alpha {
		nowUTC := time.Now().UTC()
		todayUTC := nowUTC.Truncate(24 * time.Hour)
		daysToAdd := GetRandomInt(1, max)
		endDate := todayUTC.AddDate(0, 0, int(daysToAdd))

		return pgtype.Date{
			Time:  endDate,
			Valid: true,
		}
	}

	// Return a NULL date
	return pgtype.Date{
		Valid: false,
	}
}
