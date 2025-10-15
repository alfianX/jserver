package handler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	h "github.com/alfianX/jserver/helper"
	"github.com/alfianX/jserver/types"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var secretKey = "Yk0XV9lmTg2k53klhP3cxp5JmK8VewcF"

var allowedHeaders = make(map[string]bool)

func (s *service) Actions(c *gin.Context) {
	responseErr := types.ResponseError{}

	// relativePath := c.FullPath()
	// timeStamp := c.GetHeader("X-TIMESTAMP")
	// signature := c.GetHeader("X-SIGNATURE")

	path := strings.TrimPrefix(c.Param("action"), "/")

	srv, err := s.jackdbParamService.GetServices(c, path)
	if err != nil {
		h.ErrorLog("Get url microservice: " + err.Error())
		responseErr.Status = "SERVER_FAILED"
		responseErr.Message = "Service Malfunction"
		h.Respond(c, responseErr, http.StatusInternalServerError)
		return
	}

	if c.Request.Method != srv.HttpMethod {
		responseErr.Status = "SERVER_FAILED"
		responseErr.Message = "Method Not Allowed. This endpoint only accepts " + srv.HttpMethod + " requests."
		h.Respond(c, responseErr, http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.ErrorLog("Get body: " + err.Error())
		responseErr.Status = "SERVER_FAILED"
		responseErr.Message = "Service Malfunction"
		h.Respond(c, responseErr, http.StatusInternalServerError)
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var bodyMap map[string]interface{}
	var bodyString string
	ct := c.ContentType()
	if strings.HasPrefix(ct, "application/json") {
		c.ShouldBindJSON(&bodyMap)
		var minified bytes.Buffer
		err = json.Compact(&minified, bodyBytes)
		if err != nil {
			h.ErrorLog("json minify failed: " + err.Error())
			responseErr.Status = "SERVER_FAILED"
			responseErr.Message = "Service Malfunction"
			h.Respond(c, responseErr, http.StatusInternalServerError)
			return
		}
		bodyString = minified.String()
	} else if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
		form := map[string]interface{}{}
		c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			form[k] = v[0]
		}
		bodyMap = form

		// re := regexp.MustCompile(`\r?\n`)
		// reqMessage := re.ReplaceAllString(string(bodyBytes), "")
		bodyString = string(bodyBytes)
	}

	if srv.Auth == 1 {
		fmt.Println("body to sha256:" + bodyString)
		relativePath := c.Request.URL.Path
		timeStamp := c.GetHeader("X-TIMESTAMP")
		fmt.Println("timestamp luar:" + timeStamp)
		if timeStamp == "" {
			responseErr.Status = "INVALID_REQUEST"
			responseErr.Message = "header X-TIMESTAMP empty!"
			h.Respond(c, responseErr, http.StatusBadRequest)
			return
		}
		const layout = "2006-01-02T15:04:05Z07:00"
		reqTime, err := time.Parse(layout, timeStamp)
		if err != nil {
			responseErr.Status = "INVALID_REQUEST"
			responseErr.Message = "invalid timestamp format!"
			h.Respond(c, responseErr, http.StatusBadRequest)
			return
		}
		now := time.Now().In(reqTime.Location())
		maxDiff := 5 * time.Minute

		fmt.Println("date sekarang server:" + now.String())

		diff := now.Sub(reqTime)
		if diff > maxDiff {
			responseErr.Status = "INVALID_REQUEST"
			responseErr.Message = "timestamp is too old!"
			h.Respond(c, responseErr, http.StatusBadRequest)
			return
		}

		if diff < -(2 * time.Minute) {
			responseErr.Status = "INVALID_REQUEST"
			responseErr.Message = "timestamp is in the future!"
			h.Respond(c, responseErr, http.StatusBadRequest)
			return
		}

		signature := c.GetHeader("X-SIGNATURE")
		if signature == "" {
			responseErr.Status = "INVALID_REQUEST"
			responseErr.Message = "header X-SIGNATURE empty!"
			h.Respond(c, responseErr, http.StatusBadRequest)
			return
		}
		fmt.Println("signature luar: " + signature)
		hasher := sha256.New()
		hasher.Write([]byte(bodyString))
		hashSum := hasher.Sum(nil)
		bodyHash := strings.ToLower(hex.EncodeToString(hashSum))

		toSign := fmt.Sprintf("%s:%s:%s:%s", c.Request.Method, relativePath, bodyHash, timeStamp)
		fmt.Println("yg mau di hmac: " + toSign)

		hm := hmac.New(sha512.New, []byte(secretKey))
		hm.Write([]byte(toSign))
		hmacSum := hm.Sum(nil)
		hmacBase64 := base64.StdEncoding.EncodeToString(hmacSum)

		if hmacBase64 != signature {
			responseErr.Status = "INVALID_REQUEST"
			responseErr.Message = "signature not valid"
			h.Respond(c, responseErr, http.StatusBadRequest)
			return
		}
	}

	dataHeader, err := s.jackdbParamService.AllowedHeadersGetHeaderName(c)
	if err != nil {
		h.ErrorLog("Get header name: " + err.Error())
		responseErr.Status = "SERVER_FAILED"
		responseErr.Message = "Service Malfunction"
		h.Respond(c, responseErr, http.StatusInternalServerError)
		return
	}

	for _, headerName := range dataHeader {
		allowedHeaders[headerName.HeaderName] = true
	}

	md := metadata.MD{}
	for k, v := range c.Request.Header {
		key := strings.ToLower(k)
		if allowedHeaders[key] {
			md.Set(key, v...)
		}
	}

	md.Set("x-content-type", c.ContentType())
	md.Set("x-http-method", c.GetHeader("X-Http-Method"))
	// md.Set("x-content-length", c.GetHeader("Content-Length"))

	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	var respJSON []byte
	if srv.GrpcType == 1 {
		respJSON, err = s.ProxyRequest(ctx, srv.GrpcAddress, srv.GrpcMethod, bodyMap)
	} else {
		respJSON, err = s.ProxyRequestDynamic(ctx, srv.GrpcAddress, srv.GrpcMethod, bodyMap)
	}
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			h.ErrorLog("microservice failed: " + st.Message())
			httpCode := mapGrpcCodeToHTTP(st.Code())
			responseErr.Status = "SERVER_FAILED"
			responseErr.Message = st.Message()
			lenMsg := len(st.Message())
			index := strings.Index(st.Message(), "- id:")
			if index != -1 {
				fmt.Println(st.Message()[index+5 : lenMsg])
				idStr := st.Message()[index+5 : lenMsg]
				id, err := strconv.Atoi(idStr)
				if err != nil {
					responseErr.Message = err.Error()
					h.Respond(c, responseErr, httpCode)
					return
				}

				h.Respond(c, gin.H{"id": id, "status": "SERVER_FAILED", "message": st.Message()[:index]}, httpCode)
				return
			} else {
				h.Respond(c, responseErr, httpCode)
				return
			}
		}
		h.ErrorLog("microservice failed: " + err.Error())
		responseErr.Status = "SERVER_FAILED"
		responseErr.Message = "Service Malfunction"
		h.Respond(c, responseErr, http.StatusInternalServerError)
		return
	}

	var res map[string]interface{}
	json.Unmarshal(respJSON, &res)

	if ct, ok := res["content_type"].(string); ok {
		if strings.HasPrefix(ct, "text/html") {
			html := res["html"].(string)
			c.Data(http.StatusOK, ct, []byte(html))
			return
		} else {
			var bodyRes map[string]interface{}
			json.Unmarshal([]byte(res["body"].(string)), &bodyRes)
			h.Respond(c, bodyRes, int(res["status"].(float64)))
			return
		}
	}

	h.Respond(c, res, http.StatusOK)
}

func mapGrpcCodeToHTTP(code codes.Code) int {
	switch code {
	case codes.OK:
		return 200
	case codes.NotFound:
		return 404
	case codes.InvalidArgument:
		return 400
	case codes.AlreadyExists:
		return 409
	case codes.PermissionDenied:
		return 403
	case codes.Unauthenticated:
		return 401
	case codes.ResourceExhausted:
		return 429
	case codes.FailedPrecondition:
		return 412
	case codes.Internal:
		return 500
	default:
		return 500
	}
}
