package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (r Repo) HostDebitGetAddress(ctx context.Context, name string) (string, error) {
	var hostDebit model.HostDebit
	result := r.Db.WithContext(ctx).Select("address").Where("name = ? AND status = ?", name, 1).Find(&hostDebit)

	return hostDebit.Address, result.Error
}
