package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb/model"
)

func (r Repo) QrCallbackSave(ctx context.Context, data *model.QrDynamicCallback) (int64, error) {
	result := r.Db.WithContext(ctx).Select("host", "request_data", "created_at").Create(&data)

	return data.ID, result.Error
}

func (r Repo) QrCallbackUpdateResponse(ctx context.Context, data *model.QrDynamicCallback) error {
	result := r.Db.WithContext(ctx).Model(&model.QrDynamicCallback{ID: data.ID}).Updates(&model.QrDynamicCallback{
		ResponseData: data.ResponseData,
		UpdatedAt:    data.UpdatedAt,
	})

	return result.Error
}

func (r Repo) QrCallbcakUpdateFlag(ctx context.Context, data *model.QrDynamicCallback) error {
	result := r.Db.WithContext(ctx).Model(&model.QrDynamicCallback{ID: data.ID}).Updates(&model.QrDynamicCallback{
		Flag:      data.Flag,
		UpdatedAt: data.UpdatedAt,
	})

	return result.Error
}
