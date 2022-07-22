package binance

import (
	"fmt"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

func init() {

	// websocket
	libBinance.WebsocketKeepalive = true
}

func StartWebSocketStream() {
	var (
		wsPartialDepthServerDoneC chan struct{}
		wsPartialDepthServerStopC chan struct{}

		wsAggTradeServerDoneC chan struct{}
		wsAggTradeServerStopC chan struct{}

		wsUpdateAccountDoneC chan struct{}
		wsUpdateAccountStopC chan struct{}

		wsUpdateTickerDoneC chan struct{}
		wsUpdateTickerStopC chan struct{}
	)
	updateTime()

	go func() {
		for {
			wsPartialDepthTable = buildWsPartialDepthTable()
		}
	}()

	go func() {
		for {
			wsAggTradeTable = buildAggTradeTable()
		}
	}()

	go func() {
		for {
			historyTable = buildHistoryTable()
		}
	}()

	go func() {
		for {
			tickerBTCTable, tickerUSDTTable, tickerBUSDTable = buildTickerTable()
		}
	}()

	go func() {
		for {
			updateMaps()
		}
	}()

	wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
	wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
	wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
	wsUpdateTickerDoneC, wsUpdateTickerStopC = UpdateWsTickerTable()
	StartUpdateAccount()
	go func() {
		for {
			select {
			case symbol := <-global.FreshC:
				console.ConsoleInstance.Write(fmt.Sprintf("New Symbol: %v", symbol))
				// Clear Order
				AccountInstance.Symbol = symbol
				go func() {
					wsPartialDepthServerStopC <- struct{}{}
					<-wsPartialDepthServerDoneC

					wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
				}()

				go func() {
					wsAggTradeServerStopC <- struct{}{}
					<-wsAggTradeServerDoneC

					wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
				}()

				go func() {
					wsUpdateAccountStopC <- struct{}{}
					<-wsUpdateAccountDoneC

					wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
				}()

				go func() {
					wsUpdateTickerStopC <- struct{}{}
					<-wsUpdateTickerDoneC

					wsUpdateTickerDoneC, wsUpdateTickerStopC = UpdateWsTickerTable()
				}()

				StartUpdateAccount()
			}
		}
	}()

	go func() {
		for {
			select {
			case <-global.ReConnectWsPartialDepth:
				wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
			case <-global.ReConnectWsAggTrade:
				wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
			case <-global.ReConnectWsUpdateAccount:
				wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
				StartUpdateAccount()
			case <-global.ReConnectWsTickerTable:
				wsUpdateTickerDoneC, wsUpdateTickerStopC = UpdateWsTickerTable()
			}
		}
	}()
}
