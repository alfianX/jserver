package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (r Repo) GetServices(ctx context.Context, prefix string) (model.Services, error) {
	var services model.Services
	result := r.Db.WithContext(ctx).Where("http_prefix = ?", prefix).Find(&services)

	return services, result.Error
}
