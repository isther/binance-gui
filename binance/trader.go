package binance

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/hotkey"

	libBinance "github.com/adshao/go-binance/v2"
)

var Orders = make(map[byte][]string)

type Trader struct {
	mode          giu.Modifier
	clientOrderID string

	key      byte
	tableID  int
	sideType libBinance.SideType
}

func NewTrader(mode giu.Modifier, key string) *Trader {
	var (
		tableID  int
		sideType libBinance.SideType
	)

	if id := hotkey.GetBuyKeyIndex(key[0]); id != -1 {
		tableID = id
		sideType = libBinance.SideTypeBuy
	} else if id := hotkey.GetSaleKeyIndex(key[0]); id != -1 {
		tableID = id
		sideType = libBinance.SideTypeSell
	}
	return &Trader{
		mode:          mode,
		clientOrderID: fmt.Sprintf("%d%d", time.Now().UnixNano(), key[0]),

		key:      key[0],
		tableID:  tableID,
		sideType: sideType,
	}
}

func (t *Trader) Trade() {
	switch t.mode {
	case giu.ModNone:
		// Create order on Full Warehouse
		t.createOrderOnSubWarehouse()
	case giu.ModAlt:
		// Create order on Sub Warehouse
		t.createOrderOnFullWarehouse()
	case giu.ModShift:
		// Cancel order
		t.cancelOrder()
	default:
		console.ConsoleInstance.Write(fmt.Sprintf("Error mode"))
	}
}

func (t *Trader) createOrderOnFullWarehouse() {
	var (
		priceStr    string
		quantityStr string
	)

	if t.sideType == libBinance.SideTypeBuy {
		price, _ := strconv.ParseFloat(depthTable.Bids[t.tableID].Price, 64)
		price = t.priceCorrection(price)
		priceStr = fmt.Sprintf("%f", price)

		free, _ := strconv.ParseFloat(AccountInstance.One.Free, 64)
		quantity := free / price
		quantityStr = t.quantityCorrection(quantity)
	} else {
		price, _ := strconv.ParseFloat(depthTable.Asks[t.tableID].Price, 64)
		price = t.priceCorrection(price)
		priceStr = fmt.Sprintf("%f", price)
		quantityStr = AccountInstance.One.Free
	}

	quantity, _ := strconv.ParseFloat(quantityStr, 64)
	if reflect.DeepEqual(quantity, 0.0) {
		console.ConsoleInstance.Write("账户余额不足")
		return
	}
	t.createOrder(priceStr, quantityStr)
}

func (t *Trader) createOrderOnSubWarehouse() {
	var (
		priceStr    string
		quantityStr string
	)

	if t.sideType == libBinance.SideTypeBuy {
		price, _ := strconv.ParseFloat(depthTable.Bids[t.tableID].Price, 64)
		price = t.priceCorrection(price)
		priceStr = fmt.Sprintf("%f", price)

		quantity := global.AverageSymbol2Amount / price
		quantityStr = t.quantityCorrection(quantity)
	} else {
		price, _ := strconv.ParseFloat(depthTable.Asks[t.tableID].Price, 64)
		price = t.priceCorrection(price)
		priceStr = fmt.Sprintf("%f", price)

		quantityStr = fmt.Sprintf("%f", global.AverageSymbol1Amount)
	}

	quantity, _ := strconv.ParseFloat(quantityStr, 64)
	if reflect.DeepEqual(quantity, 0.0) {
		console.ConsoleInstance.Write("账户余额不足")
		return
	}
	t.createOrder(priceStr, quantityStr)
}

func (t *Trader) createOrder(price, quantity string) {
	order, err := GetClient().NewCreateOrderService().Symbol(AccountInstance.Symbol).
		Side(t.sideType).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).NewClientOrderID(t.clientOrderID).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v price: %v quantity: %v", err, price, quantity))
		return
	}

	console.ConsoleInstance.Write(fmt.Sprintf("OK: ID: %s price: %s quantity: %s",
		order.ClientOrderID,
		price,
		quantity,
	))

	Orders[t.key] = append(Orders[t.key], t.clientOrderID)
}

func (t *Trader) cancelOrder() {
	if len(Orders[t.key]) == 0 {
		console.ConsoleInstance.Write("No order")
		return
	}

	for i := range Orders[t.key] {
		go t.cancelAOrder(Orders[t.key][i])
	}

	Orders[t.key] = []string{}
}

func (t *Trader) cancelAOrder(clientOrderID string) {
	_, err := GetClient().NewCancelOrderService().
		Symbol(AccountInstance.Symbol).OrigClientOrderID(clientOrderID).
		Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
	console.ConsoleInstance.Write(fmt.Sprintf("OK, key: %v id: %v", t.key, clientOrderID))
}

func (t *Trader) priceCorrection(price float64) float64 {
	priceStr := correction(fmt.Sprintf("%f", price), AccountInstance.PriceFilter.tickSize)

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Price correction error: %v", err))
		return 0
	}

	tickSize, err := strconv.ParseFloat(AccountInstance.PriceFilter.tickSize, 64)
	if t.sideType == libBinance.SideTypeBuy {
		return price + tickSize
	}

	return price - tickSize
}

func (t *Trader) quantityCorrection(quantity float64) string {
	return correction(fmt.Sprintf("%f", quantity), AccountInstance.LotSizeFilter.stepSize)
}

func correction(val, size string) string {
	var (
		start  bool
		length = 0
	)

	for i := len(size) - 1; i > 0; i-- {
		if size[i] == '1' {
			start = true
		}
		if start {
			length++
		}
	}

	for i := range val {
		if val[i] == '.' {
			val = val[:i+length]
			break
		}
	}
	return val
}
