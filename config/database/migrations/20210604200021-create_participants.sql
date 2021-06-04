
-- +migrate Up
CREATE TABLE IF NOT EXISTS participants(
  id SERIAL NOT NULL,
  conversation_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  type TEXT CHECK (type IN ('text')) NOT NULL,
  --
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  --
  CONSTRAINT participants_pk_id PRIMARY KEY (id),
  CONSTRAINT participants_fk_conversation_id FOREIGN KEY (conversation_id) REFERENCES conversations (id),
  CONSTRAINT participants_fk_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +migrate Down
DROP TABLE IF EXISTS participants;
