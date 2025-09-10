package types

type CardGenerateQrRequest struct {
	TransactionID string `json:"transactionId" binding:"required"`
	MerchantID    string `json:"merchantId" binding:"required"`
	CardType      string `json:"cardType" binding:"required"`
	Acquirer      string `json:"acquirer" binding:"required"`
	Amount        int64  `json:"amount" binding:"required"`
}

type CardTrxRequest struct {
	TransactionID string `json:"transactionId" binding:"required"`
	CardType      string `json:"cardType" binding:"required"`
	Acquirer      string `json:"acquirer" binding:"required"`
	ISO8583       string `json:"ISO8583" binding:"required"`
}

type StatusCardTrxRequest struct {
	TransactionID string `json:"transactionId" binding:"required"`
}
