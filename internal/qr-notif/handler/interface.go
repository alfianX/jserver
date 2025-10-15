package handler

import (
	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database/jackdb/repo"
	servicedb "github.com/alfianX/jserver/database/jackdb/service_db"
	qr_notif "github.com/alfianX/jserver/proto/qr-notif"
	"gorm.io/gorm"
)

type Service struct {
	qr_notif.UnimplementedQrNotifServiceServer
	config config.Config
	// cookieOdoo    string
	jackdbService servicedb.Service
}

// func NewHandler(cookieOdoo string, cnf config.Config, db *gorm.DB) *Service {
func NewHandler(cnf config.Config, db *gorm.DB) *Service {
	return &Service{
		config: cnf,
		// cookieOdoo:    cookieOdoo,
		jackdbService: servicedb.NewService(repo.NewRepo(db)),
	}
}
