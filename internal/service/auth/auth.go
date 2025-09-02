package auth

import (
	"database/sql"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
	"github.com/IvanDrf/polls-site/internal/repo/auth/users"
	"github.com/IvanDrf/polls-site/internal/repo/transaction"
	"github.com/IvanDrf/polls-site/internal/service/auth/email"
	"github.com/IvanDrf/polls-site/internal/service/auth/links"
	"github.com/golang-jwt/jwt"

	j "github.com/IvanDrf/polls-site/internal/repo/auth/jwt"
	"github.com/IvanDrf/polls-site/internal/service/auth/checker"
	"github.com/IvanDrf/polls-site/internal/service/auth/hasher"
	jwter "github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type Auther interface {
	RegisterUser(user *models.User) error
	VerifyEmail(link string) (string, string, error)
	LoginUser(user *models.User) (string, string, error)

	DeleteUnverifiedUsers()

	RefreshTokens(r *http.Request) (string, string, error)
}

type auth struct {
	pswChecker checker.PswChecker
	pswHasher  hasher.PswHasher
	emChecker  checker.EmailChecker

	jwter jwter.Jwter

	userRepo  users.UserRepo
	tokenRepo j.JWTRepo

	transaction transaction.Transactioner

	emailService email.EmailService
	linker       links.VerifLinker

	logger *slog.Logger
}

func NewAuthService(cfg *config.Config, db *sql.DB, logger *slog.Logger) Auther {
	return auth{
		pswChecker: checker.NewPSWChecker(),
		pswHasher:  hasher.NewPswHasher(),
		emChecker:  checker.NewEmailChecker(),

		jwter: jwter.NewJwter(cfg),

		userRepo:  users.NewRepo(cfg, db),
		tokenRepo: j.NewTokensRepo(cfg, db),

		transaction: transaction.NewTransactioner(cfg, db),

		emailService: email.NewEmailService(cfg),
		linker:       links.NewVerifLinker(cfg),

		logger: logger,
	}
}

func (a auth) RegisterUser(user *models.User) error {
	a.logger.Info("auth -> Register")

	if !a.emChecker.ValidEmail(user.Email) {
		return errs.ErrInvalidEmailInReg()
	}

	if !a.pswChecker.ValidPassword(user.Password) {
		return errs.ErrInvalidPswInReg()
	}

	if res, err := a.userRepo.FindUserByEmail(user.Email); res.Id != 0 || err == nil {
		return errs.ErrAlreadyInDB()
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		user.Password = a.pswHasher.HashPassword(user.Password)
	}()

	link, token := a.linker.CreateVerificationLink()
	if link == "" {
		return errs.ErrCantCreateVerifLink()
	}

	user.Verificated = false
	user.VerifToken = token
	user.Expired = time.Now().Add(24 * time.Hour)

	wg.Wait()
	a.transaction.StartTransaction()

	var err error
	user.Id, err = a.userRepo.AddUser(user)
	if err != nil {
		a.transaction.RollBackTransaction()

		a.logger.Error(err.Error())
		return errs.ErrCantRegister()
	}

	err = a.emailService.SendEmail(&models.EmailSending{Email: user.Email, Link: link}, email.VerifHeader, email.VerifBody)
	if err != nil {
		a.transaction.RollBackTransaction()

		a.logger.Error(err.Error())
		return err
	}

	a.transaction.CommitTransaction()

	return nil
}

func (a auth) LoginUser(user *models.User) (string, string, error) {
	a.logger.Info("auth -> Login")

	userInDB, err := a.userRepo.FindUserByEmail(user.Email)
	if err != nil {
		return "", "", errs.ErrCantFindUser()
	}

	if !userInDB.Verificated {
		return "", "", errs.ErrNotActivatedUser()
	}

	if !a.pswHasher.ComparePassword(userInDB.Password, user.Password) {
		return "", "", errs.ErrInvalidPswInLog()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&userInDB)
	if err != nil {
		return "", "", err
	}

	tokenInDB, err := a.tokenRepo.FindRefreshToken(userInDB.Id)
	if err != nil {
		err = a.tokenRepo.AddRefreshToken(userInDB.Id, refreshToken)
		if err != nil {
			return "", "", errs.ErrCantAddToken()
		}
	} else {
		err = a.tokenRepo.UpdateRefreshToken(tokenInDB.UserId, refreshToken)
	}

	return accessToken, refreshToken, err
}

func (a auth) VerifyEmail(link string) (string, string, error) {
	user, err := a.userRepo.FindUserByLink(link)
	if err != nil {
		return "", "", errs.ErrCantFindUserByLink()
	}

	if time.Now().After(user.Expired) {
		return "", "", errs.ErrExpiredLink()
	}

	err = a.userRepo.ActivateUser(&user)
	if err != nil {
		return "", "", errs.ErrCantActivateUser()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&user)
	if err != nil {
		a.transaction.RollBackTransaction()

		a.logger.Error(err.Error())
		return "", "", err
	}

	err = a.tokenRepo.AddRefreshToken(user.Id, refreshToken)
	if err != nil {
		a.transaction.RollBackTransaction()

		a.logger.Error(err.Error())
		return "", "", errs.ErrCantAddToken()
	}

	return accessToken, refreshToken, nil
}

func (a auth) DeleteUnverifiedUsers() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		if _, ok := <-ticker.C; ok {
			a.DeleteUnverifiedUsers()
		}
	}
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
