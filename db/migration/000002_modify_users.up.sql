ALTER TABLE "users"
ADD COLUMN "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
ADD COLUMN "hashed_password" varchar NOT NULL;