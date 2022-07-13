package main

import (
	"os"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/conf"
	"github.com/isther/binanceGui/console"
)

var (
	windowX int = 960
	windowY int = 1800
)

func init() {
	os.Setenv("http_proxy", conf.Conf.Proxy)
	os.Setenv("https_proxy", conf.Conf.Proxy)

	go ticker()
	startTipWindow()

	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			binance.UpdateAccount()
			<-ticker.C
		}
	}()

	go console.ConsoleInstance.Start()
	go binance.StartWsPartialDepthServer()
}

func main() {
	if conf.Conf.Pprof {
		pprof()
	}

	app := giu.NewMasterWindow("Binance-GUI", windowY, windowX, 0).RegisterKeyboardShortcuts(
		regAllUsedKey(giu.ModNone)...,
	).RegisterKeyboardShortcuts(
		regAllUsedKey(giu.ModAlt)...,
	).RegisterKeyboardShortcuts(
		regAllUsedKey(giu.ModShift)...,
	)
	app.Run(mainWindow)
}
