package handler

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	servicedb "github.com/alfianX/jserver/database/jackdb/service_db"
	h "github.com/alfianX/jserver/helper"
	"github.com/alfianX/jserver/pkg/iso"
	mini_atm "github.com/alfianX/jserver/proto/mini-atm"
	"github.com/moov-io/iso8583"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *Service) Transaction(ctx context.Context, req *mini_atm.Request) (*mini_atm.Response, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	hostA := md.Get("host-name")
	mgH := md.Get("merchant-group")
	var host string
	var mg string

	if hostA != nil {
		host = strings.ToUpper(hostA[0])
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "missing header host-name!")
	}

	if mgH != nil {
		mg = strings.ToUpper(mgH[0])
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "missing header merchant-group!")
	}

	if req.Iso == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Mandatory Field {iso}")
	}

	if req.SerialNumber == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Mandatory Field {serialNumber}")
	}

	if req.TransactionType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Mandatory Field {transactionType}")
	}

	isoRequest := req.Iso
	serialNumber := req.SerialNumber
	trxType := strings.ToUpper(req.TransactionType)

	isoMsg := iso8583.NewMessage(iso.Spec87Hex)
	isoMsg.Unpack([]byte(isoRequest[14:]))

	iso8583.Describe(isoMsg, os.Stdout)

	mti, _ := isoMsg.GetMTI()
	procode, _ := isoMsg.GetString(3)
	de4, _ := isoMsg.GetString(4)
	var amount int64
	if de4 != "" {
		amount, _ = strconv.ParseInt(de4, 10, 64)
	}
	stan, _ := isoMsg.GetString(11)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	var trxDate time.Time
	de12, _ := isoMsg.GetString(12)
	de13, _ := isoMsg.GetString(13)
	if de12 != "" && de13 != "" {
		trxDateStr := strconv.Itoa(time.Now().Year()) + "-" + de13[:2] + "-" + de13[2:4] + " " + de12[:2] + ":" + de12[2:4] + ":" + de12[4:6]
		trxDate, _ = time.ParseInLocation("2006-01-02 15:04:05", trxDateStr, loc)
	} else {
		trxDate = time.Now()
	}
	nii, _ := isoMsg.GetString(24)
	de35, _ := isoMsg.GetString(35)
	de35s := strings.Split(de35, "D")
	pan := h.MaskPan(de35s[0])
	tid, _ := isoMsg.GetString(41)
	mid, _ := isoMsg.GetString(42)
	trace, _ := isoMsg.GetString(62)
	de63, _ := isoMsg.GetString(63)

	idTrx, err := s.jackdbService.MiniAtmSave(ctx, &servicedb.MiniAtmReqParam{
		SerialNumber:    serialNumber,
		MerchantGroup:   mg,
		TrxType:         trxType,
		Mti:             mti,
		Procode:         procode,
		Stan:            stan,
		Trace:           trace,
		Tid:             tid,
		Mid:             mid,
		Pan:             pan,
		Amount:          amount,
		TransactionDate: trxDate,
		Nii:             nii,
		De63:            de63,
		Host:            host,
	})
	if err != nil {
		h.ErrorLog("Mini-Atm - Save trx : "+err.Error(), "mini_atm")
		return nil, status.Errorf(codes.Internal, "General Error [S0]")
	}

	hostAddress, err := s.jackdbParamService.HostMiniatmGetAddress(ctx, host)
	if err != nil {
		h.ErrorLog("Mini-Atm - Check host : "+err.Error(), "mini_atm")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : H0")
	}

	if hostAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Host not found!")
	}

	isoResponse, err := h.TcpSendToHost(hostAddress, isoRequest, s.config.CnfGlob.TimeoutTrx)
	if err != nil {
		h.ErrorLog("Mini-Atm - Send to host : "+err.Error(), "mini_atm")
		return nil, status.Errorf(codes.Internal, "General Error [H0]")
	}

	isoMsg.Unpack([]byte(isoResponse[14:]))
	var trxDateRes time.Time
	de12res, _ := isoMsg.GetString(12)
	de13res, _ := isoMsg.GetString(13)
	if de12res != "" && de13res != "" {
		trxDateStr := strconv.Itoa(time.Now().Year()) + "-" + de13res[:2] + "-" + de13res[2:4] + " " + de12res[:2] + ":" + de12res[2:4] + ":" + de12res[4:6]
		trxDateRes, _ = time.ParseInLocation("2006-01-02 15:04:05", trxDateStr, loc)
	}
	responseCode, _ := isoMsg.GetString(39)
	rrn, _ := isoMsg.GetString(37)

	err = s.jackdbService.MiniAtmUpdateResponse(ctx, &servicedb.MiniAtmResParam{
		ID:                      idTrx,
		ResponseCode:            responseCode,
		TransactionDateResponse: trxDateRes,
		Rrn:                     rrn,
	})
	if err != nil {
		h.ErrorLog("Mini-Atm - Update trx : "+err.Error(), "mini_atm")
		return nil, status.Errorf(codes.Internal, "General Error [U0]")
	}

	return &mini_atm.Response{Status: "SUCCESS", Message: "Trx Success", Iso: isoResponse}, nil
}
