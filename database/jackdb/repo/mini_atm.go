package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r Repo) MiniAtmSave(ctx context.Context, data *model.MiniAtm) (int64, error) {
	result := r.Db.WithContext(ctx).Select(
		"serial_number",
		"merchant_group",
		"trx_type",
		"mti",
		"procode",
		"stan",
		"trace",
		"tid",
		"mid",
		"pan",
		"amount",
		"transaction_date",
		"nii",
		"de63",
		"host",
		"created_at",
	).Create(&data)

	return data.ID, result.Error
}

func (r Repo) MiniAtmUpdateResponse(ctx context.Context, data *model.MiniAtm) error {
	result := r.Db.WithContext(ctx).Model(&model.MiniAtm{ID: data.ID}).Updates(&model.MiniAtm{
		ResponseCode:            data.ResponseCode,
		TransactionDateResponse: data.TransactionDateResponse,
		Rrn:                     data.Rrn,
		UpdatedAt:               data.UpdatedAt,
	})

	return result.Error
}

func (r Repo) MiniAtmGetDataForOdoo(ctx context.Context, tx *gorm.DB) ([]model.MiniAtm, error) {
	var data []model.MiniAtm
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("trx_type IN ('TARIK TUNAI', 'TRANSFER', 'INFO SALDO') AND response_code IS NOT NULL AND flag_odoo = ?", 0).
		Order("created_at desc").
		Limit(4).
		Find(&data).Error

	return data, err
}

func (r Repo) MiniAtmUpdateFlagOdoo(ctx context.Context, tx *gorm.DB, data *model.MiniAtm) error {
	result := tx.WithContext(ctx).Model(&model.MiniAtm{ID: data.ID}).Updates(&model.MiniAtm{
		FlagOdoo: data.FlagOdoo,
	})

	return result.Error
}
