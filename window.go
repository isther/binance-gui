package main

import (
	"fmt"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/binance"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/orderlist"
)

var (
	volatilityRatiosInputSize float32 = 100

	symbol  = binance.AccountInstance.Symbol
	symbol1 = binance.AccountInstance.One.Asset
	symbol2 = binance.AccountInstance.Two.Asset
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
				giu.Menu("快捷购买BNB").Layout(
					giu.Button("BUSD购买").OnClick(func() { go binance.NewTradeBNB("BUSD").Trade() }),
					giu.Button("USDT购买").OnClick(func() { go binance.NewTradeBNB("USDT").Trade() }),
				),
				giu.Button("BUSD一键购买USDT").OnClick(func() { go binance.NewTradeBUSDAndUSDT().Trade() }),
				giu.Style().
					SetColor(giu.StyleColorText, global.RED).
					To(
						giu.Column(
							giu.Button("服务器延迟: "+global.Ping),
						),
					),
				giu.Button("交易模式([]): "+global.GetTradeMode()),
				giu.Style().
					SetColor(giu.StyleColorBorder, global.HotKeyColor).
					To(
						giu.Column(
							giu.Button("交易热键状态(空格): "+global.GetHotKeyStatus()).OnClick(func() { global.ReverseHotKeyStatus() }),
						),
					),
				giu.Button("币安系统时间: "+binance.TimeString),
			),
			giu.SplitLayout(giu.DirectionHorizontal, 1200, // H
				giu.SplitLayout(giu.DirectionHorizontal, 300,
					giu.SplitLayout(giu.DirectionVertical, 600,
						giu.TabBar().TabItems(
							giu.TabItem("当前挂单").Layout(
								giu.SplitLayout(giu.DirectionVertical, 280, //V
									giu.Table().Freeze(0, 1).FastMode(true).Size(270, 300).Rows(orderlist.GetOpenSaleOrdersTable()...),
									giu.Table().Freeze(0, 1).FastMode(true).Size(270, 300).Rows(orderlist.GetOpenBuyOrdersTable()...),
								),
							),
						),
						giu.TabBar().TabItems(
							giu.TabItem("交易对").Layout(
								giu.Child().Layout(
									giu.TabBar().TabItems(
										giu.TabItem("BUSD").Layout(
											giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetTickerBUSDTable()...),
										),
										giu.TabItem("BTC").Layout(
											giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetTickerBTCTable()...),
										),
										giu.TabItem("USDT").Layout(
											giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetTickerUSDTTable()...),
										),
									),
								),
							),
						),
					),
					giu.SplitLayout(giu.DirectionHorizontal, 300,
						giu.SplitLayout(giu.DirectionVertical, 600,
							giu.SplitLayout(giu.DirectionVertical, 300,
								giu.TabBar().TabItems(
									giu.TabItem("成交明细").Layout(
										giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetWsAggTradeTable()...),
									),
								),
								giu.TabBar().TabItems(
									giu.TabItem("预警").Layout(
										giu.SplitLayout(giu.DirectionVertical, 125,
											giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetEarlyWaringTable1m()...),
											giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetEarlyWaringTable3m()...),
										),
									),
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
										giu.Style().
											SetColor(giu.StyleColorText, global.RED).
											SetFontSize(20).
											To(
												giu.Label(fmt.Sprintf("Trading Pair: %v", binance.AccountInstance.Symbol)),
											),
										giu.Row(
											giu.Label("输入切换交易对: "),
											giu.InputText(&symbol1).Size(72.5),
											giu.Label("/"),
											giu.InputText(&symbol2).Size(72.5),
											giu.Button("确定").OnClick(freshSymbol),
										),
										giu.Row(
											giu.Label("分仓数: "),
											giu.InputInt(&global.Average).Size(120),
											giu.Button("确定(Enter)").OnClick(binance.UpdateAverageAmount),
										),
										giu.Style().
											SetColor(giu.StyleColorBorder, global.BLUE).
											SetStyle(giu.StyleVarFramePadding, 10, 10).
											To(
												giu.Column(
													giu.Button(fmt.Sprintf("%s: %.8f",
														binance.AccountInstance.One.Asset, global.AverageSymbol1Amount)),
													giu.Button(fmt.Sprintf("%s : %.8f",
														binance.AccountInstance.Two.Asset, global.AverageSymbol2Amount)),
												),
											),
									),
								)),
							),
						),
						giu.SplitLayout(giu.DirectionHorizontal, 300,
							giu.TabBar().TabItems(
								giu.TabItem("订单簿2").Layout(
									giu.Column(
										giu.Table().FastMode(true).Size(270, 420).Rows(binance.GetWsPartialDepthBuyTable()...),
										giu.Style().
											SetColor(giu.StyleColorText, binance.AggTradePriceColor).
											To(
												giu.Label(fmt.Sprintf("实时价格: %s", binance.AggTradePrice)),
											),
										giu.Table().FastMode(true).Size(270, 420).Rows(binance.GetWsPartialDepthSaleTable()...),
									),
								),
							),
							giu.TabBar().TabItems(giu.TabItem("订单簿1").Layout(
								giu.Column(
									giu.Table().FastMode(true).Size(270, 420).Rows(binance.GetHttpDepthBuyTable()...),
									giu.Style().
										SetColor(giu.StyleColorText, binance.CostColor).
										To(
											giu.Label(fmt.Sprintf("持仓成本: %v", binance.CostInstance.UpdateAverageCode())),
										),
									giu.Table().FastMode(true).Size(270, 420).Rows(binance.GetHttpDepthSaleTable()...),
								),
							)),
						),
					),
				),
				giu.SplitLayout(giu.DirectionVertical, 600,
					giu.TabBar().TabItems(
						// giu.TabItem("K线").Layout(),
						giu.TabItem("终端").Layout(
							giu.Label(console.ConsoleInstance.Read()),
						),
						giu.TabItem("参数设置").Layout(giu.Style().
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
									giu.Label("全局下单"),
									giu.Row(
										giu.Label("分仓买F1:"),
										giu.InputFloat(&global.VolatilityRatiosF1).Size(volatilityRatiosInputSize),
										giu.Label("全仓买F2:"),
										giu.InputFloat(&global.VolatilityRatiosF2).Size(volatilityRatiosInputSize),
										giu.Label("分仓卖F5:"),
										giu.InputFloat(&global.VolatilityRatiosF5).Size(volatilityRatiosInputSize),
										giu.Label("全仓卖F6:"),
										giu.InputFloat(&global.VolatilityRatiosF6).Size(volatilityRatiosInputSize),
									),
									giu.Row(
										giu.Label("撤单清仓卖F12:"),
										giu.InputFloat(&global.VolatilityRatiosF12).Size(volatilityRatiosInputSize),
									),
									giu.Row(
										giu.Label("大单提醒档位设置"),
										giu.Label("成交明细是否蔽小单:"),
										giu.Checkbox("", &global.IsShieldSmallOrder),
									),
									giu.Row(
										giu.Label("成交明细: "),
										giu.Label("一档"),
										giu.InputInt(&global.AggTradeBigOrderReminder[1]).Size(volatilityRatiosInputSize),
										giu.Label("二档"),
										giu.InputInt(&global.AggTradeBigOrderReminder[2]).Size(volatilityRatiosInputSize),
										giu.Label("三档"),
										giu.InputInt(&global.AggTradeBigOrderReminder[3]).Size(volatilityRatiosInputSize),
										giu.Label("四档"),
										giu.InputInt(&global.AggTradeBigOrderReminder[4]).Size(volatilityRatiosInputSize),
									),
									giu.Row(
										giu.Label("订单簿2: "),
										giu.Label("一档"),
										giu.InputInt(&global.Order2BigOrderReminder[1]).Size(volatilityRatiosInputSize),
										giu.Label("二档"),
										giu.InputInt(&global.Order2BigOrderReminder[2]).Size(volatilityRatiosInputSize),
										giu.Label("三档"),
										giu.InputInt(&global.Order2BigOrderReminder[3]).Size(volatilityRatiosInputSize),
										giu.Label("四档"),
										giu.InputInt(&global.Order2BigOrderReminder[4]).Size(volatilityRatiosInputSize),
									),
									giu.Row(
										giu.Label("订单簿1: "),
										giu.Label("一档"),
										giu.InputInt(&global.Order1BigOrderReminder[1]).Size(volatilityRatiosInputSize),
										giu.Label("二档"),
										giu.InputInt(&global.Order1BigOrderReminder[2]).Size(volatilityRatiosInputSize),
										giu.Label("三档"),
										giu.InputInt(&global.Order1BigOrderReminder[3]).Size(volatilityRatiosInputSize),
										giu.Label("四档"),
										giu.InputInt(&global.Order1BigOrderReminder[4]).Size(volatilityRatiosInputSize),
									),
									giu.Row(
										giu.Label("筛选与预警   "),
										giu.Label("开启预警:"),
										giu.Checkbox("", &global.EarlyWarning),
									),
									giu.Row(
										giu.Label("1m: "),
										giu.Label("涨幅"),
										giu.InputFloat(&global.EarlyWarning1mAmplitude).Size(volatilityRatiosInputSize),
										giu.Label("成交额"),
										giu.InputFloat(&global.EarlyWarning1mTurnOver).Size(volatilityRatiosInputSize),
									),
									giu.Row(
										giu.Label("3m: "),
										giu.Label("涨幅"),
										giu.InputFloat(&global.EarlyWarning3mAmplitude).Size(volatilityRatiosInputSize),
										giu.Label("成交额"),
										giu.InputFloat(&global.EarlyWarning3mTurnOver).Size(volatilityRatiosInputSize),
									),
								),
							)),
					),
					giu.TabBar().TabItems(
						giu.TabItem("历史订单").Layout(
							giu.Table().Freeze(0, 1).FastMode(true).Rows(binance.GetHistoryTable()...),
						),
					),
				),
			),
		)
}
