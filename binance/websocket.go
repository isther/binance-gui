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

	go func() {
		wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
		for {
			<-wsPartialDepthServerDoneC
			console.ConsoleInstance.Write("Reload partial depth...")
			wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
		}
	}()

	go func() {
		wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
		for {
			<-wsAggTradeServerDoneC
			console.ConsoleInstance.Write("Reload aggTrade...")
			wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
		}
	}()

	go func() {
		wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
		go StartUpdateAccount()
		for {
			<-wsUpdateAccountDoneC
			console.ConsoleInstance.Write("Reload updateAccount...")
			wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
		}
	}()

	go func() {
		wsUpdateTickerDoneC, wsUpdateTickerStopC = UpdateWsTickerTable()
		for {
			<-wsUpdateTickerDoneC
			console.ConsoleInstance.Write("Reload ticker...")
			wsUpdateTickerDoneC, wsUpdateTickerStopC = UpdateWsTickerTable()
		}
	}()

	go func() {
		for {
			select {
			case symbol := <-global.FreshC:
				go StartUpdateAccount()
				console.ConsoleInstance.Write(fmt.Sprintf("New Symbol: %v", symbol))
				AccountInstance.Symbol = symbol
				wsPartialDepthServerStopC <- struct{}{}
				wsAggTradeServerStopC <- struct{}{}
				wsUpdateAccountStopC <- struct{}{}
				wsUpdateTickerStopC <- struct{}{}
			case <-global.ReconnectWsPartialDepthC:
				wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
			case <-global.ReconnectWsAggTradeC:
				wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
			case <-global.ReconnectWsAccountC:
				wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
				go StartUpdateAccount()
			case <-global.ReconnectWsTickerC:
				wsUpdateTickerDoneC, wsUpdateTickerStopC = UpdateWsTickerTable()
			}
		}
	}()
}
