package model

import "time"

type HostQr struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	AddressQuery string    `json:"address_query"`
	Status       int64     `json:"status"`
	CreatedAt    time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (HostQr) TableName() string {
	return "host_qr"
}

type Services struct {
	ID          int64     `json:"id"`
	HttpPrefix  string    `json:"http_prefix"`
	HttpMethod  string    `json:"http_method"`
	GrpcAddress string    `json:"grpc_address"`
	GrpcMethod  string    `json:"grpc_method"`
	GrpcType    int64     `json:"grpc_type"`
	Auth        int64     `json:"auth"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

type HostMiniatm struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Status    int64     `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (HostMiniatm) TableName() string {
	return "host_miniatm"
}

type HostDebit struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Status    int64     `json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:false" json:"updated_at"`
}

func (HostDebit) TableName() string {
	return "host_debit"
}

type CodeOdoo struct {
	ID        int64     `json:"id"`
	Host      string    `json:"host"`
	TypeTrx   string    `json:"type_trx"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `gorm:"autoCreateTime:false" json:"created_at"`
}

func (CodeOdoo) TableName() string {
	return "code_odoo"
}
