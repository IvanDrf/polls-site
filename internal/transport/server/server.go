package server

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/transport/auth"
	"github.com/IvanDrf/polls-site/internal/transport/handlers"
)

type Server struct {
	server     *http.ServeMux
	middleware auth.Middleware
	handler    handlers.Handler
}

func NewServer(cfg *config.Config, db *sql.DB, logger *slog.Logger) *Server {
	return &Server{
		server:     http.NewServeMux(),
		middleware: auth.NewMiddleware(cfg, logger),
		handler:    handlers.NewHandler(cfg, db, logger),
	}
}

func (s *Server) Start(cfg *config.Config) {
	go s.handler.DeleteUnverifiedUsers()

	addr := fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.ServerPort)
	if err := http.ListenAndServe(addr, s.server); err != nil {
		log.Fatal(errs.ErrCantStartServer())
	}
}
