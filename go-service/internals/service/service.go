package service

import (
	"go-servie/internals/repo"
)

type Service struct {
	repo *repo.Repo
}

func NewService(repo *repo.Repo) *Service {
	return &Service{
		repo: repo,
	}
}
