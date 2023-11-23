package main

import (
	"context"
	"fmt"
	"time"

	"github.com/iorn121/MyDailyRecord/kintone"
)

func main() {
	// lambda.Start(HandleRequest)
	flg := kintone.IsExisted("2024-01-01")
	fmt.Println(flg)
}

type MyEvent struct {
	Today string `json:"today"`
}

func isDate(date string) bool {
	_, err := time.Parse(time.RFC3339, date)
	return err != nil
}

func HandleRequest(ctx context.Context, event *MyEvent) (*string, error) {
	if event == nil {
		return nil, fmt.Errorf("received nil event")
	}
	if !isDate(event.Today) {
		return nil, fmt.Errorf("invalid date format: %s", event.Today)
	}
	if kintone.IsExisted("TODAY()") {
		fmt.Println("existed")
	} else {
		fmt.Println("not existed")
	}
	message := fmt.Sprintf("Hello %s!", event.Today)
	return &message, nil
}
