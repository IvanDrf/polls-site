package service

import (
	"database/sql"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo"
)

type Service interface {
	RegisterUser(req *models.UserReq) error
}

type service struct {
	repo repo.Repo
}

func NewService(cfg *config.Config, db *sql.DB) Service {
	return service{repo: repo.NewRepo(cfg, db)}
}

func (this service) RegisterUser(req *models.UserReq) error {
	return nil
}
