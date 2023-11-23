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

func HandleRequest(ctx context.Context) (*string, error) {

	var params map[string]interface{}
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
	params["deepSleepBR"] = breathData.BR[0].Value.DeepSleepSummary.BreathingRate
	params["remSleepBR"] = breathData.BR[0].Value.RemSleepSummary.BreathingRate
	params["fullSleepBR"] = breathData.BR[0].Value.FullSleepSummary.BreathingRate
	params["lightSleepBR"] = breathData.BR[0].Value.LightSleepSummary.BreathingRate

	

	// existed がサイズ0の場合
	if len(existed) == 0 {
		kintone.PostRecord(params)
	} else if len(existed) == 1 {
		kintone.UpdateRecord(existed[0], params)
	} else {
		return nil, fmt.Errorf("Today's data is duplicated")
	}
	return nil, nil
}
