package binance

import (
	"context"
	"fmt"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"

	libBinance "github.com/adshao/go-binance/v2"
)

type Trader struct {
	Mode    giu.Modifier
	Key     string
	OrderID string
}

func NewTrader(mode giu.Modifier, key string) *Trader {
	return &Trader{
		Mode:    mode,
		Key:     key,
		OrderID: fmt.Sprintf("%d%d", key[0], time.Now().UnixNano()),
	}
}

func (t *Trader) Trade() {
	switch t.Mode {
	case giu.ModNone:
		// Create order on Full Warehouse
		t.createOrderOnFullWarehouse()
	case giu.ModAlt:
		// Create order on Sub Warehouse
		t.createOrderOnSubWarehouse()
	case giu.ModShift:
		// Cancel order
		t.cancelOrder()
	}
}

func (t *Trader) createOrderOnFullWarehouse() {
	console.ConsoleInstance.Write(fmt.Sprintf("Order Created On Full, Order ID: %s", t.OrderID))
}

func (t *Trader) createOrderOnSubWarehouse() {
	console.ConsoleInstance.Write(fmt.Sprintf("Order Created On Sub, Order ID: %s", t.OrderID))
}

func (t *Trader) createOrder(price, quantity string) {
	order, err := GetClient().NewCreateOrderService().Symbol(global.Symbol).
		Side(libBinance.SideTypeBuy).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).NewClientOrderID(t.OrderID).Do(context.Background())

	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	console.ConsoleInstance.Write(fmt.Sprintf("Order Created, Order ID: %s", order.ClientOrderID))
}

func (t *Trader) cancelOrder() {
	console.ConsoleInstance.Write(fmt.Sprintf("Order Cancelled, Order ID: %s", t.OrderID))
}
