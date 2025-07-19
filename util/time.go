package util

import "github.com/jackc/pgx/v5/pgtype"

func MinutesToPGInterval(minutes int) pgtype.Interval {
	const microsecondsPerMinute int64 = 60 * 1000000

	return pgtype.Interval{
		Microseconds: int64(minutes) * microsecondsPerMinute,
		Days:         0,
		Months:       0,
		Valid:        true,
	}
}
