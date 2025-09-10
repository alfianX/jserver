package handler

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	h "github.com/alfianX/jserver/helper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Service) GenerateQr(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {

	md, _ := metadata.FromIncomingContext(ctx)

	method := getFirstOrDefault(md.Get("x-http-method"), "POST")
	contentType := getFirstOrDefault(md.Get("x-content-type"), "application/json")
	// contentLength := getFirstOrDefault(md.Get("x-content-length"), "")

	data := req.AsMap()

	hostQRA := md.Get("host-qr")
	hostCodeA := md.Get("Host-Code")
	var host string
	var hostCode string

	if hostQRA != nil {
		host = strings.ToUpper(hostQRA[0])
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

	id, err := s.jackdbService.QrGenerateSave(ctx, host, requestData)
	if err != nil {
		h.ErrorLog("QR-Dynamic-Generate - Save data : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : S0")
	}

	payload := string(reqData)
	if contentType == "application/x-www-form-urlencoded" {
		values := url.Values{}
		for key, val := range data {
			values.Set(key, val.(string))
		}

		payload = values.Encode()
	}

	hostAddress, err := s.jackdbParamService.GetAddressByName(ctx, host)
	if err != nil {
		h.ErrorLog("QR-Dynamic-Generate - Check host : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : H0")
	}

	if hostAddress == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Host not found!")
	}

	exReq, err := http.NewRequest(method, hostAddress, strings.NewReader(payload))
	if err != nil {
		h.ErrorLog("QR-Dynamic-Generate - Prepare send to host : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : N0")
	}

	exReq.Header.Set("Content-Type", contentType)
	// exReq.Header.Set("Content-Length", contentLength)

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	response, err := client.Do(exReq)
	if err != nil {
		if strings.Contains(err.Error(), "Timeout") {
			return nil, status.Errorf(codes.Internal, "request timeout!")
		}
		h.ErrorLog("QR-Dynamic-Generate - Send to host : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : N1")
	}

	defer response.Body.Close()

	var responseByte []byte
	if response.Header.Get("Content-Encoding") != "" {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			h.ErrorLog("QR-Dynamic-Generate - Read gzip encode : " + err.Error())
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : G0")
		}
		defer reader.Close()
		responseByte, err = io.ReadAll(reader)
		if err != nil {
			h.ErrorLog("QR-Dynamic-Generate - Read response host : " + err.Error())
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : R0")
		}
	} else {
		responseByte, err = io.ReadAll(response.Body)
		if err != nil {
			h.ErrorLog("QR-Dynamic-Generate - Read response host : " + err.Error())
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : R1")
		}
	}

	originResponseData := string(responseByte)
	responseData := strings.ReplaceAll(originResponseData, "\n", "")
	responseData = strings.ReplaceAll(responseData, " ", "")

	err = s.jackdbService.QrGenerateUpdateResponse(ctx, responseData, id)
	if err != nil {
		h.ErrorLog("QR-Dynamic-Generate - Update data : " + err.Error())
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : U0")
	}

	result := map[string]interface{}{
		"status":       response.StatusCode,
		"content_type": response.Header.Get("Content-Type"),
		"body":         originResponseData,
	}

	return h.RespondGrpc(result), nil
}

func getFirstOrDefault(values []string, defaultVal string) string {
	if len(values) > 0 && values[0] != "" {
		return values[0]
	}
	return defaultVal
}
