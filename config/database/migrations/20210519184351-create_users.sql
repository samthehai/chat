-- +migrate Up
CREATE TABLE IF NOT EXISTS users(
  id SERIAL,
  name VARCHAR (255) NOT NULL,
  firebase_id VARCHAR (255) UNIQUE NOT NULL,
  provider VARCHAR (255) NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT users_pk_id PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS users;
