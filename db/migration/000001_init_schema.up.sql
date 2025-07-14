CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "email" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp
);

CREATE TABLE "challenges" (
  "id" uuid PRIMARY KEY,
  "title" varchar NOT NULL,
  "user_id" uuid,
  "description" varchar,
  "start_date" date NOT NULL,
  "end_date" date,
  "duration" integer,
  "active" boolean,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "challenge_entries" (
  "id" uuid PRIMARY KEY,
  "challenge_id" uuid,
  "user_id" uuid,
  "date" date,
  "completed" boolean,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "task_days" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid,
  "date" date,
  "count" integer,
  "total_duration" time
);

CREATE TABLE "tasks" (
  "id" uuid PRIMARY KEY,
  "task_day_id" uuid,
  "user_id" uuid,
  "title" varchar NOT NULL,
  "description" varchar,
  "duration" time NOT NULL,
  "completed" boolean
);

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "challenges" ("start_date");

CREATE INDEX ON "challenge_entries" ("date");

CREATE INDEX ON "task_days" ("date");

ALTER TABLE "challenges" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "challenge_entries" ADD FOREIGN KEY ("challenge_id") REFERENCES "challenges" ("id");

ALTER TABLE "challenge_entries" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "task_days" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("task_day_id") REFERENCES "task_days" ("id");

ALTER TABLE "tasks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");