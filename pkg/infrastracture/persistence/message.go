package persistence

import (
	errs "github.com/pkg/errors"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/repository"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
)

type MessageRepo struct {
	database *io.SQLDatabase
}

var _ repository.IMessageRepository = (*MessageRepo)(nil)

func NewMessageReopsitory(db *io.SQLDatabase) *MessageRepo {
	return &MessageRepo{
		database: db,
	}
}

func (r *MessageRepo) ListMessages(userId int) ([]entity.Message, error) {
	messages := make([]entity.Message, 0)

	query := "SELECT id, user_id, message FROM message WHERE user_id = ? ORDER BY id DESC"
	stmtOut, err := r.database.Prepare(query)
	if err != nil {
		return nil, errs.WithStack(err)
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query(userId)
	if err != nil {
		return nil, errs.WithStack(err)
	}

	for rows.Next() {
		message := entity.Message{}

		err = rows.Scan(&message.Id, &message.UserId, &message.Message)
		if err != nil {
			return nil, errs.WithStack(err)
		}

		messages = append(messages, message)
	}
	if err != nil {
		return nil, errs.WithStack(err)
	}
	return messages, nil
}

func (r MessageRepo) Messages() ([]entity.Message, error) {
	messages := make([]entity.Message, 0)
	return messages, nil
}

func (r MessageRepo) CreateMessage(message *entity.Message) error {
	return nil
}
