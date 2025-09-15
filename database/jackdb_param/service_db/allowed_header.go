package servicedb

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (s Service) AllowedHeadersGetHeaderName(ctx context.Context) ([]model.AllowedHeader, error) {
	data, err := s.repo.AllowedHeadersGetHeaderName(ctx)

	return data, err
}
