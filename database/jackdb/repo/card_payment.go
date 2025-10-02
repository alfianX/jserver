package repo

import (
	"context"

	"github.com/alfianX/jserver/database/jackdb/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r Repo) CardPaymentSave(ctx context.Context, data *model.CardPayment) (int64, error) {
	result := r.Db.WithContext(ctx).Select(
		"serial_number",
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

func (r Repo) CardPaymentUpdateResponse(ctx context.Context, data *model.CardPayment) error {
	result := r.Db.WithContext(ctx).Model(&model.CardPayment{ID: data.ID}).Updates(&model.CardPayment{
		ResponseCode:            data.ResponseCode,
		ApprovalCode:            data.ApprovalCode,
		TransactionDateResponse: data.TransactionDateResponse,
		Rrn:                     data.Rrn,
		UpdatedAt:               data.UpdatedAt,
	})

	return result.Error
}

func (r Repo) CardPaymentGetDataForOdoo(ctx context.Context, tx *gorm.DB) ([]model.CardPayment, error) {
	var data []model.CardPayment
	err := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("trx_type = ? AND mti = ? AND response_code = ? AND flag_odo = ?", "SALE", "0200", "00", 0).
		Order("created_at desc").
		Limit(4).
		Find(&data).Error

	return data, err
}

func (r Repo) CardPaymentUpdateFlagOdoo(ctx context.Context, tx *gorm.DB, data *model.CardPayment) error {
	result := tx.WithContext(ctx).Model(&model.CardPayment{ID: data.ID}).Updates(&model.CardPayment{
		FlagOdo: data.FlagOdo,
	})

	return result.Error
}
