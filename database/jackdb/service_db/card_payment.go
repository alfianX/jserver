package servicedb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alfianX/jserver/config"
	"github.com/alfianX/jserver/database/jackdb/model"
	servicedb_param "github.com/alfianX/jserver/database/jackdb_param/service_db"
	"github.com/alfianX/jserver/helper"
	"github.com/alfianX/jserver/types"
)

type CardPaymentReqParam struct {
	SerialNumber    string
	TrxType         string
	Mti             string
	Procode         string
	Stan            string
	Trace           string
	Tid             string
	Mid             string
	Pan             string
	Amount          int64
	TransactionDate time.Time
	Nii             string
	De63            string
	Host            string
}

type CardPaymentResParam struct {
	ID                      int64
	ResponseCode            string
	ApprovalCode            string
	TransactionDateResponse time.Time
	Rrn                     string
}

func (s Service) CardPaymentSave(ctx context.Context, param *CardPaymentReqParam) (int64, error) {
	data := model.CardPayment{
		SerialNumber:    param.SerialNumber,
		TrxType:         param.TrxType,
		Mti:             param.Mti,
		Procode:         param.Procode,
		Stan:            param.Stan,
		Trace:           param.Trace,
		Tid:             param.Tid,
		Mid:             param.Mid,
		Pan:             param.Pan,
		Amount:          param.Amount,
		TransactionDate: param.TransactionDate,
		Nii:             param.Nii,
		De63:            param.De63,
		Host:            param.Host,
		CreatedAt:       time.Now(),
	}

	id, err := s.repo.CardPaymentSave(ctx, &data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) CardPaymentUpdateResponse(ctx context.Context, param *CardPaymentResParam) error {
	data := model.CardPayment{
		ID:                      param.ID,
		ResponseCode:            param.ResponseCode,
		ApprovalCode:            param.ApprovalCode,
		TransactionDateResponse: param.TransactionDateResponse,
		Rrn:                     param.Rrn,
		UpdatedAt:               time.Now(),
	}

	err := s.repo.CardPaymentUpdateResponse(ctx, &data)

	return err
}

func (s Service) CardPaymentUpdateFlagSuccess(ctx context.Context, id int64) error {
	err := s.repo.CardPaymentUpdateFlagSuccess(ctx, id)

	return err
}

func (s Service) CardPaymentSendToOdoo(ctx context.Context, cfg config.Config, jackdbParamService servicedb_param.Service) {
	// ok := true
	data, err := s.repo.CardPaymentGetDataForOdoo(ctx, s.repo.Db)
	if err != nil {
		fmt.Println("Cron-Card-Payment - Get data in db for odoo : " + err.Error())
		helper.ErrorLog("Cron-Card-Payment - Get data in db for odoo : " + err.Error())
		return
	}

	if len(data) == 0 {
		return
	}

	for _, row := range data {
		var payload types.CardPaymentOdoo

		payload.SerialNumber = row.SerialNumber
		payload.Mti = row.Mti
		payload.Procode = row.Procode
		payload.Stan = row.Stan
		payload.Trace = row.Trace
		payload.Tid = row.Tid
		payload.Mid = row.Mid
		payload.ExternalStoreId = row.Mid + "." + row.Tid
		payload.Pan = row.Pan
		payload.Amount = row.Amount / 100
		payload.PaymentDate = row.TransactionDate.UTC().Format(time.RFC3339)
		payload.Nii = row.Nii
		payload.De63 = row.De63
		payload.ResponseCode = row.ResponseCode
		payload.ApprovalCode = row.ApprovalCode
		payload.TransactionDateResponse = row.TransactionDateResponse.UTC().Format(time.RFC3339)
		payload.PaymentReferenceNo = row.Rrn

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			fmt.Println("Cron-Card-Payment - Marshal json odoo : " + err.Error())
			helper.ErrorLog("Cron-Card-Payment - Marshal json odoo : " + err.Error())
			break
		}

		jsonString := string(jsonBytes)

		fmt.Println("haloo: " + jsonString)

		// code := "031"
		// name := "Arthajasa PG Debits"

		name, code, err := jackdbParamService.CodeOdooGetName(ctx, row.Host, row.TrxType)
		if err != nil {
			fmt.Println("Cron-Card-Payment - Get code odoo : " + err.Error())
			helper.ErrorLog("Cron-Card-Payment - Get code odoo : " + err.Error())
			break
		}

		cookie, err := helper.AuthenticateOdoo(cfg.CnfGlob.OdooURL + "/web/session/authenticate")
		if err != nil {
			fmt.Println("Cron-Card-Payment - Get cookie odoo : " + err.Error())
			helper.ErrorLog("Cron-Card-Payment - Get cookie odoo : " + err.Error())
			break
		}

		if cookie == "" {
			fmt.Println("Cron-Card-Payment - Cookie odoo empty!")
			helper.ErrorLog("Cron-Card-Payment - Cookie odoo empty!")
			break
		}

		err = helper.SendToOdoo(cfg.CnfGlob.OdooURL+"/iid_api_manage/post_data", name, cookie, "Transaction", code, "Arthajasa", jsonString)
		if err != nil {
			fmt.Println("Cron-Card-Payment - Send to odoo : " + err.Error())
			helper.ErrorLog("Cron-Card-Payment - Send to odoo  : " + err.Error())
			break
		}

		err = s.repo.CardPaymentUpdateFlagOdoo(ctx, s.repo.Db, &model.CardPayment{
			ID:      row.ID,
			FlagOdo: 1,
		})
		if err != nil {
			fmt.Println("Cron-Card-Payment - Update flag odoo : " + err.Error())
			helper.ErrorLog("Cron-Card-Payment - Update flag odoo  : " + err.Error())
			break
		}
	}
}
