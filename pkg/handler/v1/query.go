package v1

import "github.com/99designs/gqlgen/graphql/handler"

func (h *Handler) Query() *handler.Server {
	return h.query
}
