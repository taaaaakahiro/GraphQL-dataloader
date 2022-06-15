package loader

import "github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"

type Loaders struct {
	UserLoader *UserLoader
}

func NewLoaders(r *persistence.Repositories) *Loaders {
	loaders := &Loaders{
		UserLoader: NewUserLoader(r),
	}
	return loaders
}
