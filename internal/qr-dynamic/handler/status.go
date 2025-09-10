package handler

import (
	"context"
	"encoding/json"
	"strings"

	h "github.com/alfianX/jserver/helper"
	"github.com/alfianX/jserver/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Service) Status(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	response := types.ResponseDynamicStatus{}

	md, _ := metadata.FromIncomingContext(ctx)

	data := req.AsMap()

	hostQRA := md.Get("host-qr")
	hostCodeA := md.Get("Host-Code")
	var host string
	var hostCode string

	if hostQRA != nil {
		host = hostQRA[0]
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "missing header Host-Qr!")
	}

	if hostCodeA != nil {
		hostCode = hostCodeA[0]
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "missing header Host-Code!")
	}

	cookie, err := h.AuthenticateOdoo(s.config.CnfGlob.OdooURL + "/web/session/authenticate")
	if err != nil {
		h.ErrorLog("QR-Dynamic - Get cookie odoo : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O3")
	}

	if cookie == "" {
		h.ErrorLog("QR-Dynamic - Cookie odoo empty !")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O4")
	}

	res, err := h.CheckCodeTrxJournal(hostCode, s.config.CnfGlob.OdooURL+"/iid_api_manage")
	if err != nil {
		h.ErrorLog("QR-Dynamic - Check code trx journal : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O2")
	}

	if res != "ok" {
		return nil, status.Errorf(codes.InvalidArgument, "Host-Code not registered!")
	}

	reqData, _ := json.Marshal(data)
	requestData := string(reqData)
	requestData = strings.ReplaceAll(requestData, "\n", "")
	requestData = strings.ReplaceAll(requestData, " ", "")

	id, err := s.jackdbService.QrCallbackSave(ctx, host, requestData)
	if err != nil {
		h.ErrorLog("QR-Dynamic - Save data : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : S0")
	}

	response.Status = "SUCCESS"
	response.Message = "Callback success"
	response.Data = data

	respB, _ := json.Marshal(response)
	responseData := string(respB)
	responseData = strings.ReplaceAll(responseData, "\n", "")
	responseData = strings.ReplaceAll(responseData, " ", "")

	err = s.jackdbService.QrCallbackUpdateResponse(ctx, responseData, id)
	if err != nil {
		h.ErrorLog("QR-Dynamic - Update data : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : U0")
	}

	err = h.SendToOdoo(s.config.CnfGlob.OdooURL+"/iid_api_manage/post_data", host, cookie, "Transaction", hostCode, host, requestData)
	if err != nil {
		h.ErrorLog("QR-Dynamic - Send to odoo : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O5")
	}

	return h.RespondGrpc(response), nil
}
