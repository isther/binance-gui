package gui

//order2

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"

	"github.com/isther/binance/global"
)

var (
	globalWsPartialDepthServeEventC = make(chan *libBinance.WsPartialDepthEvent)
)

type WsPartialDepthServer struct {
	Symbol string
	Levels int
}

func NewWsPartialDepthServer(symbol string, levels int) *WsPartialDepthServer {
	return &WsPartialDepthServer{
		Symbol: symbol,
		Levels: levels,
	}
}

func (ws *WsPartialDepthServer) StartWsPartialDepth(need100Ms bool) chan struct{} {
	var (
		err   error
		stopC chan struct{}
	)

	wsDepthHandler := func(event *libBinance.WsPartialDepthEvent) {
		// fmt.Println(event)
		time.Sleep(time.Millisecond * 100)
		globalWsPartialDepthServeEventC <- event
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	if need100Ms {
		_, stopC, err = libBinance.WsPartialDepthServe100Ms(ws.Symbol, strconv.Itoa(ws.Levels), wsDepthHandler, errHandler)
	} else {
		_, stopC, err = libBinance.WsPartialDepthServe(ws.Symbol, strconv.Itoa(ws.Levels), wsDepthHandler, errHandler)
	}
	if err != nil {
		fmt.Println(err)
		return stopC
	}

	return stopC
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
		giu.Label("Quantity"),
	).BgColor(&(color.RGBA{0x33, 0x33, 0xff, 0xff}))

	rows[21] = giu.TableRow(
		giu.Label("id"),
		giu.Label("Price"),
		giu.Label("Quantity"),
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
			giu.Label(fmt.Sprintf("%d", i+1)),

			giu.Label(fmt.Sprintf("%f", price)),
			giu.Label(fmt.Sprintf("%f", quantity)),
		)
	}

	for i := range eventNew.Bids {
		price, quantity, err := eventNew.Bids[i].Parse()
		if err != nil {
			panic(err)
		}
		rows[length+i+1] = giu.TableRow(
			giu.Label(fmt.Sprintf("%d", i+1)),

			giu.Label(fmt.Sprintf("%f", price)),
			giu.Label(fmt.Sprintf("%f", quantity)),
		)
	}

	giu.Update()

	return rows
}
