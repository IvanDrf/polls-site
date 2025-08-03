package cookies

import (
	"net/http"

	"github.com/IvanDrf/polls-site/internal/transport/auth/jwt"
)

type Cookier interface {
	SetAuthCookies(w http.ResponseWriter, access, refresh string)
}

type cookier struct {
}

func NewCookier() Cookier {
	return cookier{}
}

func (c cookier) SetAuthCookies(w http.ResponseWriter, access, refresh string) {
	http.SetCookie(w, &http.Cookie{
		Name:     jwter.AccessToken,
		Value:    access,
		Path:     "/",
		Secure:   true,
		HttpOnly: false,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     jwter.RefreshToken,
		Value:    refresh,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
}
