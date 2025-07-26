package auth

import (
	"database/sql"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo"
	"github.com/IvanDrf/polls-site/internal/transport/auth/checker"
)

type Auther interface {
	RegisterUser(req *models.RegisterReq) error
}

type auth struct {
	pswChecker checker.PswChecker
	pswHasher  checker.PswHasher
	emChecker  checker.EmailChecker

	repo repo.Repo
}

func NewAuthService(cfg *config.Config, db *sql.DB) Auther {
	return auth{
		pswChecker: checker.NewPSWChecker(),
		pswHasher:  checker.NewPswHasher(),
		emChecker:  checker.NewEmailChecker(),

		repo: repo.NewRepo(cfg, db),
	}
}

func (this auth) RegisterUser(req *models.RegisterReq) error {
	if !this.emChecker.ValidEmail(req.Email) {
		return errs.ErrInvalidEmailInReg()
	}

	if !this.pswChecker.ValidPassword(req.Password) {
		return errs.ErrInvalidPswInReg()
	}

	if res, err := this.repo.FindUserByEmail(req.Email); res.Id != 0 || err == nil {
		return errs.ErrAlreadyInDB()
	}

	req.Password = this.pswHasher.HashPassword(req.Password)

	if err := this.repo.RegisterUser(req); err != nil {
		return errs.ErrCantRegister()
	}

	return nil
}
