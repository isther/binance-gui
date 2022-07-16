package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

var (
	volatilityRatiosInputSize float32 = 100

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
				giu.Menu("说明").Layout(
					giu.MenuItem("快捷键").OnClick(func() {
						giu.Msgbox("快捷键说明", hotKeyTip).Buttons(giu.MsgboxButtonsOk)
					}),
				),
			),
			giu.SplitLayout(giu.DirectionHorizontal, 1350, //H
				giu.SplitLayout(giu.DirectionVertical, 600, //V
					giu.SplitLayout(giu.DirectionHorizontal, 750, //V
						giu.TabBar().TabItems(
							giu.TabItem("K线").Layout(
								giu.Plot("Plot Time Axe 时间线").Size(580, 540).AxisLimits(timeDataMin, timeDataMax, 0, 1, giu.ConditionOnce).XAxeFlags(giu.PlotAxisFlagsTime).Plots(
									giu.PlotLineXY("Time Line 时间线", timeDataX, timeDataY),
									giu.PlotScatterXY("Time Scatter 时间散点图", timeDataX, timeScatterY),
								),
							),
							giu.TabItem("持仓明细").Layout(),
							giu.TabItem("波动比设置").Layout(),
							giu.TabItem("波动比设置").Layout(giu.Style().
								SetColor(giu.StyleColorBorder, global.BLUE).
								SetStyle(giu.StyleVarFramePadding, 10, 10).
								To(
									giu.Column(
										giu.Label("卖"),
										giu.Row(
											giu.Label("06 - 10"),
											giu.InputFloat(&global.VolatilityRatiosSale[5]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[6]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[7]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[8]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[9]).Size(volatilityRatiosInputSize),
										),
										giu.Row(
											giu.Label("11 - 15"),
											giu.InputFloat(&global.VolatilityRatiosSale[10]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[11]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[12]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[13]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[14]).Size(volatilityRatiosInputSize),
										),
										giu.Row(
											giu.Label("16 - 20"),
											giu.InputFloat(&global.VolatilityRatiosSale[15]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[16]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[17]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[18]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosSale[19]).Size(volatilityRatiosInputSize),
										),
										giu.Label("买"),
										giu.Row(
											giu.Label("06 - 10"),
											giu.InputFloat(&global.VolatilityRatiosBuy[5]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[6]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[7]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[8]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[9]).Size(volatilityRatiosInputSize),
										),
										giu.Row(
											giu.Label("11 - 15"),
											giu.InputFloat(&global.VolatilityRatiosBuy[10]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[11]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[12]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[13]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[14]).Size(volatilityRatiosInputSize),
										),
										giu.Row(
											giu.Label("16 - 20"),
											giu.InputFloat(&global.VolatilityRatiosBuy[15]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[16]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[17]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[18]).Size(volatilityRatiosInputSize),
											giu.InputFloat(&global.VolatilityRatiosBuy[19]).Size(volatilityRatiosInputSize),
										),
									),
								)),
						),
						giu.SplitLayout(giu.DirectionHorizontal, 300, //H
							giu.TabBar().TabItems(
								giu.TabItem("当前挂单").Layout(
									giu.SplitLayout(giu.DirectionVertical, 280,
										giu.Table().Freeze(0, 1).FastMode(true).Size(270, 300).Rows(binance.GetOpenSaleOrdersTable()...),
										giu.Table().Freeze(0, 1).FastMode(true).Size(270, 300).Rows(binance.GetOpenBuyOrdersTable()...),
									),
								),
							),
							giu.TabBar().TabItems(
								giu.TabItem("成交明细").Layout(
									giu.Table().Freeze(0, 1).FastMode(true).Size(270, 555).Rows(binance.GetWsAggTradeTable()...),
								),
							),
						),
					),
					giu.SplitLayout(giu.DirectionHorizontal, 800, //V
						giu.TabBar().TabItems(
							giu.TabItem("终端").Layout(
								giu.Label(console.ConsoleInstance.Read()),
							),
							giu.TabItem("历史订单").Layout(),
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
										SetColor(giu.StyleColorBorder, global.BLUE).
										SetStyle(giu.StyleVarFramePadding, 10, 10).
										To(
											giu.Row(
												giu.Button(fmt.Sprintf("单仓数量: %s: %.8f ",
													binance.AccountInstance.One.Asset, global.AverageSymbol1Amount)),
											),
										),
									giu.Row(
										giu.Style().
											SetColor(giu.StyleColorBorder, global.BLUE).
											SetStyle(giu.StyleVarFramePadding, 10, 10).
											To(
												giu.Row(
													giu.Button("交易模式(Tab): "+global.GetTradeMode()).OnClick(func() { global.ReverseTradeMode() }),
												),
											),
										giu.Style().
											SetColor(giu.StyleColorBorder, global.HotKeyColor).
											SetStyle(giu.StyleVarFramePadding, 10, 10).
											To(
												giu.Button("交易热键状态(空格): "+global.GetHotKeyStatus()).OnClick(func() { global.ReverseHotKeyStatus() }),
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
								giu.Column(
									giu.Table().FastMode(true).Size(270, 420).Rows(binance.GetWsPartialDepthBuyTable()...),
									giu.Label(fmt.Sprintf("实时价格: %s", binance.AggTradePrice)),
									giu.Table().FastMode(true).Size(270, 420).Rows(binance.GetWsPartialDepthSaleTable()...),
								),
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
