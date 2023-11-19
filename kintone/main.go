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

func getRequest(urlStr string) []byte {
	header := getHeader()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlStr, nil)
	for key, value := range header {
		req.Header.Add(key, value)
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
	url := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json?app=%s&id=%d", conf.Subdomain, conf.API_ID, id)
	res := getRequest(url)
	return res
}

func GetAllRecords() []byte {
	conf := readConf()
	url := fmt.Sprintf("https://%s.cybozu.com/k/v1/records.json?app=%s", conf.Subdomain, conf.API_ID)
	res := getRequest(url)
	return res
}

func postHeader() map[string]string {
	conf := readConf()
	header := map[string]string{
		"X-Cybozu-API-Token": conf.API_TOKEN,
		"Content-Type":       "application/json",
	}
	return header
}

func postRequest(url string, params []byte) []byte {
	header := postHeader()
	client := &http.Client{}
	jsonParams, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", url, nil)
	for key, value := range header {
		req.Header.Add(key, value)
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonParams))
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

type InnerRecord struct {
	Value string `json:"value"`
}

type OuterRecord struct {
	Hello InnerRecord `json:"hello"`
}

type Record struct {
	App    string      `json:"app"`
	Record OuterRecord `json:"record"`
}

func PostRecord(params Record) []byte {
	conf := readConf()
	url := fmt.Sprintf("https://%s.cybozu.com/k/v1/record.json?app=%s", conf.Subdomain, conf.API_ID)

	jsonParams, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonParams))
	res := postRequest(url, jsonParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
	return res
}
