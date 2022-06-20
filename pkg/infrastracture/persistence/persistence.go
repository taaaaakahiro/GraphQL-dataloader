package persistence

import (
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/repository"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/io"
)

type Repositories struct {
	User    repository.IUserRepository
	Message repository.IMessageRepository
}

func NewReopsitories(db *io.SQLDatabase) (*Repositories, error) {
	return &Repositories{
		User:    NewUserRepository(db),
		Message: NewMessageReopsitory(db),
	}, nil
}
