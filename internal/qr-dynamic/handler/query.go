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

func (s *Service) Query(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) {

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

	// cookie, err := h.AuthenticateOdoo(s.config.CnfGlob.OdooURL + "/web/session/authenticate")
	// if err != nil {
	// 	h.ErrorLog("QR-Dynamic-Query - Get cookie odoo : "+err.Error(), "qr_dynamic")
	// 	return nil, status.Errorf(codes.Internal, "Service malfunction, code : O3")
	// }

	// if cookie == "" {
	// 	h.ErrorLog("QR-Dynamic-Query - Cookie odoo empty !", "qr_dynamic")
	// 	return nil, status.Errorf(codes.Internal, "Service malfunction, code : O4")
	// }

	res, err := h.CheckCodeTrxJournal(hostCode, s.config.CnfGlob.OdooURL+"/iid_api_manage")
	if err != nil {
		h.ErrorLog("QR-Dynamic-Query - Check code trx journal : "+err.Error(), "qr_dynamic")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : O2")
	}

	if res != "ok" {
		return nil, status.Errorf(codes.InvalidArgument, "Host-Code not registered!")
	}

	reqData, _ := json.Marshal(data)

	payload := string(reqData)
	if contentType == "application/x-www-form-urlencoded" {
		values := url.Values{}
		for key, val := range data {
			values.Set(key, val.(string))
		}

		payload = values.Encode()
	}

	hostAddressQuery, err := s.jackdbParamService.GetAddressQueryByName(ctx, host)
	if err != nil {
		h.ErrorLog("QR-Dynamic-Query - Check host : "+err.Error(), "qr_dynamic")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : H0")
	}

	if hostAddressQuery == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Host not found!")
	}

	exReq, err := http.NewRequest(method, hostAddressQuery, strings.NewReader(payload))
	if err != nil {
		h.ErrorLog("QR-Dynamic-Query - Prepare send to host : "+err.Error(), "qr_dynamic")
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
		h.ErrorLog("QR-Dynamic-Query - Send to host : "+err.Error(), "qr_dynamic")
		return nil, status.Errorf(codes.Internal, "Service malfunction, code : N1")
	}

	defer response.Body.Close()

	var responseByte []byte
	if response.Header.Get("Content-Encoding") != "" {
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			h.ErrorLog("QR-Dynamic-Query - Read gzip encode : "+err.Error(), "qr_dynamic")
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : G0")
		}
		defer reader.Close()
		responseByte, err = io.ReadAll(reader)
		if err != nil {
			h.ErrorLog("QR-Dynamic-Query - Read response host : "+err.Error(), "qr_dynamic")
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : R0")
		}
	} else {
		responseByte, err = io.ReadAll(response.Body)
		if err != nil {
			h.ErrorLog("QR-Dynamic-Query - Read response host : "+err.Error(), "qr_dynamic")
			return nil, status.Errorf(codes.Internal, "Service malfunction, code : R1")
		}
	}

	originResponseData := string(responseByte)

	result := map[string]interface{}{
		"status":       response.StatusCode,
		"content_type": response.Header.Get("Content-Type"),
		"body":         originResponseData,
	}

	return h.RespondGrpc(result), nil
}
