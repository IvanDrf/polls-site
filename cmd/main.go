package main

import (
	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/database"
)

func main() {
	cfg := config.InitCFG()
	db := database.InitDB(cfg)
	defer db.Close()

}
