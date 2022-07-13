package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"

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

func regAllUsedKey(mod giu.Modifier) []giu.WindowShortcut {
	return []giu.WindowShortcut{
		// Sale
		{Key: giu.KeyA, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "A").Trade() }},
		{Key: giu.KeyS, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "S").Trade() }},
		{Key: giu.KeyD, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "D").Trade() }},
		{Key: giu.KeyF, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "F").Trade() }},
		{Key: giu.KeyG, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "G").Trade() }},
		{Key: giu.KeyH, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "H").Trade() }},
		{Key: giu.KeyJ, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "J").Trade() }},
		{Key: giu.KeyK, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "K").Trade() }},
		{Key: giu.KeyL, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "L").Trade() }},
		{Key: giu.KeySemicolon, Modifier: mod, Callback: func() { go binance.NewTrader(mod, ";").Trade() }},
		{Key: giu.KeyZ, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "Z").Trade() }},
		{Key: giu.KeyX, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "X").Trade() }},
		{Key: giu.KeyC, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "C").Trade() }},
		{Key: giu.KeyV, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "V").Trade() }},
		{Key: giu.KeyB, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "B").Trade() }},
		{Key: giu.KeyN, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "N").Trade() }},
		{Key: giu.KeyM, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "M").Trade() }},
		{Key: giu.KeyComma, Modifier: mod, Callback: func() { go binance.NewTrader(mod, ",").Trade() }},
		{Key: giu.KeyPeriod, Modifier: mod, Callback: func() { go binance.NewTrader(mod, ".").Trade() }},
		{Key: giu.KeySlash, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "/").Trade() }},

		// Buy
		{Key: giu.Key1, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "1").Trade() }},
		{Key: giu.Key2, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "2").Trade() }},
		{Key: giu.Key3, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "3").Trade() }},
		{Key: giu.Key4, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "4").Trade() }},
		{Key: giu.Key5, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "5").Trade() }},
		{Key: giu.Key6, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "6").Trade() }},
		{Key: giu.Key7, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "7").Trade() }},
		{Key: giu.Key8, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "8").Trade() }},
		{Key: giu.Key9, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "9").Trade() }},
		{Key: giu.Key0, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "0").Trade() }},
		{Key: giu.KeyQ, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "Q").Trade() }},
		{Key: giu.KeyW, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "W").Trade() }},
		{Key: giu.KeyE, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "E").Trade() }},
		{Key: giu.KeyR, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "R").Trade() }},
		{Key: giu.KeyT, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "T").Trade() }},
		{Key: giu.KeyY, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "Y").Trade() }},
		{Key: giu.KeyU, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "U").Trade() }},
		{Key: giu.KeyI, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "I").Trade() }},
		{Key: giu.KeyO, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "O").Trade() }},
		{Key: giu.KeyP, Modifier: mod, Callback: func() { go binance.NewTrader(mod, "P").Trade() }},
	}
}


func pprof() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}
