-- +migrate Up
CREATE TABLE IF NOT EXISTS users(
  id SERIAL NOT NULL,
  name TEXT NOT NULL,
  picture_url TEXT DEFAULT '',
  email_address TEXT DEFAULT '',
  email_verified BOOLEAN DEFAULT FALSE,
  firebase_id TEXT UNIQUE NOT NULL,
  provider TEXT DEFAULT '',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT users_pk_id PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
