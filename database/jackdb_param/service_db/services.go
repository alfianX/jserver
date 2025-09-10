package servicedb

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (s Service) GetServices(ctx context.Context, prefix string) (model.Services, error) {
	srv, err := s.repo.GetServices(ctx, prefix)

	return srv, err
}
