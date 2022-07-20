package binance

import (
	"fmt"

	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

func StartWebSocketStream() {
	var (
		wsPartialDepthServerDoneC chan struct{}
		wsPartialDepthServerStopC chan struct{}

		wsAggTradeServerDoneC chan struct{}
		wsAggTradeServerStopC chan struct{}

		wsUpdateAccountDoneC chan struct{}
		wsUpdateAccountStopC chan struct{}
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

	wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
	wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
	wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
	StartUpdateWsTickerTable()
	StartUpdateAccount()
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

			StartUpdateAccount()
		}
	}
}
