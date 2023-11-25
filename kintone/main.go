package kintone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Subdomain string `json:"subdomain"`
	APIToken  string `json:"api_token"`
	APIID     string `json:"api_id"`
}

func readConf() Config {
	file, err := os.Open("kintone/conf.json")
	var conf Config
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		_ = decoder.Decode(&conf)
	} else {
		conf.Subdomain = os.Getenv("SUBDOMAIN")
		conf.APIToken = os.Getenv("API_TOKEN")
		conf.APIID = os.Getenv("API_ID")
	}

	return conf
}

func getHeader() map[string]string {
	conf := readConf()
	header := map[string]string{
		"X-Cybozu-API-Token": conf.APIToken,
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
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json?app=%s&id=%d", conf.Subdomain, conf.APIID, id)
	fmt.Printf("try get url:%s", urlStr)
	res := Request(urlStr, "GET", nil)
	fmt.Printf("get result:%s", string(res))
	return res
}

func GetAllRecords() []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records.json?app=%s", conf.Subdomain, conf.APIID)
	fmt.Printf("try get url:%s", urlStr)
	res := Request(urlStr, "GET", nil)
	fmt.Printf("get result:%s", string(res))
	return res
}

func PostRecord(paramMap map[string]interface{}) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json", conf.Subdomain)

	params := map[string]interface{}{
		"app": conf.APIID,
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
	fmt.Printf("try post params:%s\n", jsonParams)
	res := Request(urlStr, "POST", jsonParams)
	fmt.Printf("update result:%s\n", res)

	return res
}

func UpdateRecord(id int, paramMap map[string]interface{}) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json", conf.Subdomain)

	params := map[string]interface{}{
		"app": conf.APIID,
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
	fmt.Printf("try update params:%s\n", jsonParams)
	res := Request(urlStr, "PUT", jsonParams)
	return res
}

func DeleteRecord(ids []int) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records.json", conf.Subdomain)

	params := map[string]interface{}{
		"app": conf.APIID,
		"ids": ids,
	}
	jsonParams, _ := json.Marshal(params)
	fmt.Printf("try delete params:%s\n", jsonParams)
	res := Request(urlStr, "DELETE", jsonParams)
	fmt.Printf("delete result:%s\n", res)

	return res
}

type Cursor struct {
	Id    string `json:"id"`
	Total int    `json:"totalCount"`
}

type RecordNumber struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Record struct {
	RecordNumber RecordNumber `json:"レコード番号"`
}

type Response struct {
	Records    []Record `json:"records"`
	TotalCount string   `json:"totalCount"`
}

func ExistedIndex(date string) []int {

	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records.json?app=%s&query=date=\"%s\"&fields[0]=レコード番号&totalCount=true", conf.Subdomain, conf.APIID, date)
	res := Request(urlStr, "GET", nil)
	var response Response
	json.Unmarshal(res, &response)
	var index []int
	for _, r := range response.Records {
		recordNumber, _ := strconv.Atoi(r.RecordNumber.Value)
		index = append(index, recordNumber)
	}
	fmt.Printf("existed index:%v\n", index)
	return index
}

func GetRecordByDate(fromDate string, toDate string) []byte {
	conf := readConf()
	urlStr := fmt.Sprintf("https://%s.cybozu.com/k/v1/records/cursor.json", conf.Subdomain)
	params := map[string]interface{}{
		"app":    conf.APIID,
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
