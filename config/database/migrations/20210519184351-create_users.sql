-- +migrate Up
CREATE TABLE IF NOT EXISTS users(
  id SERIAL,
  name VARCHAR (255) DEFAULT '',
  picture_url VARCHAR (255) DEFAULT '',
  email_address VARCHAR (255) DEFAULT '',
  email_verified BOOLEAN DEFAULT FALSE,
  firebase_id VARCHAR (255) UNIQUE NOT NULL,
  provider VARCHAR (255) DEFAULT '',

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT users_pk_id PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
