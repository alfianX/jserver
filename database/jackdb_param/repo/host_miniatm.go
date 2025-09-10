package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (r Repo) HostMiniatmGetAddress(ctx context.Context, name string) (string, error) {
	var hostMiniatm model.HostMiniatm
	result := r.Db.WithContext(ctx).Select("address").Where("name = ? AND status = ?", name, 1).Find(&hostMiniatm)

	return hostMiniatm.Address, result.Error
}
