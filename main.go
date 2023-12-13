package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iorn121/MyDailyRecord/fitbit"
	"github.com/iorn121/MyDailyRecord/kintone"
)

func convertToJP(timestr string) (time.Time, error) {

	// フォーマットを指定します。
	timeFormat := "2006-01-02T15:04:05"
	t, err := time.Parse(timeFormat, timestr)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation("Japan")
	if err != nil {
		return time.Time{}, err
	}
	jstTime := t.In(location)
	adjustTime := jstTime.Add(time.Duration(-9) * time.Hour)

	return adjustTime, nil
}

func main() {
	lambda.Start(HandleRequest)
	// HandleRequest(context.Background(), &MyEvent{Name: "test"}, today)
}

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event *MyEvent) error {

	var params map[string]interface{} = make(map[string]interface{})
	now := time.Now()
	fmt.Printf("start %s excecution at %s\n", event.Name, now)
	today := now.Format("2006-01-02")
	params["date"] = today
	existed := kintone.ExistedIndex(today)

	heartData := fitbit.HeartBeat(today)
	zones := []string{"outOfRange", "fatBurn", "cardio", "peak"}
	if len(heartData.ActivitiesHeart) > 0 {
		for i, hrz := range heartData.ActivitiesHeart[0].Value.HeartRateZones {
			zone := zones[i]
			params[zone+"CaloriesOut"] = hrz.CaloriesOut
			params[zone+"Min"] = hrz.Min
			params[zone+"Max"] = hrz.Max
			params[zone+"Minutes"] = hrz.Minutes
		}
		params["restingHeartRate"] = heartData.ActivitiesHeart[0].Value.RestingHeartRate
	}

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
	if len(sleepData.Sleep) > 0 {
		params["sleepDuration"] = sleepData.Sleep[0].Duration
		params["sleepEfficiency"] = sleepData.Sleep[0].Efficiency
		startTime := sleepData.Sleep[0].StartTime
		endTime := sleepData.Sleep[0].EndTime
		startTimeJP, err := convertToJP(startTime)
		if err != nil {
			fmt.Println("Error converting start time to UTC:", err)
		}
		endTimeJP, err := convertToJP(endTime)
		if err != nil {
			fmt.Println("Error converting end time to UTC:", err)
		}
		params["sleepStartTime"] = startTimeJP.Format("2006-01-02T15:04:05.000-07:00")
		params["sleepEndTime"] = endTimeJP.Format("2006-01-02T15:04:05.000-07:00")
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
	if len(HRVData.Hrv) > 0 {
		params["dailyHRV"] = HRVData.Hrv[0].Value.DailyRmssd
		params["deepHRV"] = HRVData.Hrv[0].Value.DeepRmssd
	}

	if len(existed) == 0 {
		fmt.Println("There are no data for today.")
		kintone.PostRecord(params)
	} else if len(existed) == 1 {
		fmt.Printf("There are already data for today at %d.", existed[0])
		kintone.UpdateRecord(existed[0], params)
	} else {
		fmt.Println("There is an error searching today's data")
	}
	fmt.Printf("\nend %s execution\n", event.Name)
	return nil
}
