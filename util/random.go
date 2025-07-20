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
