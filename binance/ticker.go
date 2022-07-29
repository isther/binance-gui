package binance

import (
	"context"
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

type sortType int

type myTicker struct {
	asset      string
	turnOver   float64
	percentage float64
}

const (
	_ = iota
	sortByPercentageDescend
	sortByTurnOverDescend
	sortByAssetDescend
)

var (
	globalTickerC     = make(chan libBinance.WsAllMiniMarketsStatEvent)
	buildTickerTableC = make(chan struct{})
	updateTableC      = make(chan struct{})

	openGlobalTickerC = false

	mapBTC  = make(map[string]*myTicker)
	mapUSDT = make(map[string]*myTicker)
	mapBUSD = make(map[string]*myTicker)

	tickerBTCTable  []*giu.TableRowWidget
	tickerUSDTTable []*giu.TableRowWidget
	tickerBUSDTable []*giu.TableRowWidget

	sortTypeNow sortType = sortByPercentageDescend
)

func GetTickerBTCTable() []*giu.TableRowWidget {
	return tickerBTCTable
}

func GetTickerUSDTTable() []*giu.TableRowWidget {
	return tickerUSDTTable
}

func GetTickerBUSDTable() []*giu.TableRowWidget {
	return tickerBUSDTable
}

func UpdateWsTickerTable() (chan struct{}, chan struct{}) {
	console.ConsoleInstance.Write(fmt.Sprint("Run ticker websocket..."))
	cyclePing()

	httpGetTicker()
	buildTickerTableC <- struct{}{}

	return wsGetTicker()
}

func httpGetTicker() {
	tickers, err := GetClient().NewListPriceChangeStatsService().Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	for i := range tickers {
		var ticker = tickers[i]

		// lastPrice, _ := strconv.ParseFloat(ticker.LastPrice, 64)
		// openPrice, _ := strconv.ParseFloat(ticker.OpenPrice, 64)
		priceChangePercent, _ := strconv.ParseFloat(ticker.PriceChangePercent, 64)
		turnOver, _ := strconv.ParseFloat(ticker.QuoteVolume, 64)

		if strings.HasSuffix(ticker.Symbol, "BTC") {
			asset := ticker.Symbol[:len(ticker.Symbol)-3]
			mapBTC[asset] = &myTicker{
				asset:      asset,
				turnOver:   turnOver / 1000000,
				percentage: priceChangePercent,
			}
		} else if strings.HasSuffix(ticker.Symbol, "USDT") {
			asset := ticker.Symbol[:len(ticker.Symbol)-4]
			mapUSDT[asset] = &myTicker{
				asset:      asset,
				turnOver:   turnOver / 1000000,
				percentage: priceChangePercent,
			}
		} else if strings.HasSuffix(ticker.Symbol, "BUSD") {
			asset := ticker.Symbol[:len(ticker.Symbol)-4]
			mapBUSD[asset] = &myTicker{
				asset:      asset,
				turnOver:   turnOver / 1000000,
				percentage: priceChangePercent,
			}
		}
	}
}

func wsGetTicker() (chan struct{}, chan struct{}) {
	var (
		doneC chan struct{}
		stopC chan struct{}
		err   error
	)

	wsHandler := func(event libBinance.WsAllMiniMarketsStatEvent) {
		if openGlobalTickerC {
			globalTickerC <- event
		}
	}

	errHandler := func(err error) {
		console.ConsoleInstance.Write(fmt.Sprintf("%v", err))
	}

	doneC, stopC, err = libBinance.WsAllMiniMarketsStatServe(wsHandler, errHandler)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("%v", err))
	}

	return doneC, stopC
}

func updateMaps() {
	<-updateTableC
	var (
		tickers = <-globalTickerC
	)

	for i := range tickers {
		var ticker = tickers[i]

		lastPrice, _ := strconv.ParseFloat(ticker.LastPrice, 64)
		openPrice, _ := strconv.ParseFloat(ticker.OpenPrice, 64)
		turnOver, _ := strconv.ParseFloat(ticker.QuoteVolume, 64)

		if strings.HasSuffix(ticker.Symbol, "BTC") {
			asset := ticker.Symbol[:len(ticker.Symbol)-3]
			mapBTC[asset] = &myTicker{
				asset:      asset,
				turnOver:   turnOver / 1000000,
				percentage: 100 * (lastPrice - openPrice) / openPrice,
			}
		} else if strings.HasSuffix(ticker.Symbol, "USDT") {
			asset := ticker.Symbol[:len(ticker.Symbol)-4]
			mapUSDT[asset] = &myTicker{
				asset:      asset,
				turnOver:   turnOver / 1000000,
				percentage: 100 * (lastPrice - openPrice) / openPrice,
			}
		} else if strings.HasSuffix(ticker.Symbol, "BUSD") {
			asset := ticker.Symbol[:len(ticker.Symbol)-4]
			mapBUSD[asset] = &myTicker{
				asset:      asset,
				turnOver:   turnOver / 1000000,
				percentage: 100 * (lastPrice - openPrice) / openPrice,
			}
		}
	}
	buildTickerTableC <- struct{}{}
}

func buildTickerTable() ([]*giu.TableRowWidget, []*giu.TableRowWidget, []*giu.TableRowWidget) {
	<-buildTickerTableC

	var (
		rowsBTC  []*giu.TableRowWidget
		rowsUSDT []*giu.TableRowWidget
		rowsBUSD []*giu.TableRowWidget

		rowSelectable = giu.TableRow(
			giu.Selectable("BaseAsset").Selected(false).OnClick(func() {
				sortTypeNow = sortByAssetDescend
				openGlobalTickerC = true
				updateTableC <- struct{}{}
			}),
			giu.Selectable("24hr TurnOver").Selected(false).OnClick(func() {
				sortTypeNow = sortByTurnOverDescend
				openGlobalTickerC = true
				updateTableC <- struct{}{}
			}),
			giu.Selectable("PriceChangePercent").Selected(false).OnClick(func() {
				sortTypeNow = sortByPercentageDescend
				openGlobalTickerC = true
				updateTableC <- struct{}{}
			}),
		).BgColor(global.PURPLE)
	)

	rowsBTC = append(rowsBTC, rowSelectable)
	rowsUSDT = append(rowsUSDT, rowSelectable)
	rowsBUSD = append(rowsBUSD, rowSelectable)

	// sort
	var (
		sortBTC  []string
		sortUSDT []string
		sortBUSD []string
	)
	if sortTypeNow == sortByPercentageDescend {
		sortBTC, sortUSDT, sortBUSD = sortByPercentage()
	} else if sortTypeNow == sortByTurnOverDescend {
		sortBTC, sortUSDT, sortBUSD = sortByTurnOver()
	} else if sortTypeNow == sortByAssetDescend {
		sortBTC, sortUSDT, sortBUSD = sortByAsset()
	}

	// build
	for i := range sortBTC {
		var (
			btc    = mapBTC[sortBTC[i]]
			eColor color.RGBA
		)

		if btc.percentage >= 0 {
			eColor = global.GREEN
		} else {
			eColor = global.RED
		}
		rowsBTC = append(rowsBTC, giu.TableRow(
			giu.Selectable(btc.asset).Selected(false).OnClick(func() {
				AccountInstance.One.Asset = btc.asset
				AccountInstance.Two.Asset = "BTC"
				global.FreshC <- btc.asset + "BTC"
			}),
			giu.Label(fmt.Sprintf("%.2fm", btc.turnOver)),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, eColor).
				To(
					giu.Label(fmt.Sprintf("%.2f%%", btc.percentage)),
				),
		))
	}

	for i := range sortUSDT {
		var (
			usdt   = mapUSDT[sortUSDT[i]]
			eColor color.RGBA
		)

		if usdt.percentage >= 0 {
			eColor = global.GREEN
		} else {
			eColor = global.RED
		}
		rowsUSDT = append(rowsUSDT, giu.TableRow(
			giu.Selectable(usdt.asset).Selected(false).OnClick(func() {
				AccountInstance.One.Asset = usdt.asset
				AccountInstance.Two.Asset = "USDT"
				global.FreshC <- usdt.asset + "USDT"
			}),
			giu.Label(fmt.Sprintf("%.2fm", usdt.turnOver)),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, eColor).
				To(
					giu.Label(fmt.Sprintf("%.2f%%", usdt.percentage)),
				),
		))
	}

	for i := range sortBUSD {
		var (
			busd   = mapBUSD[sortBUSD[i]]
			eColor color.RGBA
		)

		if busd.percentage >= 0 {
			eColor = global.GREEN
		} else {
			eColor = global.RED
		}
		rowsBUSD = append(rowsBUSD, giu.TableRow(
			giu.Selectable(busd.asset).Selected(false).OnClick(func() {
				AccountInstance.One.Asset = busd.asset
				AccountInstance.Two.Asset = "BUSD"
				global.FreshC <- busd.asset + "BUSD"
			}),
			giu.Label(fmt.Sprintf("%.2fm", busd.turnOver)),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, eColor).
				To(
					giu.Label(fmt.Sprintf("%.2f%%", busd.percentage)),
				),
		))
	}
	openGlobalTickerC = false
	return rowsBTC, rowsUSDT, rowsBUSD
}

func mapToSlice() ([]*myTicker, []*myTicker, []*myTicker) {
	var (
		btcTickers  []*myTicker
		usdtTickers []*myTicker
		busdTickers []*myTicker
	)
	for _, v := range mapBTC {
		btcTickers = append(btcTickers, v)
	}

	for _, v := range mapUSDT {
		usdtTickers = append(usdtTickers, v)
	}

	for _, v := range mapBUSD {
		busdTickers = append(busdTickers, v)
	}

	return btcTickers, usdtTickers, busdTickers
}

func sortByPercentage() ([]string, []string, []string) {
	btcTickers, usdtTickers, busdTickers := mapToSlice()
	return sortOneByPercentage(btcTickers), sortOneByPercentage(usdtTickers), sortOneByPercentage(busdTickers)
}

func sortByTurnOver() ([]string, []string, []string) {
	btcTickers, usdtTickers, busdTickers := mapToSlice()
	return sortOneByTurnOver(btcTickers), sortOneByTurnOver(usdtTickers), sortOneByTurnOver(busdTickers)
}

func sortByAsset() ([]string, []string, []string) {
	btcTickers, usdtTickers, busdTickers := mapToSlice()
	return sortOneByAsset(btcTickers), sortOneByAsset(usdtTickers), sortOneByAsset(busdTickers)
}

func sortOneByPercentage(tickers []*myTicker) []string {
	sort.Slice(tickers, func(i, j int) bool {
		return tickers[i].percentage > tickers[j].percentage
	})

	var asset []string
	for i := range tickers {
		asset = append(asset, tickers[i].asset)
	}

	return asset
}

func sortOneByTurnOver(tickers []*myTicker) []string {
	sort.Slice(tickers, func(i, j int) bool {
		return tickers[i].turnOver > tickers[j].turnOver
	})

	var asset []string
	for i := range tickers {
		asset = append(asset, tickers[i].asset)
	}

	return asset
}

func sortOneByAsset(tickers []*myTicker) []string {
	sort.SliceStable(tickers, func(i, j int) bool {
		var res = strings.Compare(tickers[i].asset, tickers[j].asset)
		if res == -1 {
			return false
		} else {
			return true
		}
	})

	var asset []string
	for i := range tickers {
		asset = append(asset, tickers[i].asset)
	}

	return asset
}
