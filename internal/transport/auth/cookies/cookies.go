package cookies

import "net/http"

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
		Name:     "access_jwt",
		Value:    access,
		Path:     "/",
		Secure:   true,
		HttpOnly: false,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_jwt",
		Value:    refresh,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
}
