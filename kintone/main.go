package kintone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Config struct {
	Subdomain string
	API_TOKEN string
	API_ID    string
}

func readConf() Config {
	file, _ := os.Open("kintone/conf.json")
	defer file.Close()
	var conf Config
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&conf)
	return conf
}

func getHeader() map[string]string {
	conf := readConf()
	header := map[string]string{
		"X-Cybozu-API-Token": conf.API_TOKEN,
	}
	return header
}

func Request(urlStr string, method string, params []byte) []byte {
	header := getHeader()
	if method != "GET" {
		header["Content-Type"] = "application/json"
	}
	client := &http.Client{}
	req, _ := http.NewRequest(method, urlStr, nil)
	for key, value := range header {
		req.Header.Add(key, value)
	}
	if method != "GET" {
		req.Body = io.NopCloser(bytes.NewBuffer(params))
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	return resBody
}
func GetRecord(id int) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json?app=%s&id=%d", conf.Subdomain, conf.API_ID, id)
	res := Request(urlStr, "GET", nil)
	return res
}

func GetAllRecords() []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records.json?app=%s", conf.Subdomain, conf.API_ID)
	res := Request(urlStr, "GET", nil)
	return res
}

func PostRecord(paramMap map[string]interface{}) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json", conf.Subdomain)

	params := map[string]interface{}{
		"app": conf.API_ID,
		"record": func() map[string]interface{} {
			record := make(map[string]interface{})
			for k, v := range paramMap {
				record[k] = map[string]interface{}{
					"value": v,
				}
			}
			return record
		}(),
	}
	jsonParams, _ := json.Marshal(params)
	res := Request(urlStr, "POST", jsonParams)

	return res
}

func UpdateRecord(id int, paramMap map[string]interface{}) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json", conf.Subdomain)

	params := map[string]interface{}{
		"app": conf.API_ID,
		"id":  id,
		"record": func() map[string]interface{} {
			record := make(map[string]interface{})
			for k, v := range paramMap {
				record[k] = map[string]interface{}{
					"value": v,
				}
			}
			return record
		}(),
	}
	jsonParams, _ := json.Marshal(params)
	res := Request(urlStr, "PUT", jsonParams)

	return res
}

func DeleteRecord(ids []int) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records.json", conf.Subdomain)

	params := map[string]interface{}{
		"app": conf.API_ID,
		"ids": ids,
	}
	jsonParams, _ := json.Marshal(params)
	res := Request(urlStr, "DELETE", jsonParams)

	return res
}

type Cursor struct {
	Id    string `json:"id"`
	Total int    `json:"totalCount"`
}

func IsExisted(date string) bool {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records/cursor.json", conf.Subdomain)
	params := map[string]interface{}{
		"app":    conf.API_ID,
		"fields": []string{"レコード番号"},
		"query":  fmt.Sprintf("date=\"%s\"", date),
		"size":   500,
	}
	jsonParams, _ := json.Marshal(params)
	res := Request(urlStr, "POST", jsonParams)
	var cursor Cursor
	json.Unmarshal(res, &cursor)
	fmt.Println(cursor.Total)
	return cursor.Total > 0
}

func GetRecordByDate(fromDate string, toDate string) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records/cursor.json", conf.Subdomain)
	params := map[string]interface{}{
		"app":    conf.API_ID,
		"fields": []string{"レコード番号", "date", "sleep", "hello"},
		"query":  fmt.Sprintf("date>=\"%s\" and date<=\"%s\"", fromDate, toDate),
		"size":   500,
	}
	jsonParams, _ := json.Marshal(params)
	res := Request(urlStr, "POST", jsonParams)
	var cursor Cursor
	json.Unmarshal(res, &cursor)
	fmt.Println(cursor.Total)
	urlStr = fmt.Sprintf("https://%s.cybozu.com/k/v1/records/cursor.json?id=%s", conf.Subdomain, cursor.Id)
	res = Request(urlStr, "GET", nil)
	return res
}
