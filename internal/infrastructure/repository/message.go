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
	cacher       external.Cacher
	msgChans     map[entity.ID]chan *entity.Message
	mutex        sync.RWMutex
	dbTransactor external.Transactor
	db           *sql.DB
}

func NewMessageRepository(
	cacher external.Cacher,
	dbTransactor external.Transactor,
	db *sql.DB,
) *MessageRepository {
	return &MessageRepository{
		cacher:       cacher,
		dbTransactor: dbTransactor,
		db:           db,
		msgChans:     map[entity.ID]chan *entity.Message{},
		mutex:        sync.RWMutex{},
	}
}

func (r *MessageRepository) CreateConversationWithTransaction(
	ctx context.Context, creatorID entity.ID, conversationTitle string,
	conversationType entity.ConversationType, recipentIDs []entity.ID) (
	conversationID *entity.ID, err error) {
	fail := func(err error) (*entity.ID, error) {
		return nil, fmt.Errorf("CreateConversationWithTransaction: %w", err)
	}

	tx, ok := r.dbTransactor.GetTransactionFromCtx(ctx)
	if !ok {
		return fail(fmt.Errorf("get transaction from ctx failed"))
	}

	conversationID, err = r.createConversation(ctx, tx, creatorID,
		conversationTitle, conversationType, recipentIDs)
	if err != nil {
		return fail(err)
	}

	return conversationID, nil
}

func (r *MessageRepository) CreateConversation(
	ctx context.Context,
	creatorID entity.ID,
	conversationTitle string,
	conversationType entity.ConversationType,
	recipentIDs []entity.ID,
) (conversationID *entity.ID, err error) {
	fail := func(err error) (*entity.ID, error) {
		return nil, fmt.Errorf("CreateConversation: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fail(fmt.Errorf("begin transaction: %w", err))
	}
	defer tx.Rollback()

	conversationID, err = r.createConversation(ctx, tx, creatorID, conversationTitle,
		conversationType, recipentIDs)
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return conversationID, nil
}

func (r *MessageRepository) createConversation(
	ctx context.Context,
	tx *sql.Tx,
	creatorID entity.ID,
	conversationTitle string,
	conversationType entity.ConversationType,
	recipentIDs []entity.ID,
) (*entity.ID, error) {
	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO conversations (creator_id, title, type) VALUES ($1,$2,$3) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("prepare context: %w", err)
	}
	defer stmt.Close()

	var createdID entity.ID
	err = stmt.QueryRowContext(ctx, creatorID, conversationTitle, conversationType).Scan(&createdID)
	if err != nil {
		return nil, fmt.Errorf("exec context: %w", err)
	}

	if err := r.createParticipants(ctx, tx, createdID, recipentIDs); err != nil {
		return nil, fmt.Errorf("create participants: %w", err)
	}

	return &createdID, nil
}

func (r *MessageRepository) createParticipants(
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

func (r *MessageRepository) FindConversationsByIDsWithTransaction(
	ctx context.Context,
	conversationIDs []entity.ID,
) ([]*entity.Conversation, error) {
	fail := func(err error) ([]*entity.Conversation, error) {
		return nil, fmt.Errorf("FindConversationsByIDsWithTransaction: %w", err)
	}

	tx, ok := r.dbTransactor.GetTransactionFromCtx(ctx)
	if !ok {
		return fail(fmt.Errorf("get transaction from ctx failed"))
	}

	conversations, err := r.findConversationsByIDs(ctx, tx, conversationIDs)
	if err != nil {
		return fail(err)
	}

	return conversations, nil
}

func (r *MessageRepository) FindConversationsByIDs(
	ctx context.Context,
	conversationIDs []entity.ID,
) ([]*entity.Conversation, error) {
	fail := func(err error) ([]*entity.Conversation, error) {
		return nil, fmt.Errorf("FindConversationsByIDs: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fail(fmt.Errorf("begin transaction: %w", err))
	}
	defer tx.Rollback()

	conversations, err := r.findConversationsByIDs(ctx, tx, conversationIDs)
	if err != nil {
		return fail(err)
	}

	return conversations, nil
}

func (r *MessageRepository) findConversationsByIDs(
	ctx context.Context,
	tx *sql.Tx,
	conversationIDs []entity.ID,
) ([]*entity.Conversation, error) {
	rows, err := tx.QueryContext(
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

func (r *MessageRepository) CreateMessageWithTransaction(
	ctx context.Context,
	conversationID entity.ID,
	msgType entity.MessageType,
	senderID entity.ID,
	msg string,
) (*entity.Message, error) {
	fail := func(err error) (*entity.Message, error) {
		return nil, fmt.Errorf("CreateMessageWithTransaction: %w", err)
	}

	tx, ok := r.dbTransactor.GetTransactionFromCtx(ctx)
	if !ok {
		return nil, fmt.Errorf("get transaction from ctx failed")
	}

	message, err := r.createMessage(ctx, tx, conversationID, msgType, senderID, msg)
	if err != nil {
		return fail(err)
	}

	return message, nil
}

func (r *MessageRepository) CreateMessage(
	ctx context.Context,
	conversationID entity.ID,
	msgType entity.MessageType,
	senderID entity.ID,
	msg string,
) (*entity.Message, error) {
	fail := func(err error) (*entity.Message, error) {
		return nil, fmt.Errorf("CreateMessage: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fail(fmt.Errorf("begin transaction: %w", err))
	}
	defer tx.Rollback()

	message, err := r.createMessage(ctx, tx, conversationID, msgType, senderID, msg)
	if err != nil {
		return fail(err)
	}

	return message, nil
}

func (r *MessageRepository) createMessage(
	ctx context.Context,
	tx *sql.Tx,
	conversationID entity.ID,
	msgType entity.MessageType,
	senderID entity.ID,
	msg string,
) (*entity.Message, error) {
	stmt, err := tx.PrepareContext(
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

func (r *MessageRepository) FindAllMessagesInConversations(
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
	inputs []entity.RelayQueryInput) (map[entity.ID]*entity.IDsConnection, error) {
	// TODO: find a better solution
	res := make(map[entity.ID]*entity.IDsConnection)
	for _, input := range inputs {
		idsConnection, err := r.getConversationIDsFromUserID(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("get conversation ids from user: %w", err)
		}

		res[input.KeyID] = idsConnection
	}

	return res, nil
}

func (r *MessageRepository) getConversationIDsFromUserID(ctx context.Context,
	input entity.RelayQueryInput) (*entity.IDsConnection, error) {
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

		rows, err = r.db.QueryContext(ctx, query, input.KeyID, input.First)
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

		rows, err = r.db.QueryContext(ctx, query, input.KeyID, input.First, input.After)
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

func (r *MessageRepository) FindMessagesInConversations(ctx context.Context,
	inputs []entity.RelayQueryInput) (map[entity.ID]*entity.ConversationMessagesConnection, error) {
	// TODO: find a better solution
	res := make(map[entity.ID]*entity.ConversationMessagesConnection)
	for _, input := range inputs {
		cmc, err := r.findMessagesInConversation(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("find conversation ids from user: %w", err)
		}

		res[input.KeyID] = cmc
	}

	return res, nil
}

func (r *MessageRepository) findMessagesInConversation(ctx context.Context,
	input entity.RelayQueryInput) (*entity.ConversationMessagesConnection, error) {
	if !entity.IsValidMessagesSortByType(string(input.SortBy)) {
		return nil, fmt.Errorf("invalid sortBy: %v", input.SortBy)
	}

	sortColumn := model.GetColumnNameByMessagesSortByType(
		entity.MessagesSortByType(input.SortBy))
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if input.After == 0 {
		query =
			"SELECT id, conversation_id, sender_id, type, content, created_at, updated_at, deleted_at, " +
				"FALSE AS has_previous_page, " +

				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) " +
				"  FROM ( " +
				"   SELECT * FROM messages " +
				"   WHERE conversation_id = $1 " +
				"   ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2 + 1 " +
				"  ) as np " +
				" ) = $2 + 1 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_next_page " +

				"FROM messages " +
				"WHERE conversation_id = $1 " +
				"ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2"

		rows, err = r.db.QueryContext(ctx, query, input.KeyID, input.First)
	} else {
		query =
			"SELECT id, conversation_id, sender_id, type, content, created_at, updated_at, deleted_at, " +
				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) FROM messages " +
				"   WHERE " + sortColumn + " <= (SELECT " + sortColumn + " FROM messages WHERE id = $3) " +
				"   AND id != $3 " +
				"   AND id NOT IN " +
				"    (SELECT id FROM messages " +
				"      WHERE " + sortColumn + " = (SELECT " + sortColumn + " FROM messages WHERE id = $3 ) " +
				"      AND id >= $3 ) " +
				"   AND conversation_id = $1 " +
				" ) > 0 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_previous_page, " +

				"CASE " +
				" WHEN ( " +
				"  SELECT COUNT(*) " +
				"  FROM ( " +
				"   SELECT * FROM messages " +
				"   WHERE " + sortColumn + " >= (SELECT " + sortColumn + " FROM messages WHERE id = $3 ) " +
				"   AND id != $3 " +
				"   AND conversation_id = $1 " +
				"   ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2 + 1 " +
				"  ) as np " +
				" ) = $2 + 1 " +
				"THEN TRUE ELSE FALSE " +
				"END AS has_next_page " +

				"FROM messages " +
				"WHERE " + sortColumn + " >= (SELECT " + sortColumn + " FROM messages WHERE id = $3 ) " +
				"AND id != $3 " +
				"AND id NOT IN ( " +
				" SELECT id FROM messages " +
				" WHERE " + sortColumn + " = (SELECT " + sortColumn + " FROM messages WHERE id = $3 ) " +
				" AND id <= $3 " +
				") " +
				"AND conversation_id = $1 " +
				"ORDER BY " + sortColumn + " ASC, id ASC LIMIT $2"

		rows, err = r.db.QueryContext(ctx, query, input.KeyID, input.First, input.After)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		cmEdges         []*entity.ConversationMessagesEdge
		hasNextPage     bool
		hasPreviousPage bool
	)

	for rows.Next() {
		var edge struct {
			Message         entity.Message
			HasNextPage     bool `json:"has_next_page"`
			HasPreviousPage bool `json:"has_previous_page"`
		}

		if err := rows.Scan(
			&edge.Message.ID,
			&edge.Message.ConversationID,
			&edge.Message.SenderID,
			&edge.Message.Type,
			&edge.Message.Content,
			&edge.Message.CreatedAt,
			&edge.Message.UpdatedAt,
			&edge.Message.DeletedAt,
			&edge.HasNextPage,
			&edge.HasPreviousPage,
		); err != nil {
			return nil, err
		}

		hasNextPage = edge.HasNextPage
		hasPreviousPage = edge.HasPreviousPage
		cmEdges = append(cmEdges, &entity.ConversationMessagesEdge{
			Node:   &edge.Message,
			Cursor: edge.Message.ID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &entity.ConversationMessagesConnection{
		Edges: cmEdges,
		PageInfo: &entity.PageInfo{
			HasPreviousPage: hasPreviousPage,
			HasNextPage:     hasNextPage,
		},
		TotalCount: len(cmEdges),
	}, nil
}
