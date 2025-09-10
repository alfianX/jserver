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

type MiniAtmReqParam struct {
	SerialNumber    string
	MerchantGroup   string
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

type MiniAtmResParam struct {
	ID                      int64
	ResponseCode            string
	TransactionDateResponse time.Time
	Rrn                     string
}

func (s Service) MiniAtmSave(ctx context.Context, param *MiniAtmReqParam) (int64, error) {
	data := model.MiniAtm{
		SerialNumber:    param.SerialNumber,
		MerchantGroup:   param.MerchantGroup,
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

	id, err := s.repo.MiniAtmSave(ctx, &data)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s Service) MiniAtmUpdateResponse(ctx context.Context, param *MiniAtmResParam) error {
	data := model.MiniAtm{
		ID:                      param.ID,
		ResponseCode:            param.ResponseCode,
		TransactionDateResponse: param.TransactionDateResponse,
		Rrn:                     param.Rrn,
		UpdatedAt:               time.Now(),
	}

	err := s.repo.MiniAtmUpdateResponse(ctx, &data)

	return err
}

func (s Service) SendToOdoo(ctx context.Context, cfg config.Config, jackdbParamService servicedb_param.Service) {
	ok := true
	tx := s.repo.Db.Begin()
	data, err := s.repo.MiniAtmGetDataForOdoo(ctx, tx)
	if err != nil {
		tx.Rollback()
		fmt.Println("Cron-Mini-ATM - Get data in db for odoo : " + err.Error())
		helper.ErrorLog("Cron-Mini-ATM - Get data in db for odoo : " + err.Error())
	}

	if len(data) == 0 {
		tx.Rollback()
		return
	}

	for _, row := range data {
		var payload types.MiniAtmOdoo

		payload.SerialNumber = row.SerialNumber
		payload.TrxType = row.TrxType
		payload.Mti = row.Mti
		payload.Procode = row.Procode
		payload.Stan = row.Stan
		payload.Trace = row.Trace
		payload.Tid = row.Tid
		payload.Mid = row.Mid
		payload.ExternalStoreId = row.Mid + "." + row.Tid
		payload.Pan = row.Pan
		payload.Amount = row.Amount
		payload.PaymentDate = row.TransactionDate.UTC().Format(time.RFC3339)
		payload.Nii = row.Nii
		payload.De63 = row.De63
		payload.ResponseCode = row.ResponseCode
		payload.TransactionDateResponse = row.TransactionDateResponse.UTC().Format(time.RFC3339)
		payload.PaymentReferenceNo = row.Rrn

		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			tx.Rollback()
			fmt.Println("Cron-Mini-ATM - Marshal json odoo : " + err.Error())
			helper.ErrorLog("Cron-Mini-ATM - Marshal json odoo : " + err.Error())
			break
		}

		jsonString := string(jsonBytes)

		// var code string
		// var name string
		// switch row.TrxType {
		// case "INFO SALDO":
		// 	code = "034"
		// 	name = "Arthajasa PG Info Saldo"
		// case "TRANSFER":
		// 	code = "032"
		// 	name = "Arthajasa PG Transfer"
		// case "TARIK TUNAI":
		// 	code = "033"
		// 	name = "Arthajasa PG Tarik Tunai"
		// }
		var trxType string
		if row.MerchantGroup == "AGENT" {
			trxType = row.TrxType
		} else if row.MerchantGroup == "MERCHANT" {
			trxType = "SALE"
		}

		name, code, err := jackdbParamService.CodeOdooGetName(ctx, row.Host, trxType)
		if err != nil {
			tx.Rollback()
			fmt.Println("Cron-Card-Payment - Get code odoo : " + err.Error())
			helper.ErrorLog("Cron-Card-Payment - Get code odoo : " + err.Error())
			ok = false
		}

		cookie, err := helper.AuthenticateOdoo(cfg.CnfGlob.OdooURL + "/web/session/authenticate")
		if err != nil {
			tx.Rollback()
			fmt.Println("Cron-Mini-ATM - Get cookie odoo : " + err.Error())
			helper.ErrorLog("Cron-Mini-ATM - Get cookie odoo : " + err.Error())
			ok = false
		}

		if cookie == "" {
			tx.Rollback()
			fmt.Println("Cron-Mini-ATM - Cookie odoo empty!")
			helper.ErrorLog("Cron-Mini-ATM - Cookie odoo empty!")
			ok = false
		}

		err = helper.SendToOdoo(cfg.CnfGlob.OdooURL+"/iid_api_manage/post_data", name, cookie, "Transaction", code, "Arthajasa", jsonString)
		if err != nil {
			tx.Rollback()
			fmt.Println("Cron-Mini-ATM - Send to odoo : " + err.Error())
			helper.ErrorLog("Cron-Mini-ATM - Send to odoo  : " + err.Error())
			ok = false
		}

		if ok {
			err = s.repo.MiniAtmUpdateFlagOdoo(ctx, tx, &model.MiniAtm{
				ID:       row.ID,
				FlagOdoo: 1,
			})
			if err != nil {
				tx.Rollback()
				fmt.Println("Cron-Mini-ATM - Update flag odoo : " + err.Error())
				helper.ErrorLog("Cron-Mini-ATM - Update flag odoo  : " + err.Error())
				break
			}

			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}
