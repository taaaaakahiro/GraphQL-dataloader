package loader

import (
	"github.com/graph-gophers/dataloader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
)

type UserLoader struct {
	loader *dataloader.Loader
}

func NewUserLoader(r persistence.Repositories) *UserLoader {
	return &UserLoader{
		loader: dataloader.NewBatchedLoader(r.User.GetUsers),
	}
}
