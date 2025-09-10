package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alfianX/jserver/pkg/rsa"
	"github.com/fir1/rest-api/pkg/erru"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/protobuf/types/known/structpb"
)

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

func Respond(c *gin.Context, data interface{}, status int) {
	var respData interface{}
	switch v := data.(type) {
	case nil:
	case erru.ErrArgument:
		status = http.StatusBadRequest
		respData = ErrorResponse{ErrorMessage: v.Unwrap().Error()}
	case error:
		if http.StatusText(status) == "" {
			status = http.StatusInternalServerError
		} else {
			respData = ErrorResponse{ErrorMessage: v.Error()}
		}
	default:
		respData = data
	}

	c.JSON(status, respData)
}

func RespondGrpc(input interface{}) *structpb.Struct {
	var result map[string]interface{}
	data, _ := json.Marshal(input)

	_ = json.Unmarshal(data, &result)

	str, _ := structpb.NewStruct(result)

	return str
}

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Min is " + fe.Param()
	case "max":
		return "Max is " + fe.Param()
	case "numeric":
		return "Should be numeric"
	}
	return "Unknown error"
}

func Decode(c *gin.Context, v interface{}) error {
	if err := c.ShouldBindJSON(v); err != nil {
		return err
	}
	return nil
}

func ReadRequestBody(c *gin.Context) ([]byte, error) {
	var bodyBytes []byte
	var err error

	if c.Request.Body != nil {
		bodyBytes, err = io.ReadAll(c.Request.Body)
		if err != nil {
			err := errors.New("could not read request body")
			return nil, err
		}
	}
	return bodyBytes, nil
}

func RestoreRequestBody(c *gin.Context, bodyBytes []byte) {
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}

func SendMessageToHsm(IPPORT, message string) (string, error) {
	iso, _ := hex.DecodeString(message)

	tcpServer, err := net.ResolveTCPAddr("tcp", IPPORT)
	if err != nil {
		return "", err
	}

	conn, err := net.DialTCP("tcp", nil, tcpServer)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write(iso)
	if err != nil {
		return "", err
	}

	received := make([]byte, 1024)
	bytesRead, err := conn.Read(received)
	if err != nil {
		return "", err
	}

	conn.Close()

	messageHost := hex.EncodeToString(received[:bytesRead])

	return messageHost, nil
}

func MaskPan(pan string) string {
	length := len(pan)
	visibleCount := length / 4
	hiddenCount := length - (visibleCount * 2)

	mask := pan[:visibleCount] + strings.Repeat("*", hiddenCount) + pan[length-visibleCount:]

	return mask
}

func HSMEncrypt(IPPORT string, zek string, data string) (string, error) {
	var lenDataHex string
	lenData := len(data)
	if lenData > 0 {
		lenDataHexT := strings.ToUpper(fmt.Sprintf("%04x", lenData))
		extra := lenData % 16
		if extra > 0 {
			padSize := 16 - extra
			lenData = lenData + padSize
			lenDataHex = strings.ToUpper(fmt.Sprintf("%04x", lenData))
			data = data + strings.Repeat("0", padSize)
		} else {
			lenDataHex = strings.ToUpper(fmt.Sprintf("%04x", lenData))
		}

		command := "RSWIM0001100A" + zek + lenDataHex + data
		len := len(command)
		lenHex := strings.ToUpper(fmt.Sprintf("%04x", len))
		message := lenHex + hex.EncodeToString([]byte(command))

		response, err := SendMessageToHsm(IPPORT, message)
		if err != nil {
			return "", err
		}

		resByte, err := hex.DecodeString(response)
		if err != nil {
			return "", err
		}

		hsmResponse := string(resByte[8:10])
		if hsmResponse != "00" {
			return "", errors.New("HSM response invalid")
		}

		return lenDataHexT + string(resByte[14:]), nil
	} else {
		return "", errors.New("no data")
	}
}

func HSMDecrypt(IPPORT string, zek string, data string) (string, error) {
	var lenDataHex string
	lenData := len(data)
	if lenData > 0 {
		lenDataHex = strings.ToUpper(fmt.Sprintf("%04x", lenData))
		command := "RSWIM2001100A" + zek + lenDataHex + data
		len := len(command)
		lenHex := strings.ToUpper(fmt.Sprintf("%04x", len))
		message := lenHex + hex.EncodeToString([]byte(command))

		response, err := SendMessageToHsm(IPPORT, message)
		if err != nil {
			return "", err
		}

		resByte, err := hex.DecodeString(response)
		if err != nil {
			return "", err
		}

		hsmResponse := string(resByte[8:10])
		if hsmResponse != "00" {
			return "", errors.New("HSM response invalid")
		}

		return string(resByte[14:]), nil
	} else {
		return "", nil
	}
}

func CreateSignature(tid, mid, email, transactionDate, trace, approvalCode string) (string, error) {
	filename := "rsa/rsa_private.key"
	pem, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	pemStr := string(pem)
	privateKey, err := rsa.ParseRsaPrivateKeyFromPemStr(pemStr)
	if err != nil {
		return "", err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	trxDate, err := time.ParseInLocation("2006-01-02 15:04:05", transactionDate, loc)
	if err != nil {
		return "", err
	}
	trxDateSig := trxDate.Format("20060102150405")

	dataSignature := []byte(tid + mid + email + trxDateSig + approvalCode + trace)
	signature, err := rsa.CreateSignature(privateKey, dataSignature)
	if err != nil {
		return "", err
	}

	signatureFinal := base64.StdEncoding.EncodeToString(signature)

	return signatureFinal, nil
}
