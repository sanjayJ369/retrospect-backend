-- name: GetTimezonesWhereDayIsStarting :many
-- Finds all distinct timezones where the current local time is at the beginning of a new day (e.g., 00:00 to 00:59).
-- We can then find all users for these timezones.
SELECT DISTINCT timezone FROM users
WHERE EXTRACT(HOUR FROM (NOW() AT TIME ZONE timezone)) = 0;

-- name: CreateTaskDaysForUsersInTimezone :exec
-- Creates new task_day entries for all users in a given timezone for their "today".
-- The `(NOW() AT TIME ZONE $1)::date` correctly calculates the user's current date.
INSERT INTO task_days (user_id, date)
SELECT id, (NOW() AT TIME ZONE $1)::date
FROM users
WHERE timezone = $1
ON CONFLICT (user_id, date) DO NOTHING;


