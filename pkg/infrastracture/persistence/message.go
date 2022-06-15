package persistence

import (
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
