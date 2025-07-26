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
		middleware: auth.NewMiddleware(cfg),
		handler:    handlers.NewHandler(cfg, db, logger),
	}
}

func (this *Server) Start(cfg *config.Config) {
	addr := fmt.Sprintf("%s:%s", cfg.ServerAddress, cfg.ServerPort)
	if err := http.ListenAndServe(addr, this.server); err != nil {
		log.Fatal(errs.ErrCantStartServer())
	}
}

func (this *Server) RegisterRoutes() {
	this.server.HandleFunc("POST /register", this.handler.RegisterUser)
	this.server.HandleFunc("POST /login", this.handler.LoginUser)

}
