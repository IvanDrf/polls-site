package models

// Register User
type RegisterReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       int
	Email    string `json:"email"`
	Password string `json:"password"`
}
