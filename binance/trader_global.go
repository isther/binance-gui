package binance

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
)

/*
F1: 分仓买，实时价格*波动比
F2: 全仓买，实时价格*波动比
F5: 分仓卖，实时价格*波动比
F6: 全仓卖，实时价格*波动比

F4: 撤掉所有买单
F8: 撤掉所有卖单
F9: 撤掉所有委托
F12: 撤掉所有委托，出售所有持仓
*/
type GlobalTrader struct {
	key           string
	clientOrderID string
	sideType      libBinance.SideType
}

func NewGlobalTrader(key string) *GlobalTrader {
	var sideType libBinance.SideType
	switch key {
	case "F1":
		sideType = libBinance.SideTypeBuy
	case "F2":
		sideType = libBinance.SideTypeBuy
	case "F5":
		sideType = libBinance.SideTypeSell
	case "F6":
		sideType = libBinance.SideTypeSell
	case "F12":
		sideType = libBinance.SideTypeSell
	}

	return &GlobalTrader{
		key:           key,
		clientOrderID: fmt.Sprintf("G%d", time.Now().UnixNano()),
		sideType:      sideType,
	}
}

func (g *GlobalTrader) Trade() {
	switch g.key {
	case "F1":
		// Create order on Sub Warehouse
		g.createOrderOnSubWarehouse()
	case "F2":
		// Create order on Full Warehouse
		g.createOrderOnFullWarehouse()
	case "F5":
		// Create order on Sub Warehouse
		g.createOrderOnSubWarehouse()
	case "F6":
		// Create order on Full Warehouse
		g.createOrderOnFullWarehouse()
	case "F4":
		// Cancel all buy orders
		g.cancelOrderOnBuy()
	case "F8":
		// Cancel all sell orders
		g.cancelOrderOnSale()
	case "F9":
		// Cancel all orders
		g.cancelAllOrder()
	case "F12":
		// Cancel all orders, sell all positions
		g.cancelAllOrderAndSellAllPositions()
	}
}

/*
分仓买，实时价格*波动比
分仓卖，实时价格*波动比
*/
func (g *GlobalTrader) createOrderOnSubWarehouse() {
	var (
		priceStr    string
		quantityStr string
	)

	price, _ := strconv.ParseFloat(AggTradePrice, 64)

	if g.sideType == libBinance.SideTypeBuy {
		price = price * float64(global.VolatilityRatiosF1)
		console.ConsoleInstance.Write(fmt.Sprintf("全局分仓买入"))
	} else {
		price = price * float64(global.VolatilityRatiosF5)
		console.ConsoleInstance.Write(fmt.Sprintf("全局分仓卖出"))
	}

	price = g.priceCorrection(price)
	priceStr = fmt.Sprintf("%f", price)

	quantity := g.quantityCorrection(global.AverageSymbol1Amount)
	quantityStr = fmt.Sprintf("%f", quantity)

	g.createOrder(priceStr, quantityStr)
}

/*
全仓买，实时价格*波动比
全仓卖，实时价格*波动比
*/
func (g *GlobalTrader) createOrderOnFullWarehouse() {
	var (
		priceStr    string
		quantityStr string
	)

	price, _ := strconv.ParseFloat(AggTradePrice, 64)
	if g.sideType == libBinance.SideTypeBuy {
		price = price * float64(global.VolatilityRatiosF2)
		price = g.priceCorrection(price)
		priceStr = fmt.Sprintf("%f", price)

		free, _ := strconv.ParseFloat(AccountInstance.Two.Free, 64)
		quantity := g.quantityCorrection(float64(free / price))
		quantityStr = fmt.Sprintf("%f", quantity)
		console.ConsoleInstance.Write(fmt.Sprintf("全局全仓买入"))
	} else {
		if g.key == "F12" {
			price = price * float64(global.VolatilityRatiosF12)
			console.ConsoleInstance.Write(fmt.Sprintf("全局清仓卖出"))
		} else {
			price = price * float64(global.VolatilityRatiosF6)
			console.ConsoleInstance.Write(fmt.Sprintf("全局全仓卖出"))
		}
		price = g.priceCorrection(price)
		priceStr = fmt.Sprintf("%f", price)

		quantityStr = AccountInstance.One.Free
		// Check Balance: 检查余额是否为零
		quantity, _ := strconv.ParseFloat(quantityStr, 64)
		quantity = g.quantityCorrection(quantity)
		if reflect.DeepEqual(quantity, 0.0) {
			console.ConsoleInstance.Write(fmt.Sprintf("已全仓卖出, 无需再次操作"))
			return
		}
		quantityStr = fmt.Sprintf("%f", quantity)
	}

	g.createOrder(priceStr, quantityStr)
}

func (g *GlobalTrader) priceCorrection(price float64) float64 {
	priceStr := correction(fmt.Sprintf("%.8f", price), AccountInstance.PriceFilter.tickSize)
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Price correction error: %v", err))
		return 0
	}
	return price
}

func (g *GlobalTrader) quantityCorrection(quantity float64) float64 {
	quantityStr := correction(fmt.Sprintf("%.8f", quantity), AccountInstance.LotSizeFilter.stepSize)
	quantity, err := strconv.ParseFloat(quantityStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Quantity correction error: %v", err))
		return 0
	}
	return quantity
}

// 撤掉所有买单
func (g *GlobalTrader) cancelOrderOnBuy() {
	for _, order := range g.getAllOpenOrders() {
		if order.Side == libBinance.SideTypeBuy {
			go g.cancelAOrder(order.ClientOrderID)
		}
	}
}

// 撤掉所有卖单
func (g *GlobalTrader) cancelOrderOnSale() {
	for _, order := range g.getAllOpenOrders() {
		if order.Side == libBinance.SideTypeSell {
			go g.cancelAOrder(order.ClientOrderID)
		}
	}
}

func (g *GlobalTrader) getAllOpenOrders() []*libBinance.Order {
	res, err := GetClient().NewListOpenOrdersService().Symbol(AccountInstance.Symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	return res
}

// 撤掉所有委托
func (g *GlobalTrader) cancelAllOrder() {
	_, err := GetClient().NewCancelOpenOrdersService().Symbol(AccountInstance.Symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
}

// 出售所有持仓
func (g *GlobalTrader) cancelAllOrderAndSellAllPositions() {
	g.cancelAllOrder()
	g.createOrderOnFullWarehouse()
}

func (g *GlobalTrader) createOrder(price, quantity string) {
	_, err := GetClient().NewCreateOrderService().Symbol(AccountInstance.Symbol).
		Side(g.sideType).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v price: %v quantity: %v", err, price, quantity))
		return
	}
}

func (g *GlobalTrader) cancelAOrder(clientOrderID string) {
	_, err := GetClient().NewCancelOrderService().
		Symbol(AccountInstance.Symbol).OrigClientOrderID(clientOrderID).
		Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
}