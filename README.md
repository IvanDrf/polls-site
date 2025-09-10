# Polls-site

Backend application for creating polls and voting in them. During development, I decided to use *Monolithic* architecture, maybe in future i will rewrite this app with *Microservice* architecture, and even maybe i will write frontend, so it will be full-fledged web application. 

There are still some issues in the code that I plan to solve in the future, for example: *This is an opened transaction*:

Ofc service and app *works fine*, but one function depends on another, although each of them belongs to its own service, so there's still room for improvement.

```golang
func (v voteService) DeleteAllVotesInPoll(poll *models.Poll, r *http.Request) error
  //
  ...
  //

  v.transaction.RollBackTransaction()

  //
  ...
  //

  err = v.votesRepo.DeleteAllVotes(poll.QuestionId)
  	if err != nil {
  		v.transaction.RollBackTransaction()
  		return errs.ErrCantDeleteAllVotes()
  	}
  
  	// Don't commit transaction cuz deleting another tables
  
  	return nil
  }


func (p pollService) DeletePoll(poll *models.Poll, r *http.Request) error {
  //
  ...
  //

  if poll.UserId != question.UserId {
		return errs.ErrNotAdmin()
	}

	// Do not open transaction cuz it's already open in DeleteAllVotes

	err = p.answRepo.DeleteAllAnswers(poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteAnswer()
	}

	err = p.questRepo.DeleteQuestionById(poll.QuestionId)
	if err != nil {
		p.transaction.RollBackTransaction()
		return errs.ErrCantDeleteQuestion()
	}

	p.transaction.CommitTransaction()

	return nil
}

```

## Used Technologies
- GO
  - net/http
  - net/smtp
  - log/slog
  - database/sql
  - jwt
  - godotenv
  - crypto
    
- MySQL

## Architecture
- Transport Layer - HTTP handlers
- Service Layer - business rules/logic
- Repository Layer - database layer
  
<img width="400" height="400" alt="architecture" src="https://github.com/user-attachments/assets/2241f031-ee50-4e28-ad1e-1c297831ac61" />

  ***

## Structure
<details> <summary>Tree of Dirs/Files</summary>

```
├── cmd
│   └── main.go  
├── config
│   └── config.go 
├── go.mod
├── go.sum
├── internal
│   ├── database
│   │   └── database.go 
│   ├── errs
│   │   ├── answers.go
│   │   ├── auth.go
│   │   ├── email.go
│   │   ├── errs.go
│   │   ├── http.go
│   │   ├── jwt.go
│   │   ├── questions.go
│   │   ├── settings.go
│   │   └── votes.go
│   ├── models
│   │   └── models.go
│   ├── repo
│   │   ├── auth
│   │   │   ├── jwt
│   │   │   │   └── jwt.go
│   │   │   └── users
│   │   │       └── user.go
│   │   ├── polls
│   │   │   ├── answers
│   │   │   │   └── answers.go
│   │   │   ├── questions
│   │   │   │   └── questions.go
│   │   │   └── votes
│   │   │       └── votes.go
│   │   └── transaction
│   │       └── transaction.go
│   ├── service
│   │   ├── auth
│   │   │   ├── auth.go
│   │   │   ├── checker
│   │   │   │   ├── email.go
│   │   │   │   └── passw.go
│   │   │   ├── email
│   │   │   │   └── email.go
│   │   │   ├── hasher
│   │   │   │   └── hasher.go
│   │   │   └── links
│   │   │       └── links.go
│   │   ├── poll
│   │   │   └── poll.go
│   │   └── vote
│   │       └── vote.go
│   └── transport
│       ├── auth
│       │   ├── cookies
│       │   │   └── cookies.go
│       │   ├── jwt
│       │   │   └── jwt.go
│       │   └── middleware.go
│       ├── handlers
│       │   ├── auth.go
│       │   ├── handlers.go
│       │   ├── polls.go
│       │   └── votes.go
│       └── server
│           ├── routes.go
│           └── server.go
├── LICENSE
├── logger
│   └── logger.go
└── README.md
```
</details>

## Database 

Tables from the database for the application's functionality

<img width="500" height="500" alt="tables" src="https://github.com/user-attachments/assets/5a8c06b6-c898-48f1-af47-51b982f4de60" />

<details> <summary>Description of Tables</summary>

- ### Users

  > Table with users
  
    Columns:
    - id # user's id 
    - email # user's email
    - passw # user's password (hashed)
    - verificated # user's status true/false if he verified his email or not
    - veriftoken # user's personal link for verification, that sends in email

***

- ### Jwt

  > Table with user's refresh tokens
  
    Columns:
    - token_id #token's id
    - user_id # user id from table Users
    - token # refresh token

***

- ### Questions
  > Table with poll questions
  
    Columns:
    - id # question's id
    - question # question in poll, for example - "how old are u?"

***

- ### Answers

  > Table with user answers
  
    Columns:
    - id # answer id
    - answ # answer (text) option in poll, for example - "i'm 18 y.o."
    - question_id # question's id from table Questions

***

- ### Votes

  > Table with user votes in polls
  
    Columns:
    - question_id # question's id from table Questions
    - answ_id # answer's id from table Answers
    - user_id # user id from table Users

</details>
