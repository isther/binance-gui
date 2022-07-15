package binance

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
)

var (
	globalWsAggTradeServerC = make(chan *libBinance.WsAggTradeEvent)
	wsAggTradeTable         []*giu.TableRowWidget

	AggTradePrice string
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
	AggTradePrice = priceFloat648Point(aggTradeNew.Price)

	rows = append(rows, giu.TableRow(
		giu.Label("时间"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(&(color.RGBA{0x33, 0x33, 0xff, 0xff})))

	price, err := strconv.ParseFloat(aggTradeNew.Price, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	quantity, err := strconv.ParseFloat(aggTradeNew.Quantity, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}

	eColor := color.RGBA{0xff, 0x33, 0x00, 0xff}
	if aggTradeNew.IsBuyerMaker {
		eColor = color.RGBA{0x66, 0xcc, 0x00, 0xff}
	}

	timeStr := time.Unix(aggTradeNew.Time/1e3, 0)
	rows = append(rows, giu.TableRow(
		giu.Label(fmt.Sprintf("%v:%v", timeStr.Minute(), timeStr.Second())),
		giu.Label(fmt.Sprintf("%v", priceFloat648Point(aggTradeNew.Price))),
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
