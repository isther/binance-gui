package binance

import (
	"context"
	"fmt"
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

func priceFloat648Point(f float64) string {
	s := fmt.Sprintf("%.8f", f)
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

func UpdateAverageAmount() {
	free1, _ := strconv.ParseFloat(AccountInstance.One.Free, 64)
	free2, _ := strconv.ParseFloat(AccountInstance.Two.Free, 64)
	global.AverageSymbol1Amount = free1 / float64(global.Average)
	global.AverageSymbol2Amount = free2 / float64(global.Average)
}
