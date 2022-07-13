package main

import (
	"fmt"
	"strings"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

var (
	symbol  = global.Symbol
	symbol1 = global.Symbol1
	symbol2 = global.Symbol2
)

func tipWindow() {
	giu.SingleWindow().Layout(
		giu.Label("Network Testing..."),
	)
}

func mainWindow() {
	giu.SingleWindowWithMenuBar().
		Layout(
			giu.MenuBar().Layout(
				giu.Menu("设置").Layout(
					giu.MenuItem("Api").OnClick(func() { giu.Msgbox("Info", "构建中...") }),
					giu.MenuItem("本地代理").OnClick(func() { giu.Msgbox("Info", "构建中...") }),
				),
			),
			giu.PrepareMsgbox(),
			giu.SplitLayout(giu.DirectionHorizontal, 1100, //H
				// 	giu.Label("筛选预警"),
				giu.SplitLayout(giu.DirectionVertical, 700,
					giu.SplitLayout(giu.DirectionHorizontal, 300, //V
						giu.TabBar().TabItems(
							giu.TabItem("市场行情").Layout(),
							giu.TabItem("持仓明细").Layout(),
						),
						giu.Label("K线"),
					),
					giu.Label(console.ConsoleInstance.Read()),
				),

				giu.SplitLayout(giu.DirectionHorizontal, 200, //H
					giu.SplitLayout(giu.DirectionVertical, 200, //V
						giu.Column(
							giu.Label(global.Symbol1+": "),
							giu.Label(fmt.Sprintf("  Free: %s", global.Symbol1Free)),
							giu.Label(fmt.Sprintf("  Locked: %s", global.Symbol1Locked)),
							giu.Label(global.Symbol2+": "),
							giu.Label(fmt.Sprintf("  Free: %s", global.Symbol2Free)),
							giu.Label(fmt.Sprintf("  Locked: %s", global.Symbol2Locked)),
						),
						giu.Label("成交明细区"),
					),
					giu.SplitLayout(giu.DirectionHorizontal, 270, //H
						giu.Column(
							giu.Row(
								giu.InputText(&symbol1).Size(100),
								giu.Label("/"),
								giu.InputText(&symbol2).Size(100),
								giu.Button("确定").OnClick(func() {
									symbol1 = strings.ToUpper(symbol1)
									symbol2 = strings.ToUpper(symbol2)

									symbolNew1 := symbol1 + symbol2
									symbolNew2 := symbol2 + symbol1
									if symbolNew1 == global.Symbol || symbolNew2 == global.Symbol {
										return
									}

									if binance.SymbolExist(symbolNew1) {
										symbol = symbolNew1
										global.FreshC <- symbol
										global.Symbol1 = symbol1
										global.Symbol2 = symbol2
									} else if binance.SymbolExist(symbolNew2) {
										symbol = symbolNew2
										global.FreshC <- symbol
										global.Symbol1 = symbol2
										global.Symbol2 = symbol1
									} else {
										giu.Msgbox("Error", "不存在的交易对")
									}

								}),
							),
							giu.Table().Freeze(0, 1).FastMode(true).Size(250, 840).Rows(binance.GetWsPartialDepthTable()...),
						),
						giu.Label("订单簿1"),
					),
				),
				// ),
			),
		)
}
