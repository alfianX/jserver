package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb/model"
)

func (r Repo) QrGenerateSave(ctx context.Context, data *model.QrDynamicGenerate) (int64, error) {
	result := r.Db.WithContext(ctx).Select("host", "request_data", "created_at").Create(&data)

	return data.ID, result.Error
}

func (r Repo) QrGenerateUpdateResponse(ctx context.Context, data *model.QrDynamicGenerate) error {
	result := r.Db.WithContext(ctx).Model(&model.QrDynamicGenerate{ID: data.ID}).Updates(&model.QrDynamicGenerate{
		ResponseData: data.ResponseData,
		UpdatedAt:    data.UpdatedAt,
	})

	return result.Error
}

func (r Repo) QrGenerateUpdateFlag(ctx context.Context, data *model.QrDynamicGenerate) error {
	result := r.Db.WithContext(ctx).Model(&model.QrDynamicGenerate{ID: data.ID}).Updates(&model.QrDynamicGenerate{
		Flag:      data.Flag,
		UpdatedAt: data.UpdatedAt,
	})

	return result.Error
}
