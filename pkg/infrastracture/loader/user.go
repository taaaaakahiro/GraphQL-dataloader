package loader

import (
	"context"

	"github.com/graph-gophers/dataloader"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/graph/model"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
)

type UserLoader struct {
	loader *dataloader.Loader
}

func NewUserLoader(r *persistence.Repositories) *UserLoader {
	return &UserLoader{
		loader: dataloader.NewBatchedLoader(r.User.GetUsers),
	}
}

func (l *Loaders) GetUser(ctx context.Context, userID string) (*model.User, error) {
	// HACK:
	// userIDをキーにして、オンメモリでキャッシュ。既に該当キーがキャッシュにあればキャッシュから、なければDBから取得している。
	thunk := l.UserLoader.loader.Load(ctx, dataloader.StringKey(userID))
	result, err := thunk()
	if err != nil {
		return nil, err
	}
	return result.(*model.User), nil
}
