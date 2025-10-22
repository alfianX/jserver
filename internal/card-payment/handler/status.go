package handler

import (
	"context"

	h "github.com/alfianX/jserver/helper"
	card_payment "github.com/alfianX/jserver/proto/card-payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) Status(ctx context.Context, req *card_payment.RequestStatus) (*card_payment.ResponseStatus, error) {
	if req.SerialNumber == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Mandatory Field {serialNumber}")
	}

	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Mandatory Field {id}")
	}

	err := s.jackdbService.CardPaymentUpdateFlagSuccess(ctx, req.Id)
	if err != nil {
		h.ErrorLog("Card-Payment-Status - Update flag success : " + err.Error())
		return nil, status.Errorf(codes.Internal, "General Error [S0]")
	}

	return &card_payment.ResponseStatus{Status: "SUCCESS", Message: "Update success"}, nil
}
