package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"

	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

var (
	globalHistoryC = make(chan *libBinance.TradeV3)

	historyTable []*giu.TableRowWidget
)

func GetHistoryTable() []*giu.TableRowWidget {
	return historyTable
}

func buildHistoryTable() []*giu.TableRowWidget {
	var (
		rows       []*giu.TableRowWidget
		historyNew = <-globalHistoryC
	)

	rows = append(rows, giu.TableRow(
		giu.Label("时间"),
		giu.Label("交易对"),
		giu.Label("成交均价"),
		giu.Label("成交额"),
		giu.Label("手续费"),
	).BgColor(global.PURPLE))

	eColor := global.RED
	if historyNew.IsBuyer {
		eColor = global.GREEN
	}

	feeColor := global.RED
	if historyNew.CommissionAsset == "BNB" {
		feeColor = global.WHITE
	}

	now := parseTimeStampMs(historyNew.Time)
	rows = append(rows, giu.TableRow(
		giu.Label(fmt.Sprintf("%2d日 %v:%v:%v", now.Day(), now.Hour(), now.Minute(), now.Second())),
		giu.Label(fmt.Sprintf("%v", historyNew.Symbol)),
		giu.Style().SetColor(giu.StyleColorText, eColor).To(
			giu.Label(fmt.Sprintf("%v", historyNew.Price)),
		),
		giu.Label(fmt.Sprintf("%v", historyNew.QuoteQuantity)),
		giu.Style().SetColor(giu.StyleColorText, feeColor).To(
			giu.Label(fmt.Sprintf("%v %v", historyNew.Commission, historyNew.CommissionAsset)),
		),
	))
	if len(historyTable) > 1 {
		if len(historyTable) < 30 {
			return append(rows, historyTable[1:]...)
		} else {
			return append(rows, historyTable[1:30]...)
		}
	}

	return append(rows, historyTable...)
}

func updateTradeHistory() {
	res, err := GetClient().NewListTradesService().Symbol(AccountInstance.Symbol).Limit(30).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	for _, v := range res {
		globalHistoryC <- v
	}
}

func parseTimeStampMs(timestamp int64) time.Time {
	return time.Unix(timestamp/1000, timestamp%1000)
}

func parseTimeStampS(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}
