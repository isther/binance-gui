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
	wsDepthTable               *libBinance.WsPartialDepthEvent

	wsPartialDepthTable []*giu.TableRowWidget
)

func GetWsPartialDepthBuyTable() []*giu.TableRowWidget {
	if wsPartialDepthTable == nil {
		return wsPartialDepthTable
	}
	return wsPartialDepthTable[0:21]
}

func GetWsPartialDepthSaleTable() []*giu.TableRowWidget {
	if wsPartialDepthTable == nil {
		return wsPartialDepthTable
	}
	return wsPartialDepthTable[21:]
}

func runOneWsPartialDepth() (chan struct{}, chan struct{}) {
	var (
		err   error
		doneC chan struct{}
		stopC chan struct{}
	)

	wsDepthHandler := func(event *libBinance.WsPartialDepthEvent) {
		wsDepthTable = event
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
	).BgColor(global.PURPLE)

	rows[21] = giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE)

	if eventNew == nil {
		return rows
	}

	for i := range eventNew.Asks {
		var bgColor color.RGBA = global.WHITE
		if (i)%5 == 0 {
			bgColor = global.Order2Bg
		}

		price, quantity, err := eventNew.Asks[i].Parse()
		if err != nil {
			console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		}
		rows[length-i-1] = giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, bgColor).
				To(
					giu.Label(fmt.Sprintf("%c", hotkey.HotKeySale[i])),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, orderColorSet(price, quantity)).
				To(
					giu.Label(priceFloat648Point(fmt.Sprintf("%.8f", price))),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.RED).
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
				SetColor(giu.StyleColorText, bgColor).
				To(
					giu.Label(fmt.Sprintf("%c", hotkey.HotKeyBuy[i])),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, orderColorSet(price, quantity)).
				To(
					giu.Label(priceFloat648Point(fmt.Sprintf("%.8f", price))),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.GREEN).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		)
	}
	return rows
}

func orderColorSet(price, quantity float64) color.RGBA {
	var (
		priceColor = global.WHITE
		ff         = price * quantity
	)

	if ff >= float64(global.Order2BigOrderReminder[4]) {
		priceColor = global.BLACK
	} else if ff >= float64(global.Order2BigOrderReminder[3]) {
		priceColor = global.BLUE2
	} else if ff >= float64(global.Order2BigOrderReminder[2]) {
		priceColor = global.YELLOW
	} else if ff >= float64(global.Order2BigOrderReminder[1]) {
		priceColor = global.RED
	} else {
		priceColor = global.WHITE
	}

	return priceColor
}
