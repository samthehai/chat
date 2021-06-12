package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/lib/pq"
	"github.com/samthehai/chat/internal/domain/entity"
	"github.com/samthehai/chat/internal/infrastructure/repository/external"
	"github.com/samthehai/chat/internal/infrastructure/repository/model"
)

type MessageRepository struct {
	cacher   external.Cacher
	msgChans map[entity.ID]chan *entity.Message
	mutex    sync.Mutex
	db       *sql.DB
}

func NewMessageRepository(
	cacher external.Cacher,
	db *sql.DB,
) *MessageRepository {
	return &MessageRepository{
		cacher:   cacher,
		db:       db,
		msgChans: map[entity.ID]chan *entity.Message{},
		mutex:    sync.Mutex{},
	}
}

func (r *MessageRepository) CreateConversation(
	ctx context.Context,
	creatorID entity.ID,
	conversationTitle string,
	conversationType entity.ConversationType,
	recipentIDs []entity.ID,
) (*entity.ID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction begin: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO conversations (creator_id, title, type) VALUES ($1,$2,$3) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	var createdID entity.ID
	err = stmt.QueryRowContext(ctx, creatorID, conversationTitle, conversationType).Scan(&createdID)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}

	if err := r.CreateParticipants(ctx, tx, createdID, recipentIDs); err != nil {
		return nil, fmt.Errorf("create participants: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &createdID, nil
}

func (r *MessageRepository) CreateParticipants(
	ctx context.Context,
	tx *sql.Tx,
	conversationID entity.ID,
	recipentIDs []entity.ID,
) error {
	vStrs := []string{}
	vArgs := []interface{}{}
	for index, id := range recipentIDs {
		vStrs = append(vStrs, fmt.Sprintf("($%v, $%v)", index+1, index+2))

		vArgs = append(vArgs, conversationID)
		vArgs = append(vArgs, id)
	}

	smt := `INSERT INTO participants(conversation_id, user_id) VALUES %s`
	smt = fmt.Sprintf(smt, strings.Join(vStrs, ","))

	_, err := tx.ExecContext(ctx, smt, vArgs...)
	if err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

func (r *MessageRepository) FindConversations(
	ctx context.Context,
	conversationIDs []entity.ID,
) ([]*entity.Conversation, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, creator_id, title, type, created_at, updated_at, deleted_at
			FROM conversations
			WHERE id = ANY($1)`,
		pq.Array(conversationIDs),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*model.Conversation

	for rows.Next() {
		var c model.Conversation
		if err := rows.Scan(
			&c.ID,
			&c.CreatorID,
			&c.Title,
			&c.Type,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		); err != nil {
			return nil, err
		}

		conversations = append(conversations, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return model.ConvertModelConversations(conversations), nil
}

func (r *MessageRepository) CreateMessage(
	ctx context.Context,
	conversationID entity.ID,
	msgType entity.MessageType,
	senderID entity.ID,
	msg string,
) (*entity.Message, error) {

	stmt, err := r.db.PrepareContext(
		ctx,
		`INSERT INTO messages(conversation_id, sender_id, type, content)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, conversation_id, sender_id, type, content, created_at, updated_at, deleted_at`,
	)
	if err != nil {
		return nil, fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	var message entity.Message
	err = stmt.QueryRowContext(ctx, conversationID, senderID, msgType, msg).
		Scan(
			&message.ID,
			&message.ConversationID,
			&message.SenderID,
			&message.Type,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.DeletedAt,
		)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}

	return &message, nil
}

func (r *MessageRepository) FindMessagesInConversations(
	ctx context.Context,
	conversationIDs []entity.ID,
) (map[entity.ID][]*entity.Message, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, conversation_id, sender_id, type, content, created_at, updated_at, deleted_at
		 FROM messages
		 WHERE conversation_id = ANY($1)`,
		pq.Array(conversationIDs),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message

	for rows.Next() {
		var message model.Message
		if err := rows.Scan(
			&message.ID,
			&message.ConversationID,
			&message.SenderID,
			&message.Type,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.DeletedAt,
		); err != nil {
			return nil, err
		}

		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make(map[entity.ID][]*entity.Message)
	for _, m := range messages {
		if _, ok := result[m.ConversationID]; !ok {
			result[m.ConversationID] = make([]*entity.Message, 0)
		}

		result[m.ConversationID] = append(result[m.ConversationID], model.ConvertModelMessage(m))
	}

	return result, nil
}

func (s *MessageRepository) MessagePosted(
	ctx context.Context,
	input entity.User,
) (<-chan *entity.Message, error) {
	messages := make(chan *entity.Message, 1)
	s.mutex.Lock()
	s.msgChans[input.ID] = messages
	s.mutex.Unlock()

	go func() {
		<-ctx.Done()

		s.mutex.Lock()
		delete(s.msgChans, input.ID)
		s.mutex.Unlock()
	}()

	return messages, nil
}
