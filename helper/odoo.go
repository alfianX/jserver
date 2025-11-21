package helper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func AuthenticateOdoo(hostAddress string) (string, error) {
	type request struct {
		JsonRpc string `json:"jsonrpc"`
		Params  struct {
			Db       string `json:"db"`
			Login    string `json:"login"`
			Password string `json:"password"`
		} `json:"params"`
	}

	odooDB := os.Getenv("ODOO_DB")
	odooUser := os.Getenv("ODOO_USER")
	odooPassword := os.Getenv("ODOO_PASSWORD")

	data := request{}
	data.JsonRpc = "2.0"
	data.Params.Db = odooDB
	data.Params.Login = odooUser
	data.Params.Password = odooPassword

	req, _ := json.Marshal(data)

	exReq, err := http.NewRequest("POST", hostAddress, bytes.NewReader(req))
	if err != nil {
		return "", err
	}

	exReq.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{
		// Konfigurasi TLS
		TLSClientConfig: &tls.Config{
			// Ini adalah baris kunci untuk mematikan verifikasi sertifikat
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}
	response, err := client.Do(exReq)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	var session_id string
	cookie := response.Cookies()
	for _, c := range cookie {
		if c.Name == "session_id" {
			session_id = c.Value
		}
	}

	return session_id, nil
}

func SendToOdoo(hostAddress, name, categ, trxType, hostName, data string) error {

	type Body struct {
		Categ           string      `json:"categ"`
		Name            string      `json:"name"`
		TransactionType string      `json:"transaction_type"`
		Host            string      `json:"host"`
		Data            interface{} `json:"data"`
	}

	type sendOdoo struct {
		Model  string                 `json:"model"`
		Values map[string]interface{} `json:"values"`
	}

	send := sendOdoo{
		Model: "iid.api.manage",
		Values: map[string]interface{}{
			name: map[string]interface{}{
				"categ": categ,
				"body": Body{
					Categ:           categ,
					Name:            name,
					TransactionType: trxType,
					Host:            hostName,
					Data:            data,
				},
			},
		},
	}

	reqOdoo, _ := json.Marshal(send)
	fmt.Println("request to odoo" + string(reqOdoo))

	exReqOdoo, err := http.NewRequest("POST", hostAddress, bytes.NewReader(reqOdoo))
	if err != nil {
		return err
	}

	exReqOdoo.Header.Add("Content-Type", "application/json")
	// exReqOdoo.Header.Add("Cookie", "session_id="+session_id)
	// cookie := http.Cookie{}
	// cookie.Name = "session_id"
	// cookie.Value = session_id
	// exReqOdoo.AddCookie(&cookie)
	tr := &http.Transport{
		// Konfigurasi TLS
		TLSClientConfig: &tls.Config{
			// Ini adalah baris kunci untuk mematikan verifikasi sertifikat
			InsecureSkipVerify: true,
		},
	}

	clientOdoo := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}
	respOdoo, err := clientOdoo.Do(exReqOdoo)
	if err != nil {
		return err
	}

	defer respOdoo.Body.Close()

	responseByte, _ := io.ReadAll(respOdoo.Body)

	fmt.Println("response from odoo: " + string(responseByte))

	if strings.Contains(string(responseByte), "error") || strings.Contains(string(responseByte), "ERROR") || strings.Contains(string(responseByte), "Error") {
		return fmt.Errorf("%s", string(responseByte))
	}

	// currentTime := time.Now()
	// gmtFormat := "20060102"
	// dateString := currentTime.Format(gmtFormat)
	// filename := fmt.Sprintf("log/odoo_response_%s.log", dateString)
	// file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	// if err != nil {
	// 	return err
	// }
	// defer file.Close()

	// currentTime = time.Now()
	// gmtFormat = "15:04:05"
	// dateString = currentTime.Format(gmtFormat)
	// logRequest := fmt.Sprintf("[%s] - %s\n\n", dateString, string(responseByte))
	// file.WriteString(logRequest)

	return nil
}

func CheckCodeTrxJournal(code, hostAddress string) (string, error) {
	payload := map[string]interface{}{
		"model":  "iid.transaction.journal",
		"domain": [][]interface{}{{"code", "=", code}},
		"fields": []string{"name"},
	}

	reqOdoo, _ := json.Marshal(payload)

	exReqOdoo, err := http.NewRequest("GET", hostAddress, bytes.NewReader(reqOdoo))
	if err != nil {
		return "error", err
	}

	exReqOdoo.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{
		// Konfigurasi TLS
		TLSClientConfig: &tls.Config{
			// Ini adalah baris kunci untuk mematikan verifikasi sertifikat
			InsecureSkipVerify: true,
		},
	}

	clientOdoo := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}
	respOdoo, err := clientOdoo.Do(exReqOdoo)
	if err != nil {
		return "error", err
	}

	defer respOdoo.Body.Close()

	body, _ := io.ReadAll(respOdoo.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "error", err
	}

	var record string
	var resErr error

	// Traverse and check if records exist
	if result, ok := parsed["result"].([]interface{}); ok && len(result) > 0 {
		if dataObj, ok := result[0].(map[string]interface{}); ok {
			if _, ok := dataObj["records"]; ok {
				// ✅ records exist
				// fmt.Println("Records found:", records)
				record = "ok"
				resErr = nil
			} else {
				// ❌ records not found
				// fmt.Println("No 'records' field found")
				record = "error"
				resErr = errors.New("no 'records' field found")
			}
		}
	} else {
		// fmt.Println("Invalid or missing 'result'")
		record = "error"
		resErr = errors.New("invalid or missing 'result'")
	}

	return record, resErr
}

func CheckNameTrxJournal(name, hostAddress string) (string, error) {
	payload := map[string]interface{}{
		"model":  "iid.transaction.journal",
		"domain": [][]interface{}{{"name", "=", name}},
		"fields": []string{"code"},
	}

	reqOdoo, _ := json.Marshal(payload)

	exReqOdoo, err := http.NewRequest("GET", hostAddress, bytes.NewReader(reqOdoo))
	if err != nil {
		return "error", err
	}

	exReqOdoo.Header.Add("Content-Type", "application/json")

	tr := &http.Transport{
		// Konfigurasi TLS
		TLSClientConfig: &tls.Config{
			// Ini adalah baris kunci untuk mematikan verifikasi sertifikat
			InsecureSkipVerify: true,
		},
	}

	clientOdoo := &http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}
	respOdoo, err := clientOdoo.Do(exReqOdoo)
	if err != nil {
		return "error", err
	}

	defer respOdoo.Body.Close()

	body, _ := io.ReadAll(respOdoo.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "error", err
	}

	var record string
	var resErr error

	// Traverse and check if records exist
	if result, ok := parsed["result"].([]interface{}); ok && len(result) > 0 {
		if dataObj, ok := result[0].(map[string]interface{}); ok {
			if records, ok := dataObj["records"].([]interface{}); ok && len(records) > 0 {
				if recordMap, ok := records[0].(map[string]interface{}); ok {
					if code, ok := recordMap["code"].(string); ok {
						// ✅ Ambil code
						// fmt.Println("CODE:", code)
						record = code
						resErr = nil
					} else {
						record = "error"
						// resErr = errors.New("'code' field missing or not a string")
						resErr = errors.New("nothing")
					}
				} else {
					record = "error"
					// resErr = errors.New("record[0] is not an object")
					resErr = errors.New("nothing")
				}
			} else {
				record = "error"
				// resErr = errors.New("'records' field missing or empty")
				resErr = errors.New("nothing")
			}
		}
	} else {
		// fmt.Println("Invalid or missing 'result'")
		record = "error"
		// resErr = errors.New("invalid or missing 'result'")
		resErr = errors.New("nothing")
	}

	return record, resErr
}
