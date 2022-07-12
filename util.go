package main

import (
	"context"
	"os"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binance/binance"
	"github.com/isther/binance/global"
)

func startTipWindow() {
	startWindow := giu.NewMasterWindow("Network Testing...", 500, 100, 0)

	go func() {
		<-testDoneCh
		startWindow.Close()
		if !global.Connected {
			os.Exit(0)
		}
	}()

	go func() {
		timeOutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := binance.GetClient().NewPingService().Do(timeOutContext)
		if err != nil {
			global.Connected = false
			panic(err)
		}

		global.Connected = true
		testDoneCh <- struct{}{}
	}()
	startWindow.Run(tipWindow)
}

func ticker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C
		giu.Update()
	}
}
