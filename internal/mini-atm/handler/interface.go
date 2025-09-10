package handler

import (
	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database/jackdb/repo"
	servicedb "github.com/alfianX/jserver/database/jackdb/service_db"
	repo_param "github.com/alfianX/jserver/database/jackdb_param/repo"
	servicedb_param "github.com/alfianX/jserver/database/jackdb_param/service_db"
	mini_atm "github.com/alfianX/jserver/proto/mini-atm"
	"gorm.io/gorm"
)

type Service struct {
	mini_atm.UnimplementedMiniAtmServiceServer
	config             config.Config
	jackdbService      servicedb.Service
	jackdbParamService servicedb_param.Service
}

func NewHandler(cnf config.Config, db, dbParam *gorm.DB) *Service {
	return &Service{
		config:             cnf,
		jackdbService:      servicedb.NewService(repo.NewRepo(db)),
		jackdbParamService: servicedb_param.NewService(repo_param.NewRepo(dbParam)),
	}
}
