package qrdynamic

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database"
	"github.com/alfianX/jserver/internal"
	"github.com/alfianX/jserver/internal/qr-dynamic/handler"
	qr_dynamic "github.com/alfianX/jserver/proto/qr-dynamic"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Service struct {
	logger     *logrus.Logger
	config     config.Config
	svcHandler *handler.Service
}

func NewService() (*Service, error) {
	cnf, err := config.NewParsedConfig()
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(cnf.CnfGlob.Database)
	if err != nil {
		return nil, err
	}

	dbParam, err := database.Connect(cnf.CnfGlob.DatabaseParam)
	if err != nil {
		return nil, err
	}

	log := internal.NewLogger()

	// cookie, err := h.AuthenticateOdoo(cnf.CnfGlob.OdooURL + "/web/session/authenticate")
	// if err != nil {
	// 	return nil, err
	// }

	// if cookie == "" {
	// 	return nil, errors.New("odoo fail login")
	// }

	svcHandler := handler.NewHandler(cnf, db, dbParam)

	s := Service{
		logger:     log,
		config:     cnf,
		svcHandler: svcHandler,
	}

	return &s, nil
}

func (s *Service) Run(ctx context.Context) error {

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", s.config.CnfLoc.ListenPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	qr_dynamic.RegisterQrDynamicServiceServer(grpcServer, s.svcHandler)

	reflection.Register(grpcServer)

	stopServer := make(chan os.Signal, 1)
	signal.Notify(stopServer, syscall.SIGINT, syscall.SIGTERM)

	defer signal.Stop(stopServer)

	serverErrors := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		s.logger.Printf("REST API listening on PORT %d", s.config.CnfLoc.ListenPort)
		serverErrors <- grpcServer.Serve(lis)
	}(&wg)

	select {
	case err := <-serverErrors:
		return fmt.Errorf("error: starting REST API server: %w", err)
	case <-stopServer:
		s.logger.Warn("server received STOP signal")

		grpcServer.GracefulStop()

		wg.Wait()
		s.logger.Info("Server was shutdown gracefully")
	}
	return nil
}
