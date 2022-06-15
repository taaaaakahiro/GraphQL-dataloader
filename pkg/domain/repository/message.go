package repository

import (
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
)

type IMessageRepository interface {
	ListMessages(userId int) ([]entity.Message, error)
	Messages() ([]entity.Message, error)
	CreateMessage(message *entity.Message) error
}
