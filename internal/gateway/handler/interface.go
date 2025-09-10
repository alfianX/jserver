package handler

import (
	"sync"

	"github.com/alfianX/jserver/database/jackdb_param/repo"
	servicedb "github.com/alfianX/jserver/database/jackdb_param/service_db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type service struct {
	mu                 sync.Mutex
	logger             *logrus.Logger
	router             *gin.Engine
	jackdbParamService servicedb.Service
}

func NewHandler(lg *logrus.Logger, rtr *gin.Engine, db *gorm.DB) service {
	return service{
		logger:             lg,
		router:             rtr,
		jackdbParamService: servicedb.NewService(repo.NewRepo(db)),
	}
}
