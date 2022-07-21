package binance

import (
	"context"
	"fmt"
	"strconv"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
)

type TradeBNB struct {
	symbol string

	priceFilter struct {
		minPrice string
		maxPrice string
		tickSize string
	}

	lotSizeFilter struct {
		minQty   string
		maxQty   string
		stepSize string
	}
}

func NewTradeBNB(asset string) *TradeBNB {
	return &TradeBNB{
		symbol: "BNB" + asset,
	}
}

func (this *TradeBNB) Trade() {
	this.getExchangeInfo()

	var (
		price       = this.getPrice() * 1.03
		priceStr    = this.priceCorrection(price)
		quantityStr = this.quantityCorrection(11 / price)
	)

	this.createOrder(priceStr, quantityStr)
}

func (this *TradeBNB) createOrder(price, quantity string) {
	_, err := GetClient().NewCreateOrderService().Symbol(this.symbol).
		Side(libBinance.SideTypeBuy).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v price: %v quantity: %v", err, price, quantity))
		return
	}
}

func (this *TradeBNB) priceCorrection(price float64) string {
	return correction(price, this.priceFilter.tickSize)
}

func (this *TradeBNB) quantityCorrection(quantity float64) string {
	return correction(quantity, this.lotSizeFilter.stepSize)
}

func (this *TradeBNB) getExchangeInfo() {
	res, err := GetClient().NewExchangeInfoService().Symbol(this.symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	for i := range res.Symbols {
		if res.Symbols[i].Symbol == this.symbol {
			for j := range res.Symbols[i].Filters {
				filter := res.Symbols[i].Filters[j]

				if filter["filterType"] == "PRICE_FILTER" {
					this.priceFilter.minPrice = fmt.Sprintf("%v", filter["minPrice"])
					this.priceFilter.maxPrice = fmt.Sprintf("%v", filter["maxPrice"])
					this.priceFilter.tickSize = fmt.Sprintf("%v", filter["tickSize"])
				}
				if filter["filterType"] == "LOT_SIZE" {
					this.lotSizeFilter.minQty = fmt.Sprintf("%v", filter["minQty"])
					this.lotSizeFilter.maxQty = fmt.Sprintf("%v", filter["maxQty"])
					this.lotSizeFilter.stepSize = fmt.Sprintf("%v", filter["stepSize"])
				}
			}
		}
	}
}

func (this *TradeBNB) getPrice() float64 {
	res, err := GetClient().NewListPricesService().Symbol(this.symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return 0.0
	}
	var priceStr = res[0].Price
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return 0.0
	}
	return price
}
