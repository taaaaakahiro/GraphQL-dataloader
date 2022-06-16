package version

import (
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	res    *response
}

type response struct {
	Version string `json:"version"`
}

func NewHandler(logger *zap.Logger, ver string) *Handler {
	return &Handler{
		logger: logger,
		res: &response{
			Version: ver,
		},
	}
}
