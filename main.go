package main

import (
	"github.com/iorn121/MyDailyRecord/kintone"
)

func main() {
	// now := time.Now()
	// for i := 0; i < 1; i++ {
	// 	date := now.AddDate(0, 0, -i).Format("2006-01-02")
	// 	res := fitbit.HeartBeat(date)
	// 	fmt.Println(res)
	// }
	params := kintone.Params{
		App: 1,
		Record: kintone.Record{
			Hello: kintone.InnerRecord{
				Value: "ABC",
			},
		},
	}

	kintone.PostRecord(params)
}

// curl -X POST 'https://wbte0gl8dqbb.cybozu.com/k/v1/record.json' \
//   -H 'X-Cybozu-API-Token: bb1FXkRqkLTIi64qxUSLHi2Krudva5hacs3n49yD' \
//   -H 'Content-Type: application/json' \
//   -d '{
//     "app": 1,
//     "record": {
//       "hello": {
//         "value": "ABC"
//       }
//     }
//   }'
