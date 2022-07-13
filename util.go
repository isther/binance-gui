package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binance/binance"

	_ "net/http/pprof"
)

func startTipWindow() {
	var (
		connected  = false
		pingDoneCh = make(chan struct{})
	)

	startWindow := giu.NewMasterWindow("Network Testing...", 500, 100, 0)

	go func() {
		<-pingDoneCh
		startWindow.Close()
		if !connected {
			os.Exit(0)
		}
	}()

	go func() {
		timeOutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := binance.GetClient().NewPingService().Do(timeOutContext)
		if err != nil {
			connected = false
			panic(err)
		}

		connected = true
		pingDoneCh <- struct{}{}
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
func pprof() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
