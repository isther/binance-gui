package binance

import "github.com/isther/binanceGui/global"

func StartWebSocketStream() {
	var (
		wsPartialDepthServerDoneC chan struct{}
		wsPartialDepthServerStopC chan struct{}

		wsAggTradeServerDoneC chan struct{}
		wsAggTradeServerStopC chan struct{}

		wsUpdateAccountDoneC chan struct{}
		wsUpdateAccountStopC chan struct{}
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

	// update account
	go StartBuildOrderTable()

	wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
	wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()
	wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
	StartUpdateAccount()
	UpdateAverageAmount()
	for {
		select {
		case symbol := <-global.FreshC:
			// Clear Order
			ResetOrders()
			AccountInstance.Symbol = symbol
			go func() {
				wsPartialDepthServerStopC <- struct{}{}
				<-wsPartialDepthServerDoneC

				wsPartialDepthServerDoneC, wsPartialDepthServerStopC = runOneWsPartialDepth()
			}()

			wsAggTradeServerStopC <- struct{}{}
			<-wsAggTradeServerDoneC

			wsAggTradeServerDoneC, wsAggTradeServerStopC = runOneAggTradeDepth()

			go func() {
				wsUpdateAccountStopC <- struct{}{}
				<-wsUpdateAccountDoneC

				wsUpdateAccountDoneC, wsUpdateAccountStopC = AccountInstance.WsUpdateAccount()
			}()

			StartUpdateAccount()
			UpdateAverageAmount()
		}
	}
}
