# Polls-site

Backend application for creating polls and voting in them
Used Technologies
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

## Database 

Tables from the database for the application's functionality

<img width="500" height="500" alt="tables" src="https://github.com/user-attachments/assets/5a8c06b6-c898-48f1-af47-51b982f4de60" />


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

- ### Questions
  > Table with poll questions
  
    Columns:
    - id # question's id
    - question # question in poll, for example - "how old are u?"


- ### Answers

  > Table with user answers
  
    Columns:
    - id # answer id
    - answ # answer (text) option in poll, for example - "i'm 18 y.o."
    - question_id # question's id from table Questions


- ### Votes

  > Table with user votes in polls
  
    Columns:
    - question_id # question's id from table Questions
    - answ_id # answer's id from table Answers
    - user_id # user id from table Users
