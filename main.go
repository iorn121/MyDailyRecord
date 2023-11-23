package main

import (
	"fmt"

	"github.com/iorn121/MyDailyRecord/kintone"
)

func main() {
	res := kintone.IsExisted("2023-11-19")
	res2 := kintone.GetRecordByDate("2023-11-20", "2023-11-19")
	fmt.Println(res)
	fmt.Println(string(res2))
}
