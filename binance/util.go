package binance

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/conf"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

func Status() {
	client := GetClient()
	res, err := client.NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	fmt.Println(res)
}

func SymbolExist(symbol string) bool {
	_, err := binance.NewClient(conf.Conf.ApiKey, conf.Conf.SecretKey).NewExchangeInfoService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return false
	}

	return true
}

func priceFloat648Point(s string) string {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == '.' {
			return string(append([]byte(s[:]), '0'))
		}
		if s[i] != '0' {
			return s
		}
		s = s[:i]
	}
	return s
}

// 当前余额/分仓数/当前实时价格
func UpdateAverageAmount() {
	free, _ := strconv.ParseFloat(AccountInstance.Two.Free, 64)
	if reflect.DeepEqual(free, 0.0) {
		console.ConsoleInstance.Write(fmt.Sprintf("余额不足，无法分仓"))
	}

	res, err := GetClient().NewListPricesService().Symbol(AccountInstance.Symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	console.ConsoleInstance.Write("分仓中...")
	aggTradePrice, _ := strconv.ParseFloat(res[0].Price, 64)
	averageSymbol1AmountStr := correction((free/float64(global.Average))/aggTradePrice, AccountInstance.LotSizeFilter.stepSize)
	global.AverageSymbol1Amount, _ = strconv.ParseFloat(averageSymbol1AmountStr, 64)

	averageSymbol2AmountStr := correction(free/float64(global.Average), AccountInstance.PriceFilter.tickSize)
	global.AverageSymbol2Amount, _ = strconv.ParseFloat(averageSymbol2AmountStr, 64)
	console.ConsoleInstance.Write(fmt.Sprintf("已分仓, 单仓购买数量: %v", global.AverageSymbol1Amount))
}
