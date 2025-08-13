package email

import (
	"fmt"
	"net/smtp"

	"github.com/IvanDrf/polls-site/config"
	"github.com/IvanDrf/polls-site/internal/errs"
	"github.com/IvanDrf/polls-site/internal/models"
)

type EmailService interface {
	SendEmail(verif *models.EmailSending, header string, body string) error
}

type emailService struct {
	email    string
	password string

	smtpHost string
	smtpPort string
}

func NewEmailService(cfg *config.Config) EmailService {
	return emailService{
		email:    cfg.Email,
		password: cfg.Password,

		smtpHost: cfg.SmtpHost,
		smtpPort: cfg.SmtpPort,
	}
}

func (e emailService) SendEmail(verif *models.EmailSending, header string, body string) error {
	auth := smtp.PlainAuth("", e.email, e.password, e.smtpHost)

	message := e.CreateEmail(verif, header, body)
	if message == "" {
		return errs.ErrInvalidEmail()
	}

	return smtp.SendMail(
		fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort),
		auth,
		e.email,
		[]string{verif.Email},
		[]byte(message),
	)
}

const (
	VerifHeader = "Verify your email"
	VerifBody   = "Clicl the following link to vierify your email: "
)

func (e emailService) CreateEmail(email *models.EmailSending, header string, body string) string {
	switch header {
	case VerifHeader:
		header = VerifHeader

	default:
		return ""
	}

	switch body {
	case VerifBody:
		body = fmt.Sprintf(VerifBody+"%s", email.Email)

	default:
		return ""
	}

	msg := fmt.Sprintf("To: %s\r\n"+"Subject: %s\r\n"+"\r\n"+"%s\r\n", email.Email, header, body)

	return msg
}
