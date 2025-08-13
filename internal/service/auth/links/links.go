package links

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/IvanDrf/polls-site/config"
)

type VerifLinker interface {
	CreateVerificationLink() string
}

type verifLinker struct {
	host string
	port string
}

func NewVerifLinker(cfg *config.Config) VerifLinker {
	return verifLinker{
		host: cfg.ServerAddress,
		port: cfg.ServerPort,
	}
}

func (v verifLinker) CreateVerificationLink() string {
	buff := make([]byte, 32)
	if _, err := rand.Read(buff); err != nil {
		return ""
	}

	return fmt.Sprintf("http://%s:%s/verify-email?token=%s", v.host, v.port, hex.EncodeToString(buff))
}
