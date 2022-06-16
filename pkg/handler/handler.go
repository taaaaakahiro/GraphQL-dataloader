package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	v1 "github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/handler/v1"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/handler/version"
	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/infrastracture/persistence"
	"go.uber.org/zap"
)

type Handler struct {
	V1      *v1.Handler
	Version *version.Handler
}

func NewHandler(logger *zap.Logger, repo *persistence.Repositories, query *handler.Server, ver string) *Handler {
	h := &Handler{
		V1:      v1.NewHandler(logger.Named("v1"), repo, query),
		Version: version.NewHandler(logger.Named("version"), ver),
	}
	return h
}
