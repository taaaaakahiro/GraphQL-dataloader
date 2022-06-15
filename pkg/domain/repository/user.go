package repository

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/domain/entity"
)

type IUserRepository interface {
	ListUsers() ([]entity.User, error)
	User(userId int) (entity.User, error)
	GetUsers(ctx context.Context, keys dataloader.Keys)
}
