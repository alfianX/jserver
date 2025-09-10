package servicedb

import "github.com/alfianX/jserver/database/jackdb/repo"

type Service struct {
	repo repo.Repo
}

func NewService(r repo.Repo) Service {
	return Service{
		repo: r,
	}
}
