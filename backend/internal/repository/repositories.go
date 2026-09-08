package repository

import "github.com/VikyCham/go-boilerplate/internal/server"

type Repositories struct{}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{}
}