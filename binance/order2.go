package binance

import (
	"fmt"
	"image/color"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/hotkey"
)

var (
	globalWsPartialDepthServer = make(chan *libBinance.WsPartialDepthEvent)
	depthTable                 *libBinance.WsPartialDepthEvent
	wsPartialDepthTable        []*giu.TableRowWidget
)

func GetWsPartialDepthTable() []*giu.TableRowWidget {
	return wsPartialDepthTable
}

func StartWsPartialDepthServer() {
	var (
		doneC chan struct{}
		stopC chan struct{}
	)

	go func() {
		for {
			wsPartialDepthTable = buildWsPartialDepthTable()
		}
	}()

	doneC, stopC = runOneWsPartialDepth()
	for {
		select {
		case symbol := <-global.FreshC:
			stopC <- struct{}{}
			<-doneC

			AccountInstance.Symbol = symbol
			doneC, stopC = runOneWsPartialDepth()
			AccountInstance.ExchangeInfo()
		}
	}
}

func runOneWsPartialDepth() (chan struct{}, chan struct{}) {
	var (
		err   error
		doneC chan struct{}
		stopC chan struct{}
	)

	wsDepthHandler := func(event *libBinance.WsPartialDepthEvent) {
		depthTable = event
		globalWsPartialDepthServer <- event
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneC, stopC, err = libBinance.WsPartialDepthServe100Ms(AccountInstance.Symbol, fmt.Sprintf("%d", global.Levels), wsDepthHandler, errHandler)

	if err != nil {
		fmt.Println(err)
		return doneC, stopC
	}
	return doneC, stopC
}

func buildWsPartialDepthTable() []*giu.TableRowWidget {
	var (
		length = global.Levels + 1
	)
	eventNew := <-globalWsPartialDepthServer

	rows := make([]*giu.TableRowWidget, length*2)

	rows[0] = giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(&(color.RGBA{0x33, 0x33, 0xff, 0xff}))

	rows[21] = giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(&(color.RGBA{0x33, 0x33, 0xff, 0xff}))

	if eventNew == nil {
		return rows
	}

	for i := range eventNew.Asks {
		price, quantity, err := eventNew.Asks[i].Parse()
		if err != nil {
			console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		}
		rows[length-i-1] = giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				To(
					giu.Label(fmt.Sprintf("%c", hotkey.HotKeySale[i])),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				To(
					giu.Label(priceFloat648Point(price)),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, color.RGBA{0xff, 0x33, 0x00, 0xff}).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		)

		price, quantity, err = eventNew.Bids[i].Parse()
		if err != nil {
			console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		}
		rows[length+i+1] = giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				To(
					giu.Label(fmt.Sprintf("%c", hotkey.HotKeyBuy[i])),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				To(
					giu.Label(priceFloat648Point(price)),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, color.RGBA{0x66, 0xcc, 0x00, 0xff}).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		)
	}
	return rows
}
