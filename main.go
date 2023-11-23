package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/iorn121/MyDailyRecord/kintone"
)

func main() {
	lambda.Start(HandleRequest)
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
	today := time.Now().Format("2006-01-02")

	
	if kintone.IsExisted(today) {
		fmt.Println("existed")
	} else {
		fmt.Println("not existed")
	}
	message := fmt.Sprintf("Hello %s!", event.Today)
	return &message, nil
}
