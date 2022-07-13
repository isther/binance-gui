package binance

import (
	"context"

	"github.com/isther/binanceGui/global"
)

func UpdateAccount() {
	res, err := GetClient().NewGetAccountService().Do(context.Background())
	if err != nil {
		panic(err)
	}
	for i := range res.Balances {
		if res.Balances[i].Asset == global.Symbol1 {
			global.Symbol1Free = res.Balances[i].Free
			global.Symbol1Locked = res.Balances[i].Locked
		}

		if res.Balances[i].Asset == global.Symbol2 {
			global.Symbol2Free = res.Balances[i].Free
			global.Symbol2Locked = res.Balances[i].Locked
		}
	}
}
