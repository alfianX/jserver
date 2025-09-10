buildGatewayLinux:
	go build -o ./bin/jserver/gateway/jserver-gateway ./cmd/gateway/main.go

buildQrNotifLinux:
	go build -o ./bin/jserver/qr-notif/jserver-qr-notif ./cmd/qr-notif/main.go

buildQrDynamicLinux:
	go build -o ./bin/jserver/qr-dynamic/jserver-qr-dynamic ./cmd/qr-dynamic/main.go

buildMiniAtmLinux:
	go build -o ./bin/jserver/mini-atm/jserver-mini-atm ./cmd/mini-atm/main.go

buildCardPaymentLinux:
	go build -o ./bin/jserver/card-payment/jserver-card-payment ./cmd/card-payment/main.go