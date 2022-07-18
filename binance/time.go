package binance

import (
	"context"
	"fmt"
	"time"
)

var (
	TimeString = "00:00:00"
)

func updateTime() {
	go func() {
		fmt.Println("New http request to get system time of binance")
		var (
			ticker = time.NewTicker(1 * time.Second)
		)

		var timestamp = getSystemTime()

		for i := 0; i < 60; i++ {
			now := time.Unix(timestamp/1000, timestamp%1000)
			TimeString = now.Format("2006-01-02 15:04:05")
			<-ticker.C

			timestamp += 1000
		}

		updateTime()
	}()
}

func getSystemTime() int64 {
	timestamp, err := GetClient().NewServerTimeService().Do(context.Background())
	if err != nil {
		panic(err)
	}
	return timestamp
}
