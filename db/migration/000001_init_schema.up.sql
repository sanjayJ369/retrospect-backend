CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar NOT NULL,
  "name" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp
);

CREATE TABLE "challenges" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "title" varchar NOT NULL,
  "user_id" uuid,
  "description" varchar,
  "start_date" date NOT NULL DEFAULT (now()::date),
  "end_date" date,
  "duration" integer,
  "active" boolean DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "challenge_entries" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "challenge_id" uuid,
  "user_id" uuid,
  "date" date DEFAULT (now()::date),
  "completed" boolean DEFAULT false,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "task_days" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid,
  "date" date DEFAULT (now()::date),
  "count" integer DEFAULT 0,
  "total_duration" time
);

CREATE TABLE "tasks" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "task_day_id" uuid,
  "user_id" uuid,
  "title" varchar NOT NULL,
  "description" varchar,
  "duration" time NOT NULL,
  "completed" boolean DEFAULT false
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
