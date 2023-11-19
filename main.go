package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Config struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
}

type Error struct {
	ErrorType string `json:"errorType"`
}

type Message struct {
	Errors []Error `json:"errors"`
}

type NewToken struct {
	AccessToken string `json:"access_token"`
}

var conf Config

func main() {
	res := heartBeat("2023-11-19", "1d")
	fmt.Println("result", string(res))
}

func readConf() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&conf)
}

func bearerHeader() map[string]string {
	readConf()
	return map[string]string{
		"Authorization": "Bearer " + conf.AccessToken,
	}
}

func refresh() {
	urlStr := "https://api.fitbit.com/oauth2/token"
	readConf()
	params := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {conf.RefreshToken},
		"client_id":     {conf.ClientID},
	}
	res, err := http.PostForm(urlStr, params)

	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var newToken NewToken
	_ = json.Unmarshal(body, &newToken)
	conf.AccessToken = newToken.AccessToken
	file, _ := os.OpenFile("conf.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	defer file.Close()
	encoder := json.NewEncoder(file)
	_ = encoder.Encode(conf)
}

func isExpired(res []byte) bool {
	var msg Message
	_ = json.Unmarshal(res, &msg)

	if msg.Errors == nil {
		return false
	}

	for _, err := range msg.Errors {
		if err.ErrorType == "expired_token" {
			fmt.Println("TOKEN_EXPIRED!!!")
			return true
		}
	}

	return false
}

func request(urlStr string) []byte {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", urlStr, nil)
	headers := bearerHeader()
	for key, value := range headers {
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

	if isExpired(resBody) {
		refresh()
		resBody = request(urlStr)
	}

	return resBody
}

func heartBeat(date string, period string) []byte {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/%s.json", date, period)
	return request(url)
}
