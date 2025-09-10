package types

type MiniAtmOdoo struct {
	SerialNumber            string `json:"serial_number"`
	TrxType                 string `json:"trx_type"`
	Mti                     string `json:"mti"`
	Procode                 string `json:"procode"`
	Stan                    string `json:"stan"`
	Trace                   string `json:"trace"`
	Tid                     string `json:"tid"`
	Mid                     string `json:"mid"`
	ExternalStoreId         string `json:"externalStoreId"`
	Pan                     string `json:"pan"`
	Amount                  int64  `json:"amount"`
	PaymentDate             string `json:"paymentDate"`
	Nii                     string `json:"nii"`
	De63                    string `json:"de63"`
	ResponseCode            string `json:"response_code"`
	TransactionDateResponse string `json:"transaction_date_response"`
	PaymentReferenceNo      string `json:"paymentReferenceNo"`
}
