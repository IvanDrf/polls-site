package server

func (s *Server) RegisterRoutes() {
	s.server.HandleFunc("POST /register", s.handler.RegisterUser) // auth
	s.server.HandleFunc("POST /login", s.handler.LoginUser)       //auth
	s.server.HandleFunc("POST /refresh", s.handler.RefreshTokens) // auth

	s.server.HandleFunc("POST /poll/create", s.middleware.AuthMiddleware(s.handler.CreatePoll)) // poll
	s.server.HandleFunc("POST /poll/delete", s.middleware.AuthMiddleware(s.handler.DeletePoll)) // poll
	s.server.Handle("POST /poll/vote", s.middleware.AuthMiddleware(s.handler.VoteInPoll))       //poll
}
