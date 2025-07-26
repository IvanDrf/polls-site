package main

import (
	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/database"
	"github.com/IvanDrf/polls-site/internal/transport/server"
	"github.com/IvanDrf/polls-site/logger"
)

func main() {
	cfg := config.InitCFG()

	db := database.InitDB(cfg)
	defer db.Close()

	logger := logger.InitLogger(cfg)

	server := server.NewServer(cfg, db, logger)
	server.RegisterRoutes()
	server.Start(cfg)
}
