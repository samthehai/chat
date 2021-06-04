-- +migrate Up
CREATE TABLE IF NOT EXISTS conversations(
  id SERIAL NOT NULL,
  title TEXT DEFAULT '',
  creator_id INTEGER DEFAULT NULL,
  --
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  --
  CONSTRAINT conversations_pk_id PRIMARY KEY (id),
  CONSTRAINT conversations_fk_creator_id FOREIGN KEY (creator_id) REFERENCES users (id)
);
-- +migrate Down
DROP TABLE IF EXISTS conversations;
