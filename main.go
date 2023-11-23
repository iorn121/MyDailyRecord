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

	today := time.Now().Format("2006-01-02")
	existed := kintone.ExistedIndex(today)
	var params map[string]interface{}

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
