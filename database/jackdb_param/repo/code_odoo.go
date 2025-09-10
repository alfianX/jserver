package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (r Repo) CodeOdooGetName(ctx context.Context, host, typeTrx string) (string, string, error) {
	var codeOdoo model.CodeOdoo
	result := r.Db.WithContext(ctx).Select("name", "code").Where("host = ? AND type_trx = ?", host, typeTrx).Find(&codeOdoo)

	return codeOdoo.Name, codeOdoo.Code, result.Error
}
