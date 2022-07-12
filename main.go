package main

import (
	"os"

	"github.com/AllenDang/giu"
	"github.com/isther/binance/binance"
	"github.com/isther/binance/conf"
	"github.com/isther/binance/global"
)

var (
	windowX int = 960
	windowY int = 1800

	symbol = global.Symbol

	testDoneCh = make(chan struct{})
)

func init() {
	os.Setenv("http_proxy", conf.Conf.Proxy)
	os.Setenv("https_proxy", conf.Conf.Proxy)

	startTipWindow()
	go ticker()
	go binance.StartWsPartialDepthServer()
}

func main() {
	// binance.Order()
	// binance.Status()
	// binance.SystemTime()
	app := giu.NewMasterWindow("Binance-GUI", windowY, windowX, 0)
	app.Run(mainWindow)
}
