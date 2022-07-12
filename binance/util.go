package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
	"github.com/isther/binance/conf"
)

func Status() {
	client := GetClient()
	res, err := client.NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func symbolExist(symbol string) bool {
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
