-- Tables Definition
CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar UNIQUE NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp,
  "timezone" varchar NOT NULL DEFAULT 'UTC'
);

CREATE TABLE "challenges" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "user_id" uuid NOT NULL, 
  "description" varchar,
  "start_date" date NOT NULL DEFAULT (now()::date),
  "end_date" date,
  "active" boolean DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "challenge_entries" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "challenge_id" uuid NOT NULL, 
  "date" date DEFAULT (now()::date),
  "completed" boolean DEFAULT false,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "task_days" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL, 
  "date" date DEFAULT (now()::date),
  "count" integer DEFAULT 0,
  "total_duration" interval DEFAULT (INTERVAL '0 seconds')
);

CREATE TABLE "tasks" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "task_day_id" uuid NOT NULL, 
  "title" varchar NOT NULL,
  "description" varchar,
  "duration" interval NOT NULL,
  "completed" boolean DEFAULT false
);

-- View Definition
CREATE OR REPLACE VIEW current_challenges_view AS
SELECT
  c.id,
  c.title,
  c.user_id,
  c.description,
  c.start_date,
  c.end_date,
  c.active,
  c.created_at,
  (CASE
    WHEN c.end_date IS NOT NULL THEN (c.end_date - c.start_date) + 1
    ELSE (CURRENT_DATE - c.start_date) + 1
  END) AS duration
FROM
  challenges AS c;

-- Indexes
CREATE INDEX ON "challenges" ("start_date");
CREATE INDEX ON "challenges" ("user_id");
CREATE INDEX ON "challenge_entries" ("date");
CREATE INDEX ON "challenge_entries" ("challenge_id");
CREATE INDEX ON "task_days" ("date");
CREATE INDEX ON "task_days" ("user_id");
CREATE INDEX ON "tasks" ("task_day_id");

-- Constraints and Foreign Keys
ALTER TABLE "task_days" ADD CONSTRAINT "user_id_date_unique" UNIQUE ("user_id", "date");

ALTER TABLE "challenges" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "challenge_entries" ADD FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id") ON DELETE CASCADE;
ALTER TABLE "task_days" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "tasks" ADD FOREIGN KEY ("task_day_id") REFERENCES "task_days" ("id") ON DELETE CASCADE;
