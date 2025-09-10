package repo

import "gorm.io/gorm"

type Repo struct {
	Db *gorm.DB
}

func NewRepo(db *gorm.DB) Repo {
	return Repo{Db: db}
}
