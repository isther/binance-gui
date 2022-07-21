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
		clientOrderID: fmt.Sprintf("%d%s", time.Now().UnixNano(), key),
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

	priceStr = g.priceCorrection(price)

	quantityStr = g.quantityCorrection(global.AverageSymbol1Amount)

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
		priceStr = g.priceCorrection(price)

		free, _ := strconv.ParseFloat(AccountInstance.Two.Free, 64)
		quantityStr = g.quantityCorrection(float64(free / price))
		console.ConsoleInstance.Write(fmt.Sprintf("全局全仓买入"))
	} else {
		if g.key == "F12" {
			price = price * float64(global.VolatilityRatiosF12)
			console.ConsoleInstance.Write(fmt.Sprintf("全局清仓卖出"))
		} else {
			price = price * float64(global.VolatilityRatiosF6)
			console.ConsoleInstance.Write(fmt.Sprintf("全局全仓卖出"))
		}
		priceStr = g.priceCorrection(price)

		quantityStr = AccountInstance.One.Free
		// Check Balance: 检查余额是否为零
		quantity, _ := strconv.ParseFloat(quantityStr, 64)
		quantityStr = g.quantityCorrection(quantity)
		quantity, _ = strconv.ParseFloat(quantityStr, 64)
		if reflect.DeepEqual(quantity, 0.0) {
			console.ConsoleInstance.Write(fmt.Sprintf("已全仓卖出, 无需再次操作"))
			return
		}
		ResetCostInstance()
	}

	g.createOrder(priceStr, quantityStr)
}

func (g *GlobalTrader) priceCorrection(price float64) string {
	return correction(price, AccountInstance.PriceFilter.tickSize)
}

func (g *GlobalTrader) quantityCorrection(quantity float64) string {
	return correction(quantity, AccountInstance.LotSizeFilter.stepSize)
}

// 撤掉所有买单
func (g *GlobalTrader) cancelOrderOnBuy() {
	for _, order := range g.getAllOpenOrders() {
		if order.Side == libBinance.SideTypeBuy {
			go g.cancelAOrder(order.ClientOrderID)
		}
	}
	console.ConsoleInstance.Write("取消所有买单")
}

// 撤掉所有卖单
func (g *GlobalTrader) cancelOrderOnSale() {
	for _, order := range g.getAllOpenOrders() {
		if order.Side == libBinance.SideTypeSell {
			go g.cancelAOrder(order.ClientOrderID)
		}
	}
	console.ConsoleInstance.Write("取消所有卖单")
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
	console.ConsoleInstance.Write("取消所有委托")
}

// 出售所有持仓
func (g *GlobalTrader) cancelAllOrderAndSellAllPositions() {
	g.cancelAllOrder()

	console.ConsoleInstance.Write("清仓")
	for i := 0; i < 4; i++ {
		g.createOrderOnFullWarehouse()
		time.Sleep(200 * time.Millisecond)
	}
}

func (g *GlobalTrader) createOrder(price, quantity string) {
	console.ConsoleInstance.Write(fmt.Sprintf("%s %s", price, quantity))
	_, err := GetClient().NewCreateOrderService().Symbol(AccountInstance.Symbol).
		Side(g.sideType).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).NewClientOrderID(g.clientOrderID).Do(context.Background())
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
