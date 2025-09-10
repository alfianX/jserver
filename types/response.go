package types

type ResponseCard struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ISO8583 string `json:"ISO8583"`
}

type ResponseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseDynamicStatus struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
