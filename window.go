package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

var (
	symbol  = binance.AccountInstance.Symbol
	symbol1 = binance.AccountInstance.One.Asset
	symbol2 = binance.AccountInstance.Two.Asset

	timeDataMin  float64
	timeDataMax  float64
	timeDataX    []float64
	timeDataY    []float64
	timeScatterY []float64
)

func tipWindow() {
	giu.SingleWindow().Layout(
		giu.PrepareMsgbox(),
		giu.Style().
			SetFontSize(60).To(
			giu.Label("Network Testing..."),
		),
	)
}

func mainWindow() {
	giu.SingleWindowWithMenuBar().
		Layout(
			giu.PrepareMsgbox(),
			giu.MenuBar().Layout(
				giu.Menu("设置").Layout(
					giu.MenuItem("Api").OnClick(func() { giu.Msgbox("Info", "构建中...") }),
					giu.MenuItem("本地代理").OnClick(func() { giu.Msgbox("Info", "构建中...") }),
				),
			),
			giu.SplitLayout(giu.DirectionHorizontal, 1200, //H
				giu.SplitLayout(giu.DirectionHorizontal, 600, //H
					giu.SplitLayout(giu.DirectionVertical, 600, //V
						giu.TabBar().TabItems(
							giu.TabItem("K线").Layout(
								giu.Plot("Plot Time Axe 时间线").Size(580, 540).AxisLimits(timeDataMin, timeDataMax, 0, 1, giu.ConditionOnce).XAxeFlags(giu.PlotAxisFlagsTime).Plots(
									giu.PlotLineXY("Time Line 时间线", timeDataX, timeDataY),
									giu.PlotScatterXY("Time Scatter 时间散点图", timeDataX, timeScatterY),
								),
							),
							giu.TabItem("持仓明细").Layout(),
						),
						giu.TabBar().TabItems(
							giu.TabItem("输出").Layout(
								giu.Label(console.ConsoleInstance.Read()),
							),
						),
					),
					giu.SplitLayout(giu.DirectionVertical, 600, //V
						giu.SplitLayout(giu.DirectionHorizontal, 300, //V
							giu.TabBar().TabItems(giu.TabItem("当前挂单").Layout()),
							giu.TabBar().TabItems(
								giu.TabItem("成交明细").Layout(),
							),
						),
						giu.SplitLayout(giu.DirectionVertical, 125, //V
							giu.TabBar().TabItems(
								giu.TabItem("交易账户余额").Layout(
									giu.Table().Columns(
										giu.TableColumn("Symbol"),
										giu.TableColumn("Free"),
										giu.TableColumn("Locked"),
									).Rows(
										giu.TableRow(
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.One.Asset)),
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.One.Free)),
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.One.Locked)),
										),
										giu.TableRow(
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.Two.Asset)),
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.Two.Free)),
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.Two.Locked)),
										),
										giu.TableRow(
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.BNB.Asset)),
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.BNB.Free)),
											giu.Label(fmt.Sprintf("%s", binance.AccountInstance.BNB.Locked)),
										),
									),
								),
							),
							giu.TabBar().TabItems(giu.TabItem("下单配置").Layout(
								giu.Column(
									giu.Row(
										giu.Label("交易对: "),
										giu.InputText(&symbol1).Size(100),
										giu.Label("/"),
										giu.InputText(&symbol2).Size(100),
										giu.Button("确定").OnClick(freshSymbol),
									),
									giu.Row(
										giu.Label("分仓数: "),
										giu.InputInt(&global.Average).Size(240),
									),
									giu.Style().
										SetColor(giu.StyleColorBorder, color.RGBA{0x36, 0x74, 0xD5, 255}).
										SetStyle(giu.StyleVarFramePadding, 10, 10).
										To(
											giu.Row(
												giu.Button(fmt.Sprintf("分仓金额: %s: %.8f / %s: %.8f",
													binance.AccountInstance.One.Asset, global.AverageSymbol1Amount,
													binance.AccountInstance.Two.Asset, global.AverageSymbol2Amount)),
											),
											giu.Row(
												giu.Button("交易热键状态(Enter): "+global.GetHotKeyStatus()).OnClick(func() { global.HotKeyRun = !global.HotKeyRun }),
												giu.Button("交易模式(Tab): "+global.GetTradeMode()).OnClick(func() { global.ReverseTradeMode() }),
											),
										),
								),
							)),
						),
					),
				),
				giu.SplitLayout(giu.DirectionHorizontal, 280, //H
					giu.Column(
						giu.TabBar().TabItems(
							giu.TabItem("订单簿2").Layout(
								giu.Table().FastMode(true).Size(270, 840).Rows(binance.GetWsPartialDepthTable()...),
							),
						),
					),
					giu.TabBar().TabItems(giu.TabItem("订单簿1").Layout()),
				),
			),
		)
}

func plot() {
	for i := 0; i < 100; i++ {
		timeDataX = append(timeDataX, float64(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Hour*time.Duration(24*i)).Unix()))
		timeDataY = append(timeDataY, rand.Float64())
		timeScatterY = append(timeScatterY, rand.Float64())
	}

	timeDataMin = timeDataX[0]
	timeDataMax = timeDataX[len(timeDataX)-1]
}
