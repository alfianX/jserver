package model

import "time"

type QrDynamicCallback struct {
	ID           int64     `json:"id"`
	Host         string    `json:"host"`
	RequestData  string    `json:"request_data"`
	ResponseData string    `json:"response_data"`
	Flag         int64     `json:"flag"`
	CreatedAt    time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (QrDynamicCallback) TableName() string {
	return "qr_dynamic_callback"
}

type QrDynamicGenerate struct {
	ID           int64     `json:"id"`
	Host         string    `json:"host"`
	RequestData  string    `json:"request_data"`
	ResponseData string    `json:"response_data"`
	Flag         int64     `json:"flag"`
	CreatedAt    time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (QrDynamicGenerate) TableName() string {
	return "qr_dynamic_generate"
}

type QrNotif struct {
	ID           int64     `json:"id"`
	RequestData  string    `json:"request_data"`
	ResponseData string    `json:"response_data"`
	Flag         int64     `json:"flag"`
	CreatedAt    time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (QrNotif) TableName() string {
	return "qr_notif"
}

type MiniAtm struct {
	ID                      int64     `json:"id"`
	SerialNumber            string    `json:"serial_number"`
	MerchantGroup           string    `json:"merchant_group"`
	TrxType                 string    `json:"trx_type"`
	Mti                     string    `json:"mti"`
	Procode                 string    `json:"procode"`
	Stan                    string    `json:"stan"`
	Trace                   string    `json:"trace"`
	Tid                     string    `json:"tid"`
	Mid                     string    `json:"mid"`
	Pan                     string    `json:"pan"`
	Amount                  int64     `json:"amount"`
	TransactionDate         time.Time `gorm:"autoCreateTime:false" json:"transaction_date"`
	Nii                     string    `json:"nii"`
	De63                    string    `json:"de63"`
	ResponseCode            string    `json:"response_code"`
	TransactionDateResponse time.Time `gorm:"autoCreateTime:false" json:"transaction_date_response"`
	Rrn                     string    `json:"rrn"`
	Host                    string    `json:"host"`
	FlagOdoo                int64     `json:"flag_odoo"`
	CreatedAt               time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt               time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (MiniAtm) TableName() string {
	return "mini_atm"
}

type CardPayment struct {
	ID                      int64     `json:"id"`
	SerialNumber            string    `json:"serial_number"`
	TrxType                 string    `json:"trx_type"`
	Mti                     string    `json:"mti"`
	Procode                 string    `json:"procode"`
	Stan                    string    `json:"stan"`
	Trace                   string    `json:"trace"`
	Tid                     string    `json:"tid"`
	Mid                     string    `json:"mid"`
	Pan                     string    `json:"pan"`
	Amount                  int64     `json:"amount"`
	TransactionDate         time.Time `gorm:"autoCreateTime:false" json:"transaction_date"`
	Nii                     string    `json:"nii"`
	De63                    string    `json:"de63"`
	ResponseCode            string    `json:"response_code"`
	TransactionDateResponse time.Time `gorm:"autoCreateTime:false" json:"transaction_date_response"`
	Rrn                     string    `json:"rrn"`
	Host                    string    `json:"host"`
	FlagOdoo                int64     `json:"flag_odoo"`
	CreatedAt               time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt               time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (CardPayment) TableName() string {
	return "card_payment"
}
