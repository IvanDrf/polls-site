package auth

import (
	"database/sql"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo"
)

type Auther interface {
	RegisterUser(req *models.RegisterReq) error
}

type auth struct {
	pswChecker PSWChecker
	repo       repo.Repo
}

func NewAuthService(cfg *config.Config, db *sql.DB) Auther {
	return auth{
		repo:       repo.NewRepo(cfg, db),
		pswChecker: NewPSWChecker(),
	}
}

func (this auth) RegisterUser(req *models.RegisterReq) error {
	return nil
}
