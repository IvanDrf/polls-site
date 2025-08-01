package auth

import (
	"database/sql"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	repo "github.com/IvanDrf/polls-site/internal/repo/auth/user"
	"github.com/IvanDrf/polls-site/internal/transport/auth/checker"
	"github.com/IvanDrf/polls-site/internal/transport/auth/cookies"
	"github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type Auther interface {
	RegisterUser(req *models.UserReq) error
	LoginUser(req *models.UserReq) (string, string, error)
}

type auth struct {
	pswChecker checker.PswChecker
	pswHasher  checker.PswHasher
	emChecker  checker.EmailChecker

	jwter   jwt.Jwter
	cookier cookies.Cookier

	repo repo.UserRepo
}

func NewAuthService(cfg *config.Config, db *sql.DB) Auther {
	return auth{
		pswChecker: checker.NewPSWChecker(),
		pswHasher:  checker.NewPswHasher(),
		emChecker:  checker.NewEmailChecker(),

		jwter:   jwt.NewJwter(cfg),
		cookier: cookies.NewCookier(),

		repo: repo.NewRepo(cfg, db),
	}
}

func (a auth) RegisterUser(req *models.UserReq) error {
	if !a.emChecker.ValidEmail(req.Email) {
		return errs.ErrInvalidEmailInReg()
	}

	if !a.pswChecker.ValidPassword(req.Password) {
		return errs.ErrInvalidPswInReg()
	}

	if res, err := a.repo.FindUserByEmail(req.Email); res.Id != 0 || err == nil {
		return errs.ErrAlreadyInDB()
	}

	req.Password = a.pswHasher.HashPassword(req.Password)

	if err := a.repo.AddUser(req); err != nil {
		return errs.ErrCantRegister()
	}

	return nil
}

func (a auth) LoginUser(req *models.UserReq) (string, string, error) {
	user, err := a.repo.FindUserByEmail(req.Email)
	if err != nil {
		return "", "", errs.ErrCantFindUser()
	}

	if !a.pswHasher.ComparePassword(user.Password, req.Password) {
		return "", "", errs.ErrInvalidPswInLog()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&user)
	//TODO add tokens in database 'tokens'

	return accessToken, refreshToken, err
}

//TODO write refresh, so user doesn't need to login always
