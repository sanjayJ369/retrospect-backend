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


-- name: UpdateTaskDayTotalsForUsersInTimezone :exec
-- For a given timezone, finds all task_days from "yesterday" and updates their totals.
-- This is much more efficient than fetching IDs into Go and updating one by one.
WITH daily_summary AS (
  SELECT
    td.id as task_day_id,
    SUM(t.duration) as total_duration,
    COUNT(t.id) as task_count
  FROM tasks t
  JOIN task_days td ON t.task_day_id = td.id
  JOIN users u ON td.user_id = u.id
  WHERE
    u.timezone = $1
    AND td.date = (NOW() AT TIME ZONE $1)::date - INTERVAL '1 day' 
  GROUP BY td.id
)
UPDATE task_days
SET
  total_duration = daily_summary.total_duration,
  count = daily_summary.task_count
FROM daily_summary
WHERE task_days.id = daily_summary.task_day_id;