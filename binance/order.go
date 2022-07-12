package binance

import (
	"context"
	"fmt"

	"github.com/adshao/go-binance/v2"
)

func Order() {
	order, err := GetClient().NewCreateOrderService().Symbol("BNBETH").
		Side(binance.SideTypeBuy).Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity("5").
		Price("0.0030000").Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(order)
}
