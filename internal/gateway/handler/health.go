package handler

import (
	"net/http"

	h "github.com/alfianX/jserver/helper"
	"github.com/gin-gonic/gin"
)

func (s *service) Health(c *gin.Context) {
	h.Respond(c, gin.H{"Message": "App OK"}, http.StatusOK)
}
