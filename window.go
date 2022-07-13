package main

import (
	"strings"

	"github.com/AllenDang/giu"
	"github.com/isther/binance/binance"
	"github.com/isther/binance/global"
)

func tipWindow() {
	giu.SingleWindow().Layout(
		giu.Label("Network Testing..."),
	)
}

func mainWindow() {
	giu.SingleWindowWithMenuBar().RegisterKeyboardShortcuts(
		giu.WindowShortcut{Key: giu.Key0, Modifier: giu.ModNone, Callback: func() {}},
	).
		Layout(
			giu.MenuBar().Layout(
				giu.Menu("Setting").Layout(
					giu.MenuItem("Api").OnClick(func() {}),
					giu.MenuItem("Proxy").OnClick(func() {}),
				),
			),
			giu.PrepareMsgbox(),
			giu.SplitLayout(giu.DirectionHorizontal, 200, //H
				giu.Label("parameterSettings"),
				giu.SplitLayout(giu.DirectionHorizontal, 200, //H
					giu.SplitLayout(giu.DirectionVertical, 200, //V
						giu.Label("earlyWarning"),

						giu.Row(
							giu.InputText(&symbol),
							giu.Button("Submit").OnClick(func() {
								symbol = strings.ToUpper(symbol)
								if global.Symbol == symbol {
									return
								}
								if binance.SymbolExist(symbol) {
									global.FreshC <- symbol
									giu.Msgbox("Info", "Ok")
								} else {
									giu.Msgbox("Error", "No trading pair")
								}
							}),
						),
					),
					giu.SplitLayout(giu.DirectionHorizontal, 500, //H
						giu.SplitLayout(giu.DirectionVertical, 500, //V
							giu.Label("kLineChart"),
							giu.Label("transactionRecord"),
						),
						giu.SplitLayout(giu.DirectionHorizontal, 200, //H
							giu.SplitLayout(giu.DirectionVertical, 200, //V
								giu.Label("accountBalance"),
								giu.Label("transactionDetails"),
							),
							giu.SplitLayout(giu.DirectionHorizontal, 350, //H
								giu.Row(
									giu.Table().Freeze(0, 1).FastMode(true).Size(350, 960).Rows(binance.GetWsPartialDepthTable()...),
								),
								giu.Label("order1"),
							),
						),
					),
				),
			),
		)
}
