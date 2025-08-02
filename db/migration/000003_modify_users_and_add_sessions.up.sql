CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid  NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL, 
  "client_ip" varchar NOT NULL, 
  "is_blocked" boolean NOT NULL DEFAULT false, 
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "expires_at" timestamp NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id"); 

ALTER TABLE "users" ADD COLUMN "is_verified" BOOLEAN NOT NULL DEFAULT FALSE;