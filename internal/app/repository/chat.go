package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/antnzr/chat-go/internal/app/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type chatRepository struct {
	DB *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) domain.ChatRepository {
	return &chatRepository{
		DB: db,
	}
}

func (cr *chatRepository) CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*domain.Message, error) {
	conn, err := cr.DB.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	trx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = trx.Rollback(ctx)
		} else {
			_ = trx.Commit(ctx)
		}
	}()

	findChatIdQuery := `
		SELECT uc.chat_id FROM user_chats AS uc
		WHERE uc.user_id = $1 OR uc.user_id = $2
		GROUP BY uc.chat_id
		HAVING COUNT(*) = 2;
	`
	var chatId int
	err = trx.QueryRow(ctx, findChatIdQuery, dto.SourceUserId, dto.TargetUserId).Scan(&chatId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if chatId == 0 {
		chatId, err = cr.createChat(trx, ctx, *dto)
		if err != nil {
			return nil, err
		}
	}

	msg, err := cr.createMessage(trx, ctx, dto.SourceUserId, chatId, dto.Text)
	if err != nil {
		return nil, err
	}

	err = cr.saveLastMsgToChat(trx, ctx, chatId, msg.Id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (cr *chatRepository) saveLastMsgToChat(
	trx pgx.Tx,
	ctx context.Context,
	chatId int,
	messageId int,
) error {
	_, err := trx.Exec(ctx, "UPDATE chats SET last_message_id = $1 WHERE id = $2;", messageId, chatId)
	return err
}

func (cr *chatRepository) createChat(
	trx pgx.Tx,
	ctx context.Context,
	dto dto.SendMessageRequest,
) (int, error) {
	var chatId int
	err := trx.QueryRow(ctx, "INSERT INTO chats VALUES(DEFAULT) RETURNING id;").Scan(&chatId)
	if err != nil {
		return 0, err
	}

	_, err = trx.CopyFrom(
		ctx,
		pgx.Identifier{"user_chats"},
		[]string{"user_id", "chat_id"},
		pgx.CopyFromRows([][]interface{}{
			{dto.SourceUserId, chatId},
			{dto.TargetUserId, chatId},
		}),
	)

	if err != nil {
		return 0, err
	}

	return chatId, nil
}

func (cr *chatRepository) createMessage(
	trx pgx.Tx,
	ctx context.Context,
	ownerId int,
	chatId int,
	text string,
) (*domain.Message, error) {
	message := new(domain.Message)
	row := trx.QueryRow(
		ctx,
		`INSERT INTO "messages" ("owner_id", "chat_id", "text")
		VALUES($1, $2, $3)
		RETURNING "id", "owner_id", "chat_id", "text", "created_at";`,
		ownerId, chatId, text,
	)
	err := message.ScanRow(row)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (cr *chatRepository) FindChats(
	ctx context.Context,
	searchQuery dto.ChatSearchQuery,
) (int, []dto.ChatResponse, error) {
	var (
		fields = []string{}
		args   = []any{}
	)

	if searchQuery.UserId != 0 {
		fields = append(fields, " uc.user_id = $1")
		args = append(args, searchQuery.UserId)
	}

	var where string
	if len(fields) > 0 {
		where = " WHERE " + strings.Join(fields, " AND ")
	}

	// skip chats without last message
	totalQuery := fmt.Sprintf(`
		SELECT COUNT(*) AS total
		FROM user_chats AS uc
		JOIN chats AS c
			ON c.id = uc.chat_id
		JOIN messages AS m
			ON c.last_message_id = m.id
		%s;
	`, where)
	var total int
	err := cr.DB.QueryRow(ctx, totalQuery, args...).Scan(&total)

	if err != nil {
		return 0, nil, errs.ClarifyError(err)
	}

	if total == 0 {
		return 0, nil, nil
	}

	sql := fmt.Sprintf(`
		SELECT c.id, c.name, c.description, c.last_message_id, c.created_at,
				 m.id AS message_id, m.owner_id AS message_owner, m.text AS message_text,
				 m.chat_id AS chat_id, m.created_at AS message_created_at
		FROM user_chats AS uc
		JOIN chats AS c
			ON c.id = uc.chat_id
		JOIN messages AS m
			ON c.last_message_id = m.id
		%s
		ORDER BY c.created_at DESC
	`, where)

	args = append(args, searchQuery.Limit)
	sql += fmt.Sprintf(` LIMIT $%d`, len(args))

	args = append(args, (searchQuery.Page-1)*searchQuery.Limit)
	sql += fmt.Sprintf(` OFFSET $%d`, len(args))

	rows, err := cr.DB.Query(ctx, sql, args...)

	if err != nil {
		return 0, nil, errs.ClarifyError(err)
	}
	defer rows.Close()

	chats, err := scan(rows)
	if err != nil {
		return 0, nil, errs.ClarifyError(err)
	}

	return total, chats, nil
}

func scan(rows pgx.Rows) ([]dto.ChatResponse, error) {
	var chats []dto.ChatResponse
	for rows.Next() {
		chat, err := scanRowsIntoChat(rows)
		if err != nil {
			return nil, err
		}
		chats = append(chats, *chat)
	}
	return chats, nil
}

func scanRowsIntoChat(rows pgx.Rows) (*dto.ChatResponse, error) {
	chat := new(dto.ChatResponse)
	chat.LastMessage = new(dto.MessageResponse)

	err := rows.Scan(
		&chat.Id,
		&chat.Name,
		&chat.Description,
		&chat.LastMessageId,
		&chat.CreatedAt,
		&chat.LastMessage.Id,
		&chat.LastMessage.OwnerId,
		&chat.LastMessage.Text,
		&chat.LastMessage.ChatId,
		&chat.LastMessage.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return chat, nil
}
