package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iorn121/MyDailyRecord/fitbit"
	"github.com/iorn121/MyDailyRecord/kintone"
)

func main() {
	lambda.Start(HandleRequest)
}

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	fmt.Printf("start %s excecution", event.Name)

	var params map[string]interface{} = make(map[string]interface{})
	today := time.Now().Format("2006-01-02")
	params["date"] = today
	existed := kintone.ExistedIndex(today)

	heartData := fitbit.HeartBeat(today)
	zones := []string{"outOfRange", "fatBurn", "cardio", "peak"}
	for i, hrz := range heartData.ActivitiesHeart[0].Value.HeartRateZones {
		zone := zones[i]
		params[zone+"CaloriesOut"] = hrz.CaloriesOut
		params[zone+"Min"] = hrz.Min
		params[zone+"Max"] = hrz.Max
		params[zone+"Minutes"] = hrz.Minutes
	}
	params["restingHeartRate"] = heartData.ActivitiesHeart[0].Value.RestingHeartRate

	breathData := fitbit.BreathingRate(today)
	if len(breathData.BR) > 0 {
		params["deepSleepBR"] = breathData.BR[0].Value.DeepSleepSummary.BreathingRate
		params["remSleepBR"] = breathData.BR[0].Value.RemSleepSummary.BreathingRate
		params["fullSleepBR"] = breathData.BR[0].Value.FullSleepSummary.BreathingRate
		params["lightSleepBR"] = breathData.BR[0].Value.LightSleepSummary.BreathingRate
	}

	tempData := fitbit.SkinTemperature(today)
	if len(tempData.TempSkin) > 0 {
		params["skinTemp"] = tempData.TempSkin[0].Value.NightlyRelative
	}

	sleepData := fitbit.SleepDetail(today)
	if len(sleepData.Sleep) == 0 {
		params["sleepDuration"] = sleepData.Sleep[0].Duration
		params["sleepEfficiency"] = sleepData.Sleep[0].Efficiency
		params["sleepStartTime"] = sleepData.Sleep[0].StartTime
		params["sleepEndTime"] = sleepData.Sleep[0].EndTime
		params["timeInBed"] = sleepData.Sleep[0].TimeInBed
		sleepLevels := []string{"wake", "light", "rem", "deep"}
		for _, sl := range sleepLevels {
			params[sl+"SleepCount"] = sleepData.Sleep[0].Levels.Summary[sl].Count
			params[sl+"SleepMinutes"] = sleepData.Sleep[0].Levels.Summary[sl].Minutes
			params[sl+"SleepThirtyDayAvgMinutes"] = sleepData.Sleep[0].Levels.Summary[sl].ThirtyDayAvgMinutes
		}
	}

	spO2Data := fitbit.SpO2(today)
	params["SpO2Min"] = spO2Data.Value.Min
	params["SpO2Max"] = spO2Data.Value.Max
	params["SpO2Avg"] = spO2Data.Value.Avg

	VO2Data := fitbit.VO2Max(today)
	if len(VO2Data.CardioScore) > 0 {
		params["VO2Max"] = VO2Data.CardioScore[0].Value.Vo2Max
	}

	HRVData := fitbit.HRV(today)
	if len(HRVData.Hrv) == 0 {
		params["dailyHRV"] = HRVData.Hrv[0].Value.DailyRmssd
		params["deepHRV"] = HRVData.Hrv[0].Value.DeepRmssd
	}

	// existed がサイズ0の場合
	if len(existed) == 0 {
		kintone.PostRecord(params)
	} else if len(existed) == 1 {
		kintone.UpdateRecord(existed[0], params)
	} else {
		return nil, fmt.Errorf("Today's data is duplicated")
	}
	fmt.Printf("end %s excecution", event.Name)
	return nil, nil
}
