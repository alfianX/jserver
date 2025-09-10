package gateway

import (
	"github.com/alfianX/jserver/internal/gateway/handler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, lg *logrus.Logger, db *gorm.DB) {
	handler := handler.NewHandler(lg, r, db)

	r.Use(handler.MiddlewareLogger())
	r.GET("/healthz", handler.Health)
	s := r.Group("/api/jserver")
	s.Any("/*action", handler.Actions)
	// s.GET("/createKeyRSA", handler.CreateKeyRSA)
}
