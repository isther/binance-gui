package binance

import (
	"context"
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/utils"

	libBinance "github.com/adshao/go-binance/v2"
)

var (
	depthBuyTable  []*giu.TableRowWidget
	depthSaleTable []*giu.TableRowWidget
)

func GetHttpDepthBuyTable() []*giu.TableRowWidget {
	return depthBuyTable
}

func GetHttpDepthSaleTable() []*giu.TableRowWidget {
	return depthSaleTable
}

func StartHttpDepthTable() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			refreshOrder1()
			<-ticker.C
		}
	}()
}

func refreshOrder1() {
	var res = getDepth()
	depthBuyTable = buildHttpDepthBuyTable(res)
	depthSaleTable = buildHttpDepthSaleTable(res)
}

func getDepth() *libBinance.DepthResponse {
	res, err := GetClient().NewDepthService().Symbol(AccountInstance.Symbol).Limit(100).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return nil
	}

	return res
}

func buildHttpDepthBuyTable(res *libBinance.DepthResponse) []*giu.TableRowWidget {
	var (
		rows     = make([]*giu.TableRowWidget, 21)
		strSlice []string
		countSet = make(map[string]float64)
	)

	rows[0] = giu.TableRow(
		giu.Label("涨跌幅度"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE)

	for i := 1; i < 21; i++ {
		rows[i] = giu.TableRow(
			giu.Label(""),
			giu.Label(""),
			giu.Label(""),
		)
	}

	if res == nil {
		return rows
	}

	for i := range res.Asks {
		price, quantity, err := res.Asks[i].Parse()
		if err != nil {
			console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		}
		priceStr := utils.Float64ToStringLen3(price)
		countSet[priceStr] += price * quantity / 1000
	}

	for key := range countSet {
		strSlice = append(strSlice, key)
	}

	sort.SliceStable(strSlice, func(i, j int) bool {
		var res = strings.Compare(strSlice[i], strSlice[j])
		if res == -1 {
			return true
		} else {
			return false
		}
	})

	for i := range strSlice {
		if i >= 20 {
			break
		}
		v := countSet[strSlice[i]]
		price, _ := strconv.ParseFloat(strSlice[i], 64)
		aggTradePrice, _ := strconv.ParseFloat(AggTradePrice, 64)
		rows[20-i] = giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.RED).
				To(
					giu.Label(fmt.Sprintf("%.2f%%", price/aggTradePrice*100.0-100)),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, order1ColorSet(v)).
				To(
					giu.Label(priceFloat648Point(fmt.Sprintf("%s", strSlice[i]))),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.RED).
				To(
					giu.Label(fmt.Sprintf("%.2fK", v)),
				),
		)
	}
	return rows
}

func buildHttpDepthSaleTable(res *libBinance.DepthResponse) []*giu.TableRowWidget {
	var (
		rows     []*giu.TableRowWidget
		strSlice []string
		countSet = make(map[string]float64)
	)

	rows = append(rows, giu.TableRow(
		giu.Label("涨跌幅度"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE))

	if res == nil {
		return rows
	}

	for i := range res.Asks {
		price, quantity, err := res.Bids[i].Parse()
		if err != nil {
			console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		}
		priceStr := utils.Float64ToStringLen3(price)
		countSet[priceStr] += price * quantity / 1000
	}

	for key := range countSet {
		strSlice = append(strSlice, key)
	}

	sort.SliceStable(strSlice, func(i, j int) bool {
		var res = strings.Compare(strSlice[i], strSlice[j])
		if res == -1 {
			return false
		} else {
			return true
		}
	})

	for i := range strSlice {
		if i > 20 {
			break
		}
		v := countSet[strSlice[i]]
		price, _ := strconv.ParseFloat(strSlice[i], 64)
		aggTradePrice, _ := strconv.ParseFloat(AggTradePrice, 64)
		rows = append(rows, giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.GREEN).
				To(
					giu.Label(fmt.Sprintf("%.2f%%", price/aggTradePrice*100.0-100)),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, order1ColorSet(v)).
				To(
					giu.Label(priceFloat648Point(fmt.Sprintf("%s", strSlice[i]))),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.GREEN).
				To(
					giu.Label(fmt.Sprintf("%.2fK", v)),
				),
		))
	}
	return rows
}

func order1ColorSet(ff float64) color.RGBA {
	var (
		priceColor = global.WHITE
	)
	ff *= 1000

	if ff >= float64(global.Order1BigOrderReminder[4]) {
		priceColor = global.BLACK
	} else if ff >= float64(global.Order1BigOrderReminder[3]) {
		priceColor = global.BLUE2
	} else if ff >= float64(global.Order1BigOrderReminder[2]) {
		priceColor = global.YELLOW
	} else if ff >= float64(global.Order1BigOrderReminder[1]) {
		priceColor = global.RED
	} else {
		priceColor = global.WHITE
	}

	return priceColor
}
