package fitbit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Config is a struct for conf.json
type Config struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
}

// Error is a struct for error message
type Error struct {
	ErrorType string `json:"errorType"`
}

// Message is a struct for error message
type Message struct {
	Errors []Error `json:"errors"`
}

// NewToken is a struct for new token
type NewToken struct {
	AccessToken string `json:"access_token"`
}

// readConf reads conf.json and set to conf
func readConf() Config {
	file, _ := os.Open("conf.json")
	defer file.Close()
	var conf Config
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&conf)
	return conf
}

// bearerHeader returns header for bearer token
func bearerHeader() map[string]string {
	conf := readConf()
	return map[string]string{
		"Authorization": "Bearer " + conf.AccessToken,
	}
}

// refresh refreshes access token and write to conf.json
func refresh() {
	urlStr := "https://api.fitbit.com/oauth2/token"
	conf := readConf()
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

// isExpired checks if access token is expired
// token lifetime is 28800 seconds (8 hours)
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

// request sends request to urlStr and returns response body
// if access token is expired, refresh and retry
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

// HeartRateZone is a struct for heart rate zone
// Retrieves the heart rate time series data over a period of time by specifying a date and time period.
type HeartRateZone struct {
	CaloriesOut float64 `json:"caloriesOut"`
	Max         int     `json:"max"`
	Min         int     `json:"min"`
	Minutes     int     `json:"minutes"`
	Name        string  `json:"name"`
}

type HeartValue struct {
	CustomHeartRateZones []HeartRateZone `json:"customHeartRateZones"`
	HeartRateZones       []HeartRateZone `json:"heartRateZones"`
	RestingHeartRate     int             `json:"restingHeartRate"`
}

type ActivitiesHeart struct {
	DateTime string    `json:"dateTime"`
	Value    HeartData `json:"value"`
}

type everytime struct {
	Time  string `json:"time"`
	Value int    `json:"value"`
}

type Intraday struct {
	Dataset         []everytime `json:"dataset"`
	DatasetInterval int         `json:"datasetInterval"`
	DatasetType     string      `json:"datasetType"`
}

type HeartData struct {
	ActivitiesHeart []ActivitiesHeart `json:"activities-heart"`
	Intraday        Intraday          `json:"activities-heart-intraday"`
}

// heartBeat returns heart beat data
func heartBeat(date string) HeartData {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/activities/heart/date/%s/1d.json", date)
	res := request(url)
	var heartData HeartData
	err := json.Unmarshal(res, &heartData)
	if err != nil {
		fmt.Println(err)
	}
	return heartData
}

// BreathingRate measures the average breathing rate throughout the day and categories your breathing rate by sleep stage.
// The breathing rate is measured in breaths per minute (BPM).
type SleepSummary struct {
	BreathingRate float64 `json:"breathingRate"`
}

type SleepValue struct {
	DeepSleepSummary  SleepSummary `json:"deepSleepSummary"`
	RemSleepSummary   SleepSummary `json:"remSleepSummary"`
	FullSleepSummary  SleepSummary `json:"fullSleepSummary"`
	LightSleepSummary SleepSummary `json:"lightSleepSummary"`
}

type BR struct {
	Value    SleepValue `json:"value"`
	DateTime string     `json:"dateTime"`
}

type BreathingRateData struct {
	BR []BR `json:"br"`
}

// breathingRate returns breathing rate data for a given day
func breathingRate(date string) BreathingRateData {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/br/date/%s/all.json", date)
	res := request(url)
	var brData BreathingRateData
	err := json.Unmarshal(res, &brData)
	if err != nil {
		fmt.Println(err)
	}
	return brData
}

// SkinTemperature measures the average skin temperature conmpared to your personal baseline.
type TempValue struct {
	NightlyRelative float64 `json:"nightlyRelative"`
}

type TempSkin struct {
	DateTime string    `json:"dateTime"`
	Value    TempValue `json:"value"`
	LogType  string    `json:"logType"`
}

type TempSkinData struct {
	TempSkin []TempSkin `json:"tempSkin"`
}

// skinTemperature returns skin temperature data for a given day
func skinTemperature(date string) TempSkinData {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/temp/skin/date/%s.json", date)
	res := request(url)
	var tempSkinData TempSkinData
	err := json.Unmarshal(res, &tempSkinData)
	if err != nil {
		fmt.Println(err)
	}
	return tempSkinData
}

type LevelData struct {
	DateTime string `json:"dateTime"` // Timestamp the user started in sleep level.
	Level    string `json:"level"`    // The sleep level the user entered.
	Seconds  int    `json:"seconds"`  // The length of time the user was in the sleep level. Displayed in seconds.
}

type SummaryData struct {
	Count               int `json:"count"`               // Total number of times the user entered the sleep level.
	Minutes             int `json:"minutes"`             // Total number of minutes the user appeared in the sleep level.
	ThirtyDayAvgMinutes int `json:"thirtyDayAvgMinutes"` // The average sleep stage time over the past 30 days.
}

type Levels struct {
	Data      []LevelData            `json:"data"`      // Data about each sleep level.
	ShortData []LevelData            `json:"shortData"` // Short data about each sleep level.
	Summary   map[string]SummaryData `json:"summary"`   // Summary of each sleep level.
}

type Sleep struct {
	DateOfSleep         string `json:"dateOfSleep"`         // The date the sleep log ended.
	Duration            int    `json:"duration"`            // Length of the sleep in milliseconds.
	Efficiency          int    `json:"efficiency"`          // Calculated sleep efficiency score.
	EndTime             string `json:"endTime"`             // Time the sleep log ended.
	InfoCode            int    `json:"infoCode"`            // An integer value representing the quality of data collected within the sleep log.
	IsMainSleep         bool   `json:"isMainSleep"`         // Boolean value: true or false
	Levels              Levels `json:"levels"`              // Levels of sleep.
	LogId               int64  `json:"logId"`               // Sleep log ID.
	MinutesAfterWakeup  int    `json:"minutesAfterWakeup"`  // The total number of minutes after the user woke up.
	MinutesAsleep       int    `json:"minutesAsleep"`       // The total number of minutes the user was asleep.
	MinutesAwake        int    `json:"minutesAwake"`        // The total sum of "wake" minutes only.
	MinutesToFallAsleep int    `json:"minutesToFallAsleep"` // The total number of minutes before the user falls asleep.
	LogType             string `json:"logType"`             // The type of sleep in terms of how it was logged.
	StartTime           string `json:"startTime"`           // Time the sleep log begins.
	TimeInBed           int    `json:"timeInBed"`           // Total number of minutes the user was in bed.
	Type                string `json:"type"`                // The type of sleep log.
}

type Stages struct {
	Deep  int `json:"deep"`
	Light int `json:"light"`
	Rem   int `json:"rem"`
	Wake  int `json:"wake"`
}

type Summary struct {
	Stages             Stages `json:"stages"`
	TotalMinutesAsleep int    `json:"totalMinutesAsleep"`
	TotalSleepRecords  int    `json:"totalSleepRecords"`
	TotalTimeInBed     int    `json:"totalTimeInBed"`
}

type SleepData struct {
	Sleep   []Sleep `json:"sleep"`
	Summary Summary `json:"summary"`
}

// sleep returns sleep data for a given day
func sleep(date string) SleepData {
	url := fmt.Sprintf("https://api.fitbit.com/1.2/user/-/sleep/date/%s.json", date)
	res := request(url)
	var sleepData SleepData
	err := json.Unmarshal(res, &sleepData)
	if err != nil {
		fmt.Println(err)
	}
	return sleepData
}

// SpO2 (Oxygen Saturation) is an estimate of the amount of oxygen in your blood.
type SpO2Value struct {
	Avg float64 `json:"avg"`
	Max float64 `json:"max"`
	Min float64 `json:"min"`
}

type SpO2Data struct {
	DateTime string    `json:"dateTime"`
	Value    SpO2Value `json:"value"`
}

func SpO2(date string) SpO2Data {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/spo2/date/%s.json", date)
	res := request(url)
	var spo2Data SpO2Data
	err := json.Unmarshal(res, &spo2Data)
	if err != nil {
		fmt.Println(err)
	}
	return spo2Data
}

// VO2Max is a measure of the maximum volume of oxygen that an athlete can use.
type CardioScoreValue struct {
	Vo2Max string `json:"vo2Max"`
}

type CardioScore struct {
	DateTime string           `json:"dateTime"`
	Value    CardioScoreValue `json:"value"`
}

type CardioScoreData struct {
	CardioScore []CardioScore `json:"cardioScore"`
}

func VO2Max(date string) CardioScoreData {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/cardioscore/date/%s.json", date)
	res := request(url)
	var cardioScoreData CardioScoreData
	err := json.Unmarshal(res, &cardioScoreData)
	if err != nil {
		fmt.Println(err)
	}
	return cardioScoreData
}

// Heart Rate Variability (HRV) data applies specifically to a user’s “main sleep,” which is the longest single period of time asleep on a given date.
type HrvValue struct {
	DailyRmssd float64 `json:"dailyRmssd"`
	DeepRmssd  float64 `json:"deepRmssd"`
}

type Hrv struct {
	Value    HrvValue `json:"value"`
	DateTime string   `json:"dateTime"`
}

type HrvData struct {
	Hrv []Hrv `json:"hrv"`
}

func HRV(date string) HrvData {
	url := fmt.Sprintf("https://api.fitbit.com/1/user/-/hrv/date/%s.json", date)
	res := request(url)
	var hrvData HrvData
	err := json.Unmarshal(res, &hrvData)
	if err != nil {
		fmt.Println(err)
	}
	return hrvData
}
