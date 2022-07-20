package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/global"

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
		if !connected {
			giu.Msgbox("Error", "网络连接失败，请重试！").Buttons(giu.MsgboxButtonsOk)
			time.Sleep(5 * time.Second)
			os.Exit(-1)
		}
		startWindow.Close()
	}()

	go func() {
		timeOutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := binance.GetClient().NewPingService().Do(timeOutContext)
		if err != nil {
			connected = false
		} else {
			connected = true
		}

		pingDoneCh <- struct{}{}
	}()

	startWindow.Run(tipWindow)
}

func giuUpdateTicker() {
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		<-ticker.C
		giu.Update()
	}
}

func regAllUsedKey(mode giu.Modifier) []giu.WindowShortcut {
	return []giu.WindowShortcut{
		// Sale
		{Key: giu.KeyA, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "A") }},
		{Key: giu.KeyS, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "S") }},
		{Key: giu.KeyD, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "D") }},
		{Key: giu.KeyF, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "F") }},
		{Key: giu.KeyG, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "G") }},
		{Key: giu.KeyH, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "H") }},
		{Key: giu.KeyJ, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "J") }},
		{Key: giu.KeyK, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "K") }},
		{Key: giu.KeyL, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "L") }},
		{Key: giu.KeySemicolon, Modifier: mode, Callback: func() { hotKeyCallBack(mode, ";") }},
		{Key: giu.KeyZ, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "Z") }},
		{Key: giu.KeyX, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "X") }},
		{Key: giu.KeyC, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "C") }},
		{Key: giu.KeyV, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "V") }},
		{Key: giu.KeyB, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "B") }},
		{Key: giu.KeyN, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "N") }},
		{Key: giu.KeyM, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "M") }},
		{Key: giu.KeyComma, Modifier: mode, Callback: func() { hotKeyCallBack(mode, ",") }},
		{Key: giu.KeyPeriod, Modifier: mode, Callback: func() { hotKeyCallBack(mode, ".") }},
		{Key: giu.KeySlash, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "/") }},

		// Buy
		{Key: giu.Key1, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "1") }},
		{Key: giu.Key2, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "2") }},
		{Key: giu.Key3, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "3") }},
		{Key: giu.Key4, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "4") }},
		{Key: giu.Key5, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "5") }},
		{Key: giu.Key6, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "6") }},
		{Key: giu.Key7, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "7") }},
		{Key: giu.Key8, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "8") }},
		{Key: giu.Key9, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "9") }},
		{Key: giu.Key0, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "0") }},
		{Key: giu.KeyQ, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "Q") }},
		{Key: giu.KeyW, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "W") }},
		{Key: giu.KeyE, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "E") }},
		{Key: giu.KeyR, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "R") }},
		{Key: giu.KeyT, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "T") }},
		{Key: giu.KeyY, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "Y") }},
		{Key: giu.KeyU, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "U") }},
		{Key: giu.KeyI, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "I") }},
		{Key: giu.KeyO, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "O") }},
		{Key: giu.KeyP, Modifier: mode, Callback: func() { hotKeyCallBack(mode, "P") }},
	}
}

func hotKeyCallBack(mode giu.Modifier, key string) {
	if global.HotKeyRun {
		go binance.NewTrader(mode, key).Trade()
	}
}

func pprof() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

func ping() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		timeOutContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		start := time.Now()
		err := binance.GetClient().NewPingService().Do(timeOutContext)
		end := time.Now()
		if err != nil {
			global.Ping = "999ms"
		}
		global.Ping = fmt.Sprintf("%.2fms", float32(end.Sub(start).Microseconds())/1e3)
		<-ticker.C
	}
}

func freshSymbol() {
	symbol1 = strings.ToUpper(symbol1)
	symbol2 = strings.ToUpper(symbol2)

	symbolNew1 := symbol1 + symbol2
	symbolNew2 := symbol2 + symbol1
	if symbolNew1 == binance.AccountInstance.One.Asset || symbolNew2 == binance.AccountInstance.Two.Asset {
		return
	}

	if binance.SymbolExist(symbolNew1) {
		symbol = symbolNew1
		global.FreshC <- symbol
		binance.AccountInstance.One.Asset = symbol1
		binance.AccountInstance.Two.Asset = symbol2
	} else if binance.SymbolExist(symbolNew2) {
		symbol = symbolNew2
		global.FreshC <- symbol
		binance.AccountInstance.One.Asset = symbol2
		binance.AccountInstance.Two.Asset = symbol1
	} else {
		giu.Msgbox("Error", "不存在的交易对")
	}

}

func IsExpired() bool {
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		panic(err)
	}
	return endTime.Before(time.Now())
}
