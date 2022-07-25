package binance

import (
	"context"
	"fmt"
	"strconv"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
)

type TradeBUSDAndUSDT struct {
	symbol string

	Free string

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

func NewTradeBUSDAndUSDT() *TradeBUSDAndUSDT { return &TradeBUSDAndUSDT{symbol: "BUSDUSDT"} }

func (this *TradeBUSDAndUSDT) Trade() {
	this.updateAccount()
	this.getExchangeInfo()

	var (
		price    = this.getPrice()
		priceStr = fmt.Sprintf("%v", this.priceCorrection(price*0.995))

		quantityStr string
	)
	free, _ := strconv.ParseFloat(this.Free, 64)
	quantityStr = fmt.Sprintf("%d", int(free/price))

	this.createOrder(priceStr, quantityStr)
}

func (this *TradeBUSDAndUSDT) createOrder(price, quantity string) {
	_, err := GetClient().NewCreateOrderService().Symbol(this.symbol).
		Side(libBinance.SideTypeSell).Type(libBinance.OrderTypeLimit).
		TimeInForce(libBinance.TimeInForceTypeGTC).Quantity(quantity).
		Price(price).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v price: %v quantity: %v", err, price, quantity))
		return
	}
}

func (this *TradeBUSDAndUSDT) priceCorrection(price float64) string {
	return correction(price, this.priceFilter.tickSize)
}

func (this *TradeBUSDAndUSDT) updateAccount() {
	res, err := GetClient().NewGetAccountService().Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	for _, balance := range res.Balances {
		if balance.Asset == "BUSD" {
			this.Free = balance.Free
		}
	}
}

func (this *TradeBUSDAndUSDT) getExchangeInfo() {
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

func (this *TradeBUSDAndUSDT) getPrice() float64 {
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
