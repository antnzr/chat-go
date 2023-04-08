package repository

import (
	"context"

	"github.com/antnzr/chat-go/internal/app/domain"
	"github.com/antnzr/chat-go/internal/app/dto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type messageRepository struct {
	DB *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) domain.MessageRepository {
	return &messageRepository{
		DB: db,
	}
}

func (mr *messageRepository) CreateMessage(ctx context.Context, dto *dto.SendMessageRequest) (*domain.Message, error) {
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
			trx.Rollback(ctx)
		} else {
			trx.Commit(ctx)
		}
	}()

	findDialogIdQuery := `
		SELECT ud.dialog_id FROM user_dialogs AS ud
		WHERE ud.user_id = $1 OR ud.user_id = $2
		GROUP BY ud.dialog_id
		HAVING COUNT(*) = 2;
	`
	var dialogId int
	trx.QueryRow(ctx, findDialogIdQuery, dto.SourceUserId, dto.TargetUserId).Scan(&dialogId)

	if dialogId == 0 {
		dialogId, err = mr.createDialog(trx, ctx, *dto)
		if err != nil {
			return nil, err
		}
	}

	msg, err := mr.createMessage(trx, ctx, dto.SourceUserId, dialogId, dto.Text)
	if err != nil {
		return nil, err
	}

	err = mr.saveLastMsgToDialog(trx, ctx, dialogId, msg.Id)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (mr *messageRepository) saveLastMsgToDialog(
	trx pgx.Tx,
	ctx context.Context,
	dialogId int,
	messageId int,
) error {
	_, err := trx.Exec(ctx, "UPDATE dialogs SET last_message_id = $1 WHERE id = $2;", messageId, dialogId)
	return err
}

func (mr *messageRepository) createDialog(
	trx pgx.Tx,
	ctx context.Context,
	dto dto.SendMessageRequest,
) (int, error) {
	var dialogId int
	trx.QueryRow(ctx, "INSERT INTO dialogs VALUES(DEFAULT) RETURNING id;").Scan(&dialogId)

	_, err := trx.CopyFrom(
		ctx,
		pgx.Identifier{"user_dialogs"},
		[]string{"user_id", "dialog_id"},
		pgx.CopyFromRows([][]interface{}{
			{dto.SourceUserId, dialogId},
			{dto.TargetUserId, dialogId},
		}),
	)

	if err != nil {
		return 0, err
	}

	return dialogId, nil
}

func (mr *messageRepository) createMessage(
	trx pgx.Tx,
	ctx context.Context,
	ownerId int,
	dialogId int,
	text string,
) (*domain.Message, error) {
	message := new(domain.Message)
	row := trx.QueryRow(
		ctx,
		`INSERT INTO "messages" ("owner_id", "dialog_id", "text")
		VALUES($1, $2, $3)
		RETURNING "id", "owner_id", "dialog_id", "text", "created_at";`,
		ownerId, dialogId, text,
	)
	err := message.ScanRow(row)
	if err != nil {
		return nil, err
	}
	return message, nil
}
