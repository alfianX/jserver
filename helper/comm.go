package helper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/alfianX/jserver/config"
	"github.com/gin-gonic/gin"
)

func RestSendToHost(c *gin.Context, cnf config.Config, payload []byte, issuerService string) (map[string]interface{}, error) {
	var resp *http.Response

	exReq, err := http.NewRequest("POST", issuerService+"/sale", bytes.NewReader(payload))
	if err != nil {
		return nil, errors.New("Prepare send to microservice: " + err.Error())
	}

	exReq.Header = c.Request.Header

	client := &http.Client{
		Timeout: time.Duration(cnf.CnfGlob.TimeoutTrx) * time.Second,
	}
	resp, err = client.Do(exReq)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var extResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&extResp)
	if err != nil {
		return nil, errors.New("Decode response: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Response: " + extResp["message"].(string))
	}

	return extResp, nil
}

func TcpSendToHost(IPPORT string, isoReq string, timeInt int) (string, error) {
	var isoRes []byte
	tcpServer, err := net.ResolveTCPAddr("tcp", IPPORT)
	if err != nil {
		return "", err
	}

	timeout := time.Duration(timeInt) * time.Second

	conn, err := net.DialTimeout("tcp", tcpServer.String(), timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	iso, _ := hex.DecodeString(isoReq)

	_, err = conn.Write(iso)
	if err != nil {
		return "", err
	}

	conn.SetReadDeadline(time.Now().Add(timeout))

	header := make([]byte, 3)
	_, err = io.ReadFull(conn, header)
	if err != nil {
		return "", err
	}

	headerStr := hex.EncodeToString(header)
	msqLength, _ := strconv.ParseInt(headerStr[:4], 16, 64)
	messageBytes := make([]byte, msqLength-1)

	TP := headerStr[4:]

	if TP != "60" {
		return "", errors.New("not iso data")
	}

	_, err = io.ReadFull(conn, messageBytes)
	if err != nil {
		return "", err
	}

	isoRes = append(header, messageBytes...)

	isoResStr := hex.EncodeToString(isoRes)

	return isoResStr, nil
}
