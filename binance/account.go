package binance

import (
	"context"
	"fmt"

	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"

	libBinance "github.com/adshao/go-binance/v2"
)

var AccountInstance *Account

type Account struct {
	Symbol string

	One *libBinance.Balance
	Two *libBinance.Balance
	BNB *libBinance.Balance

	Balances []libBinance.Balance

	PriceFilter struct {
		minPrice string
		maxPrice string
		tickSize string
	}

	LotSizeFilter struct {
		minQty   string
		maxQty   string
		stepSize string
	}
}

func init() {
	AccountInstance = NewAccount()
}

func StartUpdateAccount() {
	AccountInstance.ExchangeInfo()
	AccountInstance.UpdateAccount()
	AccountInstance.UpdateOrderList()
}

func NewAccount() *Account {
	return &Account{
		Symbol: "BUSDUSDT",
		One:    &libBinance.Balance{Asset: "BUSD", Free: "0", Locked: "0"},
		Two:    &libBinance.Balance{Asset: "USDT", Free: "0", Locked: "0"},
		BNB:    &libBinance.Balance{Asset: "BNB", Free: "0", Locked: "0"},
	}
}

func (account *Account) UpdateAccount() {
	res, err := GetClient().NewGetAccountService().Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
	account.Balances = res.Balances
	for _, balance := range account.Balances {
		if balance.Asset == account.One.Asset {
			account.One.Free = balance.Free
			account.One.Locked = balance.Locked
		}
		if balance.Asset == account.Two.Asset {
			account.Two.Free = balance.Free
			account.Two.Locked = balance.Locked
		}

		if balance.Asset == "BNB" {
			account.BNB.Free = balance.Free
			account.BNB.Locked = balance.Locked
		}
	}
}

//获取交易信息
func (account *Account) ExchangeInfo() {
	res, err := GetClient().NewExchangeInfoService().Symbol(account.Symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	for i := range res.Symbols {
		if res.Symbols[i].Symbol == account.Symbol {
			for j := range res.Symbols[i].Filters {
				filter := res.Symbols[i].Filters[j]

				if filter["filterType"] == "PRICE_FILTER" {
					account.PriceFilter.minPrice = fmt.Sprintf("%v", filter["minPrice"])
					account.PriceFilter.maxPrice = fmt.Sprintf("%v", filter["maxPrice"])
					account.PriceFilter.tickSize = fmt.Sprintf("%v", filter["tickSize"])
				}
				if filter["filterType"] == "LOT_SIZE" {
					account.LotSizeFilter.minQty = fmt.Sprintf("%v", filter["minQty"])
					account.LotSizeFilter.maxQty = fmt.Sprintf("%v", filter["maxQty"])
					account.LotSizeFilter.stepSize = fmt.Sprintf("%v", filter["stepSize"])
				}
			}
		}
	}
}

func (account *Account) UpdateOrderList() {
	res, err := GetClient().NewListOpenOrdersService().Symbol(account.Symbol).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}

	if len(res) == 0 {
		console.ConsoleInstance.Write("No open orders")
		buildC <- struct{}{}
		return
	}

	console.ConsoleInstance.Write("载入订单...")
	for i := range res {
		order := res[i]
		console.ConsoleInstance.Write(fmt.Sprintf("Symbol: %v OrderID: %v",
			order.Symbol,
			order.ClientOrderID,
		))
		OpenOrdersInstance.AddOrders(&libBinance.Order{
			ClientOrderID: order.ClientOrderID,
			Symbol:        order.Symbol,
			Side:          libBinance.SideType(order.Side),
			Price:         order.Price,
			OrigQuantity:  order.OrigQuantity,
		})
	}
}

func (account *Account) WsUpdateAccount() (chan struct{}, chan struct{}) {
	var (
		listenKey string

		err   error
		doneC chan struct{}
		stopC chan struct{}
	)

	for {
		if listenKey, err = GetClient().NewStartUserStreamService().Do(context.Background()); err == nil {
			break
		}
	}

	wsHandler := func(event *libBinance.WsUserDataEvent) {
		account.parseAccountUpdate(event.AccountUpdate)
		account.parseBalanceUpdate(event.BalanceUpdate)
		account.parseOrderUpdate(event.OrderUpdate)
	}
	errHandler := func(err error) {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	doneC, stopC, err = libBinance.WsUserDataServe(listenKey, wsHandler, errHandler)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
	}
	return doneC, stopC
}

func (account *Account) parseAccountUpdate(accountUpdates []libBinance.WsAccountUpdate) {
	if len(accountUpdates) == 0 {
		return
	}

	for _, accountUpdate := range accountUpdates {
		if accountUpdate.Asset == account.One.Asset {
			account.One.Free = accountUpdate.Free
			account.One.Locked = accountUpdate.Locked
		}
		if accountUpdate.Asset == account.Two.Asset {
			account.Two.Free = accountUpdate.Free
			account.Two.Locked = accountUpdate.Locked
		}

		if accountUpdate.Asset == "BNB" {
			account.BNB.Free = accountUpdate.Free
			account.BNB.Locked = accountUpdate.Locked
		}
	}

	if global.Debug {
		console.ConsoleInstance.Write(fmt.Sprintf("账户余额更新: %v: %v %v %v: %v %v %v: %v %v",
			account.One.Asset, account.One.Free, account.One.Locked,
			account.Two.Asset, account.Two.Free, account.Two.Locked,
			account.BNB.Asset, account.BNB.Free, account.BNB.Locked))
	}
}

func (account *Account) parseBalanceUpdate(balanceUpdate libBinance.WsBalanceUpdate) {
	if balanceUpdate.Asset == "" {
		return
	}

	console.ConsoleInstance.Write(fmt.Sprintf("BalanceUpdate: %v %v", balanceUpdate.Asset, balanceUpdate.Change))
}

func (account *Account) parseOrderUpdate(orderUpdate libBinance.WsOrderUpdate) {
	if orderUpdate.Symbol == "" {
		return
	}

	if orderUpdate.Status == "NEW" {
		console.ConsoleInstance.Write(fmt.Sprintf("Create Order, OrderID: %v", orderUpdate.ClientOrderId))
		OpenOrdersInstance.AddOrders(&libBinance.Order{
			ClientOrderID: orderUpdate.ClientOrderId,
			Side:          libBinance.SideType(orderUpdate.Side),
			Price:         orderUpdate.Price,
			OrigQuantity:  orderUpdate.Volume,
		})
	} else if orderUpdate.Status == "CANCELED" {
		console.ConsoleInstance.Write(fmt.Sprintf("Cancel Order, OrderID: %v", orderUpdate.OrigCustomOrderId))
		OpenOrdersInstance.CancelOrders(&libBinance.Order{
			Side:          libBinance.SideType(orderUpdate.Side),
			ClientOrderID: orderUpdate.OrigCustomOrderId,
		})
	} else if orderUpdate.Status == "FILLED" {
		OpenOrdersInstance.CancelOrders(&libBinance.Order{
			Side:          libBinance.SideType(orderUpdate.Side),
			ClientOrderID: orderUpdate.ClientOrderId,
		})
	} else {
		console.ConsoleInstance.Write(fmt.Sprintf("%v", orderUpdate))
	}
}
