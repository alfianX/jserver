package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (r Repo) AllowedHeadersGetHeaderName(ctx context.Context) ([]model.AllowedHeader, error) {
	var data []model.AllowedHeader

	result := r.Db.WithContext(ctx).Find(&data)

	return data, result.Error
}
