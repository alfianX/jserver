package cardpayment

import (
	"context"

	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database/jackdb/repo"
	servicedb "github.com/alfianX/jserver/database/jackdb/service_db"
	repo_param "github.com/alfianX/jserver/database/jackdb_param/repo"
	servicedb_param "github.com/alfianX/jserver/database/jackdb_param/service_db"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronService struct {
	cnf config.Config
	// cookieOdoo         string
	jackdbService      servicedb.Service
	jackdbParamService servicedb_param.Service
}

func NewCronJob(cnf config.Config, db, dbParam *gorm.DB) CronService {
	return CronService{
		cnf:                cnf,
		jackdbService:      servicedb.NewService(repo.NewRepo(db)),
		jackdbParamService: servicedb_param.NewService(repo_param.NewRepo(dbParam)),
	}
}

func (cs *CronService) CronJob(ctx context.Context) {
	c := cron.New()

	c.AddFunc("@every 10s", func() {
		go cs.jackdbService.CardPaymentSendToOdoo(ctx, cs.cnf, cs.jackdbParamService)
	})

	c.Start()

	<-ctx.Done()
	c.Stop()
}
