package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb_param/model"
)

func (r Repo) GetAddressByName(ctx context.Context, name string) (string, error) {
	var hostQr model.HostQr
	result := r.Db.WithContext(ctx).Select("address").Where("name = ? AND status = ?", name, "1").Find(&hostQr)

	return hostQr.Address, result.Error
}

func (r Repo) GetAddressQueryByName(ctx context.Context, name string) (string, error) {
	var hostQr model.HostQr
	result := r.Db.WithContext(ctx).Select("address_query").Where("name = ? AND status = ?", name, "1").Find(&hostQr)

	return hostQr.AddressQuery, result.Error
}
