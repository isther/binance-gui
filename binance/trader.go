package binance

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/hotkey"
	"github.com/isther/binanceGui/orderlist"
	"github.com/isther/binanceGui/utils"

	libBinance "github.com/adshao/go-binance/v2"
)

// var Orders = make(map[byte][]string)

type Trader struct {
	mode          giu.Modifier
	clientOrderID string

	// key      byte
	key      string
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

	switch key {
	case ";":
		key = "Semicolon"
	case ",":
		key = "Comma"
	case ".":
		key = "Period"
	case "/":
		key = "Slash"
	}

	return &Trader{
		mode:          mode,
		clientOrderID: fmt.Sprintf("%d%s", time.Now().UnixNano(), key),

		key:      key,
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

/*
全仓购买:
	模式一:
		price: 价格均为按键价格加减一步进
		quantity: 当前余额可买卖的所有数量
	模式二:
	 	price: 前五档按照按键价格加减一步进, 六到二十当前市价*对应波动比 -- global.VolatilityRatiosBuy global.VolatilityRatiosSale
		quantity: 当前余额可买卖的所有数量
*/
func (t *Trader) createOrderOnFullWarehouse() {
	var (
		priceStr    string
		quantityStr string
	)

	if t.sideType == libBinance.SideTypeBuy {
		//{{{ Prepare price:
		var (
			price float64
		)
		if global.TradeMode == global.AllPlusOneSize || t.tableID < 5 {
			//模式一: 按键价格+1步进
			price, _ = strconv.ParseFloat(wsDepthTable.Bids[t.tableID].Price, 64)
			price = t.pricePlusTickSize(price)
		} else {
			//模式二: 前五档按照按键价格加一步进, 六到二十当前市价*对应波动比 -- global.VolatilityRatiosBuy global.VolatilityRatiosSale
			price, _ = strconv.ParseFloat(AggTradePrice, 64)
			price = price * float64(global.VolatilityRatiosBuy[t.tableID])
		}
		priceStr = t.priceCorrection(price)
		//}}}

		//{{{ Prepare quantity: 当前余额/按键价格=数量
		free, _ := strconv.ParseFloat(AccountInstance.Two.Free, 64)
		quantityStr = t.quantityCorrection(float64(free / price * 0.99))
		//}}}
		console.ConsoleInstance.Write(fmt.Sprintf("全仓买入"))
	} else {
		//{{{ Prepare price:
		var (
			price float64
		)
		if global.TradeMode == global.AllPlusOneSize || t.tableID < 5 {
			//模式一: 按键价格-1步进
			price, _ = strconv.ParseFloat(wsDepthTable.Asks[t.tableID].Price, 64)
			price = t.priceSubTickSize(price)
		} else {
			price, _ = strconv.ParseFloat(AggTradePrice, 64)
			price = price * float64(global.VolatilityRatiosSale[t.tableID])
		}
		priceStr = t.priceCorrection(price)

		//}}}

		//{{{ Prepare quantity: 当前持有全部卖出
		quantityStr = AccountInstance.One.Free

		// Check Balance: 检查余额是否为零
		quantity, _ := strconv.ParseFloat(quantityStr, 64)
		quantityStr = t.quantityCorrection(quantity)
		if reflect.DeepEqual(quantity, 0.0) {
			console.ConsoleInstance.Write(fmt.Sprintf("已全仓卖出, 无需再次操作"))
			return
		}
		//}}}
	}

	t.createOrder(priceStr, quantityStr)
}

/*
分仓购买:
	模式一:
		price: 价格均为按键价格加减一步进
		quantity: 分仓后的固定数量 global.AverageSymbol1Amount
	模式二:
		price: 前五档按照按键价格加减一步进, 六到二十当前市价*对应波动比 global.VolatilityRatiosBuy global.VolatilityRatiosSale
		quantity: 分仓后的固定数量 global.AverageSymbol1Amount
*/
func (t *Trader) createOrderOnSubWarehouse() {
	var (
		priceStr    string
		quantityStr string
	)

	if t.sideType == libBinance.SideTypeBuy {
		//{{{ Prepare price:
		var (
			price float64
		)
		if global.TradeMode == global.AllPlusOneSize || t.tableID < 5 {
			//模式一: 按键价格加一步进
			price, _ = strconv.ParseFloat(wsDepthTable.Bids[t.tableID].Price, 64)
			price = t.pricePlusTickSize(price)
		} else {
			//模式二: 前五档按照按键价格加一步进, 六到二十当前市价*对应波动比 -- global.VolatilityRatiosBuy global.VolatilityRatiosSale
			price, _ = strconv.ParseFloat(AggTradePrice, 64)
			price = price * float64(global.VolatilityRatiosBuy[t.tableID])
		}
		priceStr = t.priceCorrection(price)
		//}}}

		//{{{ Prepare quantity: 分仓的固定数量 global.AverageSymbol1Amount
		quantityStr = t.quantityCorrection(global.AverageSymbol1Amount)
		//}}}

		console.ConsoleInstance.Write(fmt.Sprintf("分仓买入"))
	} else {
		//{{{ Prepare price:
		var (
			price float64
		)
		if global.TradeMode == global.AllPlusOneSize || t.tableID < 5 {
			//模式一: 按键价格减一步进
			price, _ = strconv.ParseFloat(wsDepthTable.Asks[t.tableID].Price, 64)
			price = t.priceSubTickSize(price)
		} else {
			//模式二: 前五档按照按键价格加一步进, 六到二十当前市价*对应波动比 -- global.VolatilityRatiosBuy global.VolatilityRatiosSale
			price, _ = strconv.ParseFloat(AggTradePrice, 64)
			price = price * float64(global.VolatilityRatiosSale[t.tableID])
		}
		priceStr = t.priceCorrection(price)

		//}}}

		// {{{ Prepare quantity: 分仓的固定数量 global.AverageSymbol1Amount
		quantityStr = t.quantityCorrection(global.AverageSymbol1Amount)
		//}}}

		// {{{ Check Balance: 检查余额是否充足
		// quantityAll, _ := strconv.ParseFloat(AccountInstance.One.Free, 64)
		// if !float64CompareSmallerOrEqual(quantity, quantityAll, AccountInstance.LotSizeFilter.stepSize) {
		// 	console.ConsoleInstance.Write(fmt.Sprintf("余额不足, 持仓数量: %v下单数量: %v", quantityAll, quantity))
		// 	return
		// }
		//}}}

		console.ConsoleInstance.Write(fmt.Sprintf("分仓卖出"))
	}

	t.createOrder(priceStr, quantityStr)
}

func (t *Trader) createOrder(price, quantity string) {
	_, err := GetClient().NewCreateOrderService().Symbol(AccountInstance.Symbol).
		Side(t.sideType).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).NewClientOrderID(t.clientOrderID).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v price: %v quantity: %v", err, price, quantity))
		return
	}
}

func (t *Trader) cancelOrder() {
	orders := orderlist.OrderListInstance.GetOrders(t.key)
	if len(orders) == 0 {
		console.ConsoleInstance.Write("No order")
		return
	}

	for i := range orders {
		go t.cancelAOrder(orders[i].ClientOrderID)
	}
}

func (t *Trader) cancelAOrder(clientOrderID string) {
	_, err := GetClient().NewCancelOrderService().
		Symbol(AccountInstance.Symbol).OrigClientOrderID(clientOrderID).
		Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
}

func (t *Trader) pricePlusTickSize(price float64) float64 {
	tickSize, err := strconv.ParseFloat(AccountInstance.PriceFilter.tickSize, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Price plus tickSize error: %v", err))
		return tickSize
	}

	return price + tickSize
}

func (t *Trader) priceSubTickSize(price float64) float64 {
	tickSize, err := strconv.ParseFloat(AccountInstance.PriceFilter.tickSize, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Price plus tickSize error: %v", err))
		return tickSize
	}

	return price - tickSize
}

func (t *Trader) priceCorrection(price float64) string {
	return correction(price, AccountInstance.PriceFilter.tickSize)
}

func (t *Trader) quantityCorrection(quantity float64) string {
	return correction(quantity, AccountInstance.LotSizeFilter.stepSize)
}

func correction(val float64, size string) string {
	var (
		oneIdx    = strings.Index(size, "1")
		pointIdx  = strings.Index(size, ".")
		precision int
		resStr    string
	)
	if oneIdx < pointIdx {
		precision = oneIdx - pointIdx + 1
		resStr = fmt.Sprintf("%.8f", utils.RoundLower(val, precision))
		pointIdx = strings.Index(resStr, ".")
		resStr = resStr[:pointIdx]
		return resStr
	}

	precision = oneIdx - pointIdx
	resStr = fmt.Sprintf("%.8f", utils.RoundLower(val, precision))

	return resStr
}

func float64CompareSmallerOrEqual(smaller, greater float64, accuracyStr string) bool {
	accuracy, err := strconv.ParseFloat(accuracyStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Compare float64 failed: %v", err))
		return false
	}
	return math.Max(smaller, greater) == greater || math.Abs(greater-smaller) < accuracy
}
