package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb/model"
)

func (r Repo) QrNotifSave(ctx context.Context, data *model.QrNotif) (int64, error) {
	result := r.Db.WithContext(ctx).Select("request_data", "created_at").Create(&data)

	return data.ID, result.Error
}

func (r Repo) QrNotifUpdateResponse(ctx context.Context, data *model.QrNotif) error {
	result := r.Db.WithContext(ctx).Model(&model.QrNotif{ID: data.ID}).Updates(&model.QrNotif{
		ResponseData: data.ResponseData,
		UpdatedAt:    data.UpdatedAt,
	})

	return result.Error
}

func (r Repo) QrNotifUpdateFlag(ctx context.Context, data *model.QrNotif) error {
	result := r.Db.WithContext(ctx).Model(&model.QrNotif{ID: data.ID}).Updates(&model.QrNotif{
		Flag:      data.Flag,
		UpdatedAt: data.UpdatedAt,
	})

	return result.Error
}
