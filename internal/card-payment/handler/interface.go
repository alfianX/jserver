package handler

import (
	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database/jackdb/repo"
	servicedb "github.com/alfianX/jserver/database/jackdb/service_db"
	repo_param "github.com/alfianX/jserver/database/jackdb_param/repo"
	servicedb_param "github.com/alfianX/jserver/database/jackdb_param/service_db"
	card_payment "github.com/alfianX/jserver/proto/card-payment"
	"gorm.io/gorm"
)

type Service struct {
	card_payment.UnimplementedCardPaymentServiceServer
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
