-- +migrate Up
CREATE TABLE IF NOT EXISTS messages(
  id SERIAL NOT NULL,
  conversation_id INTEGER NOT NULL,
  sender_id INTEGER NOT NULL,
  type TEXT CHECK (type IN ('text')) NOT NULL,
  content TEXT NOT NULL DEFAULT '',
  --
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ DEFAULT NULL,
  --
  CONSTRAINT messages_pk_id PRIMARY KEY (id),
  CONSTRAINT messages_fk_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversations (id),
  CONSTRAINT messages_fk_sender_id FOREIGN KEY (sender_id) REFERENCES users (id)
);
-- +migrate Down
DROP TABLE IF EXISTS messages;
