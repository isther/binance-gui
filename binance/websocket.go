package binance

import "github.com/isther/binanceGui/global"

func StartWebSocketStream() {
	var (
		wsPartialDepthServerDoneC chan struct{}
		wsPartialDepthServerStopC chan struct{}

		wsAggTradeServerDoneC chan struct{}
		wsAggTradeServerStopC chan struct{}
	)

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

	wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
	wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
	for {
		select {
		case symbol := <-global.FreshC:
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

			AccountInstance.ExchangeInfo()
		}
	}

}
