package models

// Register/Login User
type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// User in db
type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT in db
type JWT struct {
	Id      int    `json:"-"`
	UserId  int    `json:"-"`
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type Question struct {
	Id       int    `json:"id"`
	UserId   int    `json:"-"`
	Question string `json:"question"`
}

type Answer struct {
	Id         int    `json:"id"`
	QuestionId int    `json:"q_id"`
	Answer     string `json:"answer"`
}

type Vote struct {
	Id         int `json:"id"`
	UserId     int `json:"user_id"`
	QuestionId int `json:"question_id"`
	AnswerId   int `json:"answer_id"`
}

type PollRes struct {
	QuestionId int         `json:"question_id"`
	Answers    map[int]int `json:"answers"`
}

type Poll struct {
	QuestionId int      `json:"question_id"`
	Question   string   `json:"question"`
	UserId     int      `json:"-"`
	Answers    []string `json:"answers"`
}

type PollId struct {
	Id        int   `json:"question_id"`
	AnswersId []int `json:"answers_id"`
}

type EmailSending struct {
	Email string
	Link  string
}
