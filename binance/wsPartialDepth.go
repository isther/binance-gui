package binance

//order2

import (
	"fmt"
	"image/color"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"

	"github.com/isther/binance/conf"
	"github.com/isther/binance/global"
)

var (
	globalWsPartialDepthServeEventC = make(chan *libBinance.WsPartialDepthEvent, 1000)
)

func StartWsPartialDepthServer() {
	var (
		doneC chan struct{}
		stopC chan struct{}
	)

	doneC, stopC = runOneWsPartialDepth(true)
	for {
		select {
		case symbol := <-global.FreshC:
			if symbolExist(symbol) {
				stopC <- struct{}{}
				<-doneC

				global.Symbol = symbol
				doneC, stopC = runOneWsPartialDepth(true)
			}
		}
	}
}

func BuildWsPartialDepthServeRows() []*giu.TableRowWidget {
	var (
		length   = 1 + global.Levels
		eventNew = <-globalWsPartialDepthServeEventC
	)

	rows := make([]*giu.TableRowWidget, length*2)

	rows[0] = giu.TableRow(
		giu.Label("id"),
		giu.Label("Price"),
		giu.Label("Turn over"),
	).BgColor(&(color.RGBA{0x33, 0x33, 0xff, 0xff}))

	rows[21] = giu.TableRow(
		giu.Label("id"),
		giu.Label("Price"),
		giu.Label("Turn over"),
	).BgColor(&(color.RGBA{0x33, 0x33, 0xff, 0xff}))

	if eventNew == nil {
		return rows
	}

	for i := range eventNew.Asks {
		price, quantity, err := eventNew.Asks[i].Parse()
		if err != nil {
			panic(err)
		}
		rows[length-i-1] = giu.TableRow(
			giu.Label(conf.Conf.HotKey.Sale[i]),

			giu.Label(priceFloat648Point(price)),
			giu.Style().
				SetColor(giu.StyleColorText, color.RGBA{0xff, 0x33, 0x00, 0xff}).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		)
	}

	for i := range eventNew.Bids {
		price, quantity, err := eventNew.Bids[i].Parse()
		if err != nil {
			panic(err)
		}
		rows[length+i+1] = giu.TableRow(
			giu.Label(conf.Conf.HotKey.Buy[i]),

			giu.Label(priceFloat648Point(price)),
			giu.Style().
				SetColor(giu.StyleColorText, color.RGBA{0x66, 0xcc, 0x00, 0xff}).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		)
	}

	return rows
}

func runOneWsPartialDepth(need100Ms bool) (chan struct{}, chan struct{}) {
	var (
		err   error
		doneC chan struct{}
		stopC chan struct{}
	)

	wsDepthHandler := func(event *libBinance.WsPartialDepthEvent) {
		globalWsPartialDepthServeEventC <- event
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	if need100Ms {
		doneC, stopC, err = libBinance.WsPartialDepthServe100Ms(global.Symbol, fmt.Sprintf("%d", global.Levels), wsDepthHandler, errHandler)
	} else {
		doneC, stopC, err = libBinance.WsPartialDepthServe(global.Symbol, fmt.Sprintf("%d", global.Levels), wsDepthHandler, errHandler)
	}
	if err != nil {
		fmt.Println(err)
		return doneC, stopC
	}
	return doneC, stopC
}
