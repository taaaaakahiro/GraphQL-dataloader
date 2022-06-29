package server

import (
	"net"
	"net/http"

	"github.com/taaaaakahiro/GraphQL-dataloader-MySQL/pkg/handler"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

type Config struct {
	Log *zap.Logger
}

type Server struct {
	Mux     *http.ServeMux
	Handler http.Handler
	server  *http.Server
	handler *handler.Handler
	log     *zap.Logger
}

func NewServer(registry *handler.Handler, cfg *Config) *Server {
	s := &Server{
		Mux:     http.NewServeMux(),
		handler: registry,
	}
	if cfg != nil {
		if log := cfg.Log; log != nil {
			s.log = log
		}
	}
	s.registerHandler()
	return s
}

func (s *Server) registerHandler() {
	// graph ql
	s.Mux.Handle("/gql", playground.Handler("GraphQL playground", "/query"))
	s.Mux.Handle("/query", s.handler.V1.Query())
	// rest api
	s.Mux.Handle("/user/list", s.handler.V1.GetUsers())
	s.Mux.Handle("/message/list", s.handler.V1.GetMessages())
	// others
	s.Mux.Handle("/healthz", s.healthCheckHandler())
	s.Mux.Handle("/version", s.handler.Version.GerVersion())
}

func (s *Server) Serve(listener net.Listener) error {
	server := &http.Server{
		Handler: cors.Default().Handler(s.Mux),
	}
	s.server = server
	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil

}

func (s *Server) GracefulShutdown() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) healthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
