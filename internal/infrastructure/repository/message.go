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
	mutex    sync.RWMutex
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
		mutex:    sync.RWMutex{},
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

func (r *MessageRepository) FindConversationsByIDs(
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

func (r *MessageRepository) FindConversationsByUserID(
	ctx context.Context,
	userID entity.ID,
	first int,
	after entity.ID,
) ([]*entity.Conversation, error) {
	// TODO: impl
	return nil, nil
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

func (r *MessageRepository) FindParticipantsInConversations(ctx context.Context,
	conversationIDs []entity.ID) (map[entity.ID][]*entity.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT u.id, u.name, u.picture_url, u.firebase_id, u.provider, u.email_address, u.email_verified, p.conversation_id
		 FROM users AS u
		 INNER JOIN participants AS p ON u.id = p.user_id
		 WHERE p.conversation_id = ANY($1)`,
		pq.Array(conversationIDs),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type participantWithConversationID struct {
		ConversationID entity.ID
		User           model.User
	}

	var participants []participantWithConversationID

	for rows.Next() {
		var participant participantWithConversationID
		if err := rows.Scan(
			&participant.User.ID,
			&participant.User.Name,
			&participant.User.PictureUrl,
			&participant.User.FirebaseID,
			&participant.User.Provider,
			&participant.User.EmailAddress,
			&participant.User.EmailVerified,
			&participant.ConversationID,
		); err != nil {
			return nil, err
		}

		participants = append(participants, participant)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	result := make(map[entity.ID][]*entity.User)
	for _, p := range participants {
		if _, ok := result[p.ConversationID]; !ok {
			result[p.ConversationID] = make([]*entity.User, 0)
		}

		result[p.ConversationID] = append(result[p.ConversationID], model.ConvertModelUser(&p.User))
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

func (s *MessageRepository) FanoutMessage(
	ctx context.Context,
	message *entity.Message,
) {
	s.mutex.RLock()
	for _, c := range s.msgChans {
		c <- message
	}
	s.mutex.RUnlock()
}

func (r *MessageRepository) FindConversationIDsFromUserIDs(ctx context.Context,
	inputs []entity.UserQueryInput) (map[entity.ID]*entity.IDsConnection, error) {
	// TODO: find a better solution
	res := make(map[entity.ID]*entity.IDsConnection)
	for _, input := range inputs {
		idsConnection, err := r.getConversationIDsFromUserID(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("get conversation ids from user: %w", err)
		}

		res[input.UserID] = idsConnection
	}

	return res, nil
}

func (r *MessageRepository) getConversationIDsFromUserID(ctx context.Context,
	input entity.UserQueryInput) (*entity.IDsConnection, error) {
	if !entity.IsValidConversationsSortByType(string(input.SortBy)) {
		return nil, fmt.Errorf("invalid sortBy: %v", input.SortBy)
	}

	sortColumn := model.GetColumnNameByConversationsSortByType(
		entity.ConversationsSortByType(input.SortBy))
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if input.After == 0 {
		query =
			"SELECT conversation_id, " +
				"FALSE AS has_previous_page, " +

				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) " +
				"  FROM ( " +
				"   SELECT * FROM participants " +
				"   WHERE user_id = $1 " +
				"   ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2 + 1 " +
				"  ) as np " +
				" ) = $2 + 1 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_next_page " +

				"FROM participants " +
				"WHERE user_id = $1 " +
				"ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2"

		rows, err = r.db.QueryContext(ctx, query, input.UserID, input.First)
	} else {
		query =
			"SELECT id, " +
				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) FROM participants " +
				"   WHERE " + sortColumn + " <= (SELECT " + sortColumn + " FROM participants WHERE id = $3) " +
				"   AND id != $3 " +
				"   AND id NOT IN " +
				"    (SELECT id FROM participants " +
				"      WHERE " + sortColumn + " = (SELECT " + sortColumn + " FROM participants WHERE id = $3 ) " +
				"      AND id >= $3 ) " +
				"   AND user_id = $1 " +
				" ) > 0 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_previous_page, " +

				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) " +
				"  FROM ( " +
				"   SELECT * FROM participants " +
				"   WHERE " + sortColumn + " >= (SELECT " + sortColumn + " FROM participants WHERE id = $3 ) " +
				"   AND id != $3 " +
				"   AND user_id = $1 " +
				"   ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2 + 1 " +
				"  ) as np " +
				" ) = $2 + 1 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_next_page " +

				"FROM participants " +
				"WHERE " + sortColumn + " >= (SELECT " + sortColumn + " FROM participants WHERE id = $3 ) " +
				"AND id != $3 " +
				"AND id NOT IN ( " +
				" SELECT id FROM participants " +
				" WHERE " + sortColumn + " = (SELECT " + sortColumn + " FROM participants WHERE id = $3 ) " +
				" AND id <= $3 " +
				") " +
				"AND user_id = $1 " +
				"ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2"

		rows, err = r.db.QueryContext(ctx, query, input.UserID, input.First, input.After)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		idEdges         []*entity.IDsEdge
		hasNextPage     bool
		hasPreviousPage bool
	)

	for rows.Next() {
		var edge struct {
			ID              entity.ID `json:"id"`
			HasNextPage     bool      `json:"has_next_page"`
			HasPreviousPage bool      `json:"has_previous_page"`
		}

		if err := rows.Scan(
			&edge.ID,
			&edge.HasNextPage,
			&edge.HasPreviousPage,
		); err != nil {
			return nil, err
		}

		hasNextPage = edge.HasNextPage
		hasPreviousPage = edge.HasPreviousPage
		idEdges = append(idEdges, &entity.IDsEdge{
			Node:   edge.ID,
			Cursor: edge.ID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &entity.IDsConnection{
		Edges: idEdges,
		PageInfo: &entity.PageInfo{
			HasPreviousPage: hasPreviousPage,
			HasNextPage:     hasNextPage,
		},
		TotalCount: len(idEdges),
	}, nil
}
