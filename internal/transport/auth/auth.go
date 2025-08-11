package auth

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo/auth/users"
	"github.com/golang-jwt/jwt"

	"github.com/IvanDrf/polls-site/internal/repo/auth/tokens"
	"github.com/IvanDrf/polls-site/internal/transport/auth/checker"
	"github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type Auther interface {
	RegisterUser(req *models.UserReq) (string, string, error)
	LoginUser(req *models.UserReq) (string, string, error)

	RefreshTokens(r *http.Request) (string, string, error)
}

type auth struct {
	pswChecker checker.PswChecker
	pswHasher  checker.PswHasher
	emChecker  checker.EmailChecker

	jwter jwter.Jwter

	userRepo  users.UserRepo
	tokenRepo tokens.TokensRepo

	logger *slog.Logger
}

func NewAuthService(cfg *config.Config, db *sql.DB, logger *slog.Logger) Auther {
	return auth{
		pswChecker: checker.NewPSWChecker(),
		pswHasher:  checker.NewPswHasher(),
		emChecker:  checker.NewEmailChecker(),

		jwter: jwter.NewJwter(cfg),

		userRepo:  users.NewRepo(cfg, db),
		tokenRepo: tokens.NewTokensRepo(cfg, db),

		logger: logger,
	}
}

func (a auth) RegisterUser(req *models.UserReq) (string, string, error) {
	a.logger.Info("auth -> Register")

	if !a.emChecker.ValidEmail(req.Email) {
		return "", "", errs.ErrInvalidEmailInReg()
	}

	if !a.pswChecker.ValidPassword(req.Password) {
		return "", "", errs.ErrInvalidPswInReg()
	}

	if res, err := a.userRepo.FindUserByEmail(req.Email); res.Id != 0 || err == nil {
		return "", "", errs.ErrAlreadyInDB()
	}

	req.Password = a.pswHasher.HashPassword(req.Password)

	if err := a.userRepo.AddUser(req); err != nil {
		return "", "", errs.ErrCantRegister()
	}

	user, err := a.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", "", errs.ErrCantRegister()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&user)
	if err != nil {
		return "", "", err
	}

	err = a.tokenRepo.AddRefreshToken(user.Id, refreshToken)
	if err != nil {
		return "", "", errs.ErrCantAddToken()
	}

	return accessToken, refreshToken, nil
}

func (a auth) LoginUser(req *models.UserReq) (string, string, error) {
	a.logger.Info("auth -> Login")

	user, err := a.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", "", errs.ErrCantFindUser()
	}

	if !a.pswHasher.ComparePassword(user.Password, req.Password) {
		return "", "", errs.ErrInvalidPswInLog()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&user)
	if err != nil {
		return "", "", err
	}

	tokenInDB, err := a.tokenRepo.FindRefreshToken(user.Id)
	if err != nil {
		err = a.tokenRepo.AddRefreshToken(user.Id, refreshToken)
		if err != nil {
			return "", "", errs.ErrCantAddToken()
		}
	} else {
		err = a.tokenRepo.UpdateRefreshToken(tokenInDB.UserId, refreshToken)
	}

	return accessToken, refreshToken, err
}

func (a auth) RefreshTokens(r *http.Request) (string, string, error) {
	a.logger.Info("auth -> Refresh")
	refreshToken, err := a.jwter.GetToken(r, jwter.RefreshToken)
	if err != nil {
		return "", "", err
	}

	token, err := a.jwter.ParseToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errs.ErrInValidToken()
	}

	userId := int(claims[jwter.UserId].(float64))

	user, err := a.userRepo.FindUserById(userId)
	if err != nil {
		return "", "", errs.ErrInValidToken()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&user)
	if err != nil {
		return "", "", errs.ErrCantCreateToken()
	}

	err = a.tokenRepo.UpdateRefreshToken(userId, refreshToken)
	if err != nil {
		return "", "", errs.ErrCantAddToken()
	}

	return accessToken, refreshToken, nil
}
