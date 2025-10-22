package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database"
	"github.com/alfianX/jserver/internal"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger *logrus.Logger
	router *gin.Engine
	config config.Config
}

func NewService() (*Service, error) {
	cnf, err := config.NewParsedConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(cnf.CnfGlob.DatabaseParam)
	if err != nil {
		return nil, err
	}

	log := internal.NewLogger()

	gin.SetMode(cnf.CnfGlob.Mode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	RegisterRoutes(router, log, db)

	s := Service{
		logger: log,
		config: cnf,
		router: router,
	}

	return &s, nil
}

func (s *Service) Run(ctx context.Context) error {
	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.config.CnfLoc.ListenPort),
		Handler: s.router,
	}

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("REST API listening on PORT %d", s.config.CnfLoc.ListenPort)
		serverErrors <- server.ListenAndServe()
	}(&wg)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received STOP signal")

		err := server.Shutdown(ctx)
		if err != nil {
			return fmt.Errorf("graceful shutdown did not complete: %w", err)
		}
		wg.Wait()
		s.logger.Info("Server was shutdown gracefully")
	}
	return nil
}
