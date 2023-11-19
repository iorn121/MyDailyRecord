package main

import (
	"fmt"
	"time"
	"github.com/iorn121/MyDailyRecord/fitbit"
)

func main() {
	now := time.Now()
	for i := 0; i < 1; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		res := fitbit.heartBeat(date)
		fmt.Println(res)
	}
}