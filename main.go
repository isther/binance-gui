package main

import (
	"os"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/conf"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

var (
	windowX int = 960
	windowY int = 1800
)

func init() {
	os.Setenv("http_proxy", conf.Conf.Proxy)
	os.Setenv("https_proxy", conf.Conf.Proxy)

	// global giu refresh
	go giuUpdateTicker()

	// network test
	startTipWindow()

	// console
	go console.ConsoleInstance.Start()

	plot()

	// update account
	binance.StartUpdateAccount()

	// start ws for partial depth server
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
	).RegisterKeyboardShortcuts(
		giu.WindowShortcut{Key: giu.KeyEnter, Modifier: giu.ModNone, Callback: func() { global.HotKeyRun = !global.HotKeyRun }},
		giu.WindowShortcut{Key: giu.KeyTab, Modifier: giu.ModNone, Callback: func() { global.ReverseTradeMode() }},
		giu.WindowShortcut{Key: giu.KeyMinus, Modifier: giu.ModNone, Callback: func() {
			if global.Average > 1 {
				global.Average--
			}
		}},
		giu.WindowShortcut{Key: giu.KeyEqual, Modifier: giu.ModNone, Callback: func() { global.Average++ }},
	)
	app.Run(mainWindow)
}
