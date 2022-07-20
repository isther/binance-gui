package binance

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

var (
	globalWsAggTradeServerC = make(chan *libBinance.WsAggTradeEvent)
	wsAggTradeTable         []*giu.TableRowWidget

	AggTradePrice      string
	AggTradePriceColor color.Color = global.RED
)

func GetWsAggTradeTable() []*giu.TableRowWidget {
	return wsAggTradeTable
}

func runOneAggTradeDepth() (chan struct{}, chan struct{}) {
	var (
		err   error
		doneC chan struct{}
		stopC chan struct{}
	)

	wsAggTradeHandler := func(event *libBinance.WsAggTradeEvent) {
		AggTradePrice = priceFloat648Point(event.Price)
		if !event.IsBuyerMaker {
			AggTradePriceColor = global.GREEN
		} else {
			AggTradePriceColor = global.RED
		}
		globalWsAggTradeServerC <- event
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneC, stopC, err = libBinance.WsAggTradeServe(AccountInstance.Symbol, wsAggTradeHandler, errHandler)

	if err != nil {
		fmt.Println(err)
		return doneC, stopC
	}
	return doneC, stopC
}

func buildAggTradeTable() []*giu.TableRowWidget {
	var (
		rows        []*giu.TableRowWidget
		aggTradeNew = <-globalWsAggTradeServerC
	)

	rows = append(rows, giu.TableRow(
		giu.Label("Time"),
		giu.Label("Price"),
		giu.Label("TurnOver"),
	).BgColor(global.PURPLE))

	price, err := strconv.ParseFloat(aggTradeNew.Price, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	quantity, err := strconv.ParseFloat(aggTradeNew.Quantity, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}

	if global.IsShieldSmallOrder && price*quantity < 1000 {
		return wsAggTradeTable
	}

	eColor := global.GREEN
	if aggTradeNew.IsBuyerMaker {
		eColor = global.RED
	}

	timeStr := time.Unix(aggTradeNew.Time/1e3, 0)
	rows = append(rows, giu.TableRow(
		giu.Label(fmt.Sprintf("%v:%v", timeStr.Minute(), timeStr.Second())),
		giu.Style().SetColor(giu.StyleColorText, aggTradeColorSet(price, quantity)).
			To(
				giu.Label(fmt.Sprintf("%v", priceFloat648Point(aggTradeNew.Price))),
			),
		giu.Style().
			SetColor(giu.StyleColorText, eColor).
			To(
				giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
			),
	))
	if len(wsAggTradeTable) > 1 {
		if len(wsAggTradeTable) < 100 {
			return append(rows, wsAggTradeTable[1:]...)
		} else {
			return append(rows, wsAggTradeTable[1:100]...)
		}
	}

	return append(rows, wsAggTradeTable...)
}

func aggTradeColorSet(price, quantity float64) color.RGBA {
	var (
		priceColor = global.WHITE
		ff         = price * quantity
	)

	if ff >= float64(global.AggTradeBigOrderReminder[4]) {
		priceColor = global.BLACK
	} else if ff >= float64(global.AggTradeBigOrderReminder[3]) {
		priceColor = global.BLUE2
	} else if ff >= float64(global.AggTradeBigOrderReminder[2]) {
		priceColor = global.YELLOW
	} else if ff >= float64(global.AggTradeBigOrderReminder[1]) {
		priceColor = global.RED
	} else {
		priceColor = global.WHITE
	}

	return priceColor
}
