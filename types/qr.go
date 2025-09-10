package types

type QrDynamicCallbackReq struct {
	Host string                 `json:"host"`
	Data map[string]interface{} `json:"data"`
}
