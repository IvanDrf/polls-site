package models

type JWT struct {
	Id      int    `json:"-"`
	UserId  int    `json:"-"`
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}
