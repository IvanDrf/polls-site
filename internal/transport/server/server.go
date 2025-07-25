package server

import "net/http"

type Server struct {
	server *http.ServeMux
}
