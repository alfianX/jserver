package handler

import (
	"context"
	"encoding/json"
	"strings"

	h "github.com/alfianX/jserver/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Service) Status(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	type response struct {
		ID      int64  `json:"id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	md, _ := metadata.FromIncomingContext(ctx)

	data := req.AsMap()

	hostQRA := md.Get("host-qr")
	hostCodeA := md.Get("Host-Code")
	var host string
	var hostCode string

	if hostQRA != nil {
		host = hostQRA[0]
	}

	if hostCodeA != nil {
		hostCode = hostCodeA[0]
	}

	var idReq int64
	if idFloat, ok := data["id"].(float64); ok {
		idReq = int64(idFloat)
	} else {
		h.ErrorLog("QR-Notif - Could not extract ID as float64 !")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : I0")
	}

	res, err := h.CheckNameTrxJournal(data["issuer"].(string), s.config.CnfGlob.OdooURL+"/iid_api_manage")
	if err != nil && err.Error() != "nothing" {
		h.ErrorLog("QR-Notif - Check code trx journal : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O1 - id:%d", idReq)
	}

	if res != "error" {
		hostCode = res
	} else {
		if hostCode == "" {
			return nil, status.Errorf(codes.InvalidArgument, "missing header Host-Code! - id:%d", idReq)
		}

		res, err := h.CheckCodeTrxJournal(hostCode, s.config.CnfGlob.OdooURL+"/iid_api_manage")
		if err != nil {
			h.ErrorLog("QR-Notif - Check code trx journal : " + err.Error())
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : O2 - id:%d", idReq)
		}

		if res != "ok" {
			return nil, status.Errorf(codes.InvalidArgument, "Host-Code not registered! - id:%d", idReq)
		}
	}

	cookie, err := h.AuthenticateOdoo(s.config.CnfGlob.OdooURL + "/web/session/authenticate")
	if err != nil {
		h.ErrorLog("QR-Notif - Get cookie odoo : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O3 - id:%d", idReq)
	}

	if cookie == "" {
		h.ErrorLog("QR-Notif - Cookie odoo empty !")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O4 - id:%d", idReq)
	}

	reqData, _ := json.Marshal(data)
	requestData := string(reqData)
	requestData = strings.ReplaceAll(requestData, "\n", "")
	requestData = strings.ReplaceAll(requestData, " ", "")

	id, err := s.jackdbService.QrNotifSave(ctx, requestData)
	if err != nil {
		h.ErrorLog("QR-Notif - Save data : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : S0 - id:%d", idReq)
	}

	resp := response{
		ID:      idReq,
		Status:  "SUCCESS",
		Message: "Callback success",
	}

	respB, _ := json.Marshal(resp)
	responseData := string(respB)
	responseData = strings.ReplaceAll(responseData, "\n", "")
	responseData = strings.ReplaceAll(responseData, " ", "")

	err = s.jackdbService.QrNotifUpdateResponse(ctx, responseData, id)
	if err != nil {
		h.ErrorLog("QR-Notif - Update data : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : U0 - id:%d", idReq)
	}

	err = h.SendToOdoo(s.config.CnfGlob.OdooURL+"/iid_api_manage/post_data", "soundbox", cookie, "Transaction", hostCode, host, requestData)
	if err != nil {
		h.ErrorLog("QR-Notif - Send to odoo : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O5 - id:%d", idReq)
	}

	return h.RespondGrpc(resp), nil
}
