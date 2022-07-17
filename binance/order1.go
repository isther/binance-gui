package binance

import (
	"context"
	"fmt"
	"sort"
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
		rows     []*giu.TableRowWidget
		strSlice []string
		countSet = make(map[string]float64)
	)

	rows = append(rows, giu.TableRow(
		giu.Label("标号"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE))

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
			return false
		} else {
			return true
		}
	})

	var length = len(strSlice)
	for i := range strSlice {
		v := countSet[strSlice[i]]
		rows = append(rows, giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				To(
					giu.Label(fmt.Sprintf("%d", length-i)),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				// SetColor(giu.StyleColorText, orderColorSet()).
				To(
					giu.Label(priceFloat648Point(fmt.Sprintf("%s", strSlice[i]))),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				SetColor(giu.StyleColorText, global.RED).
				To(
					giu.Label(fmt.Sprintf("%.2fK", v)),
				),
		))
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
		giu.Label("标号"),
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
		v := countSet[strSlice[i]]
		rows = append(rows, giu.TableRow(
			giu.Style().
				SetFontSize(global.Order2FontSize).
				To(
					giu.Label(fmt.Sprintf("%d", i+1)),
				),
			giu.Style().
				SetFontSize(global.Order2FontSize).
				// SetColor(giu.StyleColorText, orderColorSet()).
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
