package auth

import (
	"database/sql"
	"log/slog"
	"net/http"
	"sync"

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
	RegisterUser(user *models.UserReq) (string, string, error)
	LoginUser(user *models.UserReq) (string, string, error)

	RefreshTokens(r *http.Request) (string, string, error)

	//ResetPassword(req *models.User) (models.JWT, error)
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

func (a auth) RegisterUser(user *models.UserReq) (string, string, error) {
	a.logger.Info("auth -> Register")

	if !a.emChecker.ValidEmail(user.Email) {
		return "", "", errs.ErrInvalidEmailInReg()
	}

	if !a.pswChecker.ValidPassword(user.Password) {
		return "", "", errs.ErrInvalidPswInReg()
	}

	if res, err := a.userRepo.FindUserByEmail(user.Email); res.Id != 0 || err == nil {
		return "", "", errs.ErrAlreadyInDB()
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		user.Password = a.pswHasher.HashPassword(user.Password)
	}()

	link := a.linker.CreateVerificationLink()
	if link == "" {
		return "", "", errs.ErrCantCreateVerifLink()
	}

	a.emailService.SendEmail(&models.EmailSending{Email: user.Email, Link: link}, email.VerifHeader, email.VerifHeader)

	wg.Wait()
	a.transaction.StartTransaction()
	if err := a.userRepo.AddUser(user); err != nil {
		a.transaction.RollBackTransaction()
		return "", "", errs.ErrCantRegister()
	}

	userInDB, err := a.userRepo.FindUserByEmail(user.Email)
	if err != nil {
		a.transaction.RollBackTransaction()
		return "", "", errs.ErrCantRegister()
	}

	accessToken, refreshToken, err := a.jwter.GenerateTokens(&userInDB)
	if err != nil {
		a.transaction.RollBackTransaction()
		return "", "", err
	}

	err = a.tokenRepo.AddRefreshToken(userInDB.Id, refreshToken)
	if err != nil {
		a.transaction.RollBackTransaction()
		return "", "", errs.ErrCantAddToken()
	}

	a.transaction.CommitTransaction()

	return accessToken, refreshToken, nil
}

func (a auth) LoginUser(user *models.UserReq) (string, string, error) {
	a.logger.Info("auth -> Login")

	userInDB, err := a.userRepo.FindUserByEmail(user.Email)
	if err != nil {
		return "", "", errs.ErrCantFindUser()
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

//func (a auth) ResetPassword(user *models.User) (models.JWT, error) {
//	userInDB, err := a.userRepo.FindUserByEmail(user.Email)
//	if err != nil {
//		return models.JWT{}, errs.ErrCantFindUser()
//	}
//
//	user.Id = userInDB.Id
//
//	err = a.userRepo.ResetPassword(a.pswHasher.HashPassword(user.Password), user.Id)
//	if err != nil {
//		return models.JWT{}, errs.ErrCantResetPassword()
//	}
//
//	//TODO: think about email verify and reset password by email
//}
