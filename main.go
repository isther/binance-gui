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
	windowX int = 1080
	windowY int = 1920
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

	// start
	go binance.StartWebSocketStream()
}

func main() {
	if conf.Conf.Pprof {
		pprof()
	}

	app := giu.NewMasterWindow("Binance-GUI", windowX, windowX, giu.MasterWindowFlagsMaximized).
		RegisterKeyboardShortcuts( // 分仓下单
			regAllUsedKey(giu.ModNone)...,
		).RegisterKeyboardShortcuts( // 分仓数设置快捷键
		giu.WindowShortcut{Key: giu.KeyMinus, Modifier: giu.ModNone, Callback: func() {
			if global.Average > 1 {
				global.Average--
			}
		}},
		giu.WindowShortcut{Key: giu.KeyEqual, Modifier: giu.ModNone, Callback: func() { global.Average++ }},
		giu.WindowShortcut{Key: giu.KeyEnter, Modifier: giu.ModNone, Callback: func() { binance.UpdateAverageAmount() }},
	).RegisterKeyboardShortcuts( //全仓下单
		regAllUsedKey(giu.ModAlt)...,
	).RegisterKeyboardShortcuts( //撤单
		regAllUsedKey(giu.ModShift)...,
	).RegisterKeyboardShortcuts( // 打开关闭热键
		giu.WindowShortcut{Key: giu.KeySpace, Modifier: giu.ModNone, Callback: func() { global.ReverseHotKeyStatus() }},
	).RegisterKeyboardShortcuts( // 切换模式一二
		giu.WindowShortcut{Key: giu.KeyTab, Modifier: giu.ModNone, Callback: func() { global.ReverseTradeMode() }},
	).RegisterKeyboardShortcuts( // 全局买入
		giu.WindowShortcut{Key: giu.KeyF1, Modifier: giu.ModNone, Callback: func() {}}, //分仓买
		giu.WindowShortcut{Key: giu.KeyF2, Modifier: giu.ModNone, Callback: func() {}}, //全仓买
		giu.WindowShortcut{Key: giu.KeyF5, Modifier: giu.ModNone, Callback: func() {}}, //分仓卖
		giu.WindowShortcut{Key: giu.KeyF6, Modifier: giu.ModNone, Callback: func() {}}, //全仓卖
	).RegisterKeyboardShortcuts(
		giu.WindowShortcut{Key: giu.KeyF4, Modifier: giu.ModNone, Callback: func() {}},  // 撤销所有买单
		giu.WindowShortcut{Key: giu.KeyF8, Modifier: giu.ModNone, Callback: func() {}},  // 撤销所有卖单
		giu.WindowShortcut{Key: giu.KeyF9, Modifier: giu.ModNone, Callback: func() {}},  // 撤销所有单
		giu.WindowShortcut{Key: giu.KeyF12, Modifier: giu.ModNone, Callback: func() {}}, // 撤销所有单后市价卖出
	)
	app.Run(mainWindow)
}

const hotKeyTip = `
空格: 打开或关闭交易快捷键，默认关闭
Tab: 切换交易模式
-: 减少分仓数
=: 增加分仓数
Enter: 确认分仓数

全局下单:
	F1: 当前市价*波动比 分仓买入
	F2: 当前市价*波动比 全仓买入
	F5: 当前市价*波动比 分仓卖出
	F6: 当前市价*波动比 全仓卖出

全局撤单:
	F4: 撤当前交易对所有买单
	F8: 撤当前交易对所有卖单
	F9: 撤当前交易对所有单
	F12: 撤当前交易对所有单, 当前持有全部按市价卖出

默认分仓下单，请确保合理的分仓份数
快捷卖买: (请对应订单簿2)
	买：
	1 2 3 4 5 6 7 8 9 0
	q w e r t y u i o p

	卖：
	a s d f g h j k l ;
	z x c v b n m , . /

组合键:
	Alt + 快捷卖买键: 全仓下单
	Shift + 快捷买卖键: 撤销订单

`
