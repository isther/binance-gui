package main

import (
	"fmt"

	// "github.com/isther/binance/gui"

	"github.com/AllenDang/giu"
	"github.com/isther/binance/global"
	"github.com/isther/binance/gui"
)

var (
	freshC = make(chan struct{})
)

func init() {
	go wsPartialDepthEvent()
}

func loop() {
	giu.SingleWindowWithMenuBar().RegisterKeyboardShortcuts(
		giu.WindowShortcut{Key: giu.Key0, Modifier: giu.ModNone, Callback: func() {}},
	).
		Layout(
			giu.MenuBar().Layout(
				giu.Menu("File").Layout(
					giu.MenuItem("Open"),
					giu.MenuItem("Save"),
				),
			),
			giu.SplitLayout(giu.DirectionHorizontal, 200, //H
				giu.Label("parameterSettings"),
				giu.SplitLayout(giu.DirectionHorizontal, 200, //H
					giu.SplitLayout(giu.DirectionVertical, 200, //V
						giu.Label("earlyWarning"),

						giu.Row(
							giu.InputText(&global.Symbol),
							giu.Button("Submit").OnClick(func() {
								freshC <- struct{}{}
							}),
						),
					),
					giu.SplitLayout(giu.DirectionHorizontal, 400, //H
						giu.SplitLayout(giu.DirectionVertical, 500, //V
							giu.Label("kLineChart"),
							giu.Label("transactionRecord"),
						),
						giu.SplitLayout(giu.DirectionHorizontal, 500, //H
							giu.SplitLayout(giu.DirectionVertical, 200, //V
								giu.Label("accountBalance"),
								giu.Label("transactionDetails"),
							),
							giu.SplitLayout(giu.DirectionHorizontal, 320, //H
								giu.Row(
									giu.Table().Freeze(0, 1).FastMode(true).Size(300, 800).Rows(gui.BuildWsPartialDepthServeRows()...),
								),
								giu.Label("order1"),
							),
						),
					),
				),
			),
		)
}

func main() {
	fmt.Println("Start Binance-GUI...")

	wnd := giu.NewMasterWindow("Binance-GUI", 1600, 960, 0)
	wnd.Run(loop)
}

func wsPartialDepthEvent() {
	var stopC chan struct{}

	stopC = gui.NewWsPartialDepthServer(global.Symbol, 20).StartWsPartialDepth(true)
	for {
		select {
		case <-freshC:
			stopC <- struct{}{}
			stopC = gui.NewWsPartialDepthServer(global.Symbol, 20).StartWsPartialDepth(true)
			giu.Update()
		}
	}
}
