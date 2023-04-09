package repository

import (
	"context"
	"errors"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
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

func (mr *chatRepository) CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*domain.Message, error) {
	conn, err := mr.DB.Acquire(ctx)
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
		chatId, err = mr.createChat(trx, ctx, *dto)
		if err != nil {
			return nil, err
		}
	}

	msg, err := mr.createMessage(trx, ctx, dto.SourceUserId, chatId, dto.Text)
	if err != nil {
		return nil, err
	}

	err = mr.saveLastMsgToChat(trx, ctx, chatId, msg.Id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (mr *chatRepository) saveLastMsgToChat(
	trx pgx.Tx,
	ctx context.Context,
	chatId int,
	messageId int,
) error {
	_, err := trx.Exec(ctx, "UPDATE chats SET last_message_id = $1 WHERE id = $2;", messageId, chatId)
	return err
}

func (mr *chatRepository) createChat(
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

func (mr *chatRepository) createMessage(
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
