package models

// Register/Login User
type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWT struct {
	UserId  int    `json:"-"`
	Id      int    `json:"-"`
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}
