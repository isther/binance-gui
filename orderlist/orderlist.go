package orderlist

import (
	"fmt"
	"sync"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/hotkey"
)

var (
	OrderListInstance *OrderList
)

func init() {
	OrderListInstance = NewOrderList()
}

type OrderList struct {
	BuyOrders  map[string][]*libBinance.Order
	SaleOrders map[string][]*libBinance.Order
	mu         sync.Mutex
}

func NewOrderList() *OrderList {
	return &OrderList{
		BuyOrders:  make(map[string][]*libBinance.Order),
		SaleOrders: make(map[string][]*libBinance.Order),
	}
}

func (orderList *OrderList) Lock() {
	orderList.mu.Lock()
}

func (orderList *OrderList) Unlock() {
	orderList.mu.Unlock()
}

// Add order to order list
func (orderList *OrderList) AddOrders(order *libBinance.Order) {
	key := parseOrderIDToKey(order.ClientOrderID)
	if order.Side == libBinance.SideTypeBuy {
		orderList.checkBuyOrderIsExistIfNotToMake(key)
		orderList.mu.Lock()
		defer orderList.mu.Unlock()

		orderList.BuyOrders[key] = append(orderList.BuyOrders[key], order)
	} else if order.Side == libBinance.SideTypeSell {
		orderList.checkSaleOrderIsExistIfNotToMake(key)
		orderList.mu.Lock()
		defer orderList.mu.Unlock()

		orderList.SaleOrders[key] = append(orderList.SaleOrders[key], order)
	}
}

// Add orders by key
func (orderList *OrderList) GetOrders(key string) []*libBinance.Order {
	orderList.mu.Lock()
	defer orderList.mu.Unlock()

	var (
		sideType = orderList.getKeyType(key)
	)

	key = keyToSymbol(key)
	if sideType == libBinance.SideTypeBuy {
		return orderList.BuyOrders[key]
	} else if sideType == libBinance.SideTypeSell {
		return orderList.SaleOrders[key]
	}
	return nil
}

// Remove orders by key
func (orderList *OrderList) CancelOrdersByKey(key string) {
	orderList.mu.Lock()
	defer orderList.mu.Unlock()

	var (
		sideType = orderList.getKeyType(key)
	)

	if sideType == libBinance.SideTypeBuy {
		orderList.BuyOrders[key] = make([]*libBinance.Order, 0)
	} else if sideType == libBinance.SideTypeSell {
		orderList.SaleOrders[key] = make([]*libBinance.Order, 0)
	}
}

// Remove orders by order's id
func (orderList *OrderList) CancelOrdersByID(orderID ...string) {
	for _, id := range orderID {
		orderList.cancelOrderByID(id)
	}
}

func (orderList *OrderList) Clear() {
	orderList.mu.Lock()
	defer orderList.mu.Unlock()

	orderList.BuyOrders = make(map[string][]*libBinance.Order)
	orderList.SaleOrders = make(map[string][]*libBinance.Order)
}

func (orderList *OrderList) cancelOrderByID(orderID string) {
	orderList.mu.Lock()
	defer orderList.mu.Unlock()

	var (
		key      = parseOrderIDToKey(orderID)
		sideType = orderList.getKeyType(key)
	)

	if sideType == libBinance.SideTypeBuy {
		orders := orderList.BuyOrders[key]
		orderList.BuyOrders[key] = orderList.deleteAnOrderInList(orders, orderID)
	} else if sideType == libBinance.SideTypeSell {
		orders := orderList.SaleOrders[key]
		orderList.SaleOrders[key] = orderList.deleteAnOrderInList(orders, orderID)
	}
}

func (orderList *OrderList) deleteAnOrderInList(orders []*libBinance.Order, orderID string) []*libBinance.Order {
	newOrders := make([]*libBinance.Order, 0)
	for _, order := range orders {
		if order.ClientOrderID != orderID {
			newOrders = append(newOrders, order)
		}
	}
	return newOrders
}

func (orderList *OrderList) OutPutBuyOrders() []*libBinance.Order {
	OrderListInstance.Lock()
	defer OrderListInstance.Unlock()

	var (
		buyOrders []*libBinance.Order
	)
	for _, order := range orderList.BuyOrders {
		buyOrders = append(buyOrders, order...)
	}
	return buyOrders
}

func (orderList *OrderList) OutPutSaleOrders() []*libBinance.Order {
	OrderListInstance.Lock()
	defer OrderListInstance.Unlock()

	var (
		saleOrders []*libBinance.Order
	)
	for _, order := range orderList.SaleOrders {
		saleOrders = append(saleOrders, order...)
	}
	return saleOrders
}

func parseOrderIDToKey(clientOrderID string) string {
	return keyToSymbol(clientOrderID[19:])
}

func keyToSymbol(key string) string {
	switch key {
	case "Semicolon":
		return ";"
	case "Comma":
		return ","
	case "Period":
		return "."
	case "Slash":
		return "/"
	default:
		return key
	}
}

func (orderList *OrderList) getKeyType(key string) libBinance.SideType {
	var (
		sideType libBinance.SideType
	)

	if len(key) > 1 {
		switch key {
		case "F1", "F2":
			sideType = libBinance.SideTypeBuy
		case "F5", "F6":
			sideType = libBinance.SideTypeSell
		case "F12":
			sideType = libBinance.SideTypeSell
		case "Semicolon", "Comma", "Period", "Slash":
			sideType = libBinance.SideTypeSell
		}
	} else {
		sideType = hotkey.GetHotKeyType(key[0])
	}
	return sideType
}

func (orderList *OrderList) checkBuyOrderIsExistIfNotToMake(key string) {
	orderList.mu.Lock()
	defer orderList.mu.Unlock()

	_, ok := orderList.BuyOrders[key]
	if !ok {
		orderList.BuyOrders[key] = make([]*libBinance.Order, 0)
	}
}

func (orderList *OrderList) checkSaleOrderIsExistIfNotToMake(key string) {
	orderList.mu.Lock()
	defer orderList.mu.Unlock()

	_, ok := orderList.SaleOrders[key]
	if !ok {
		orderList.SaleOrders[key] = make([]*libBinance.Order, 0)
	}
}

func (orderList *OrderList) print() {
	fmt.Println("Buy orders: ")
	for k, v := range orderList.BuyOrders {
		fmt.Println(k, ": ")
		for _, vv := range v {
			fmt.Println(vv)
		}
	}
	fmt.Println("Sale orders: ")
	for k, v := range orderList.SaleOrders {
		fmt.Println(k, ": ")
		for _, vv := range v {
			fmt.Println(vv)
		}
	}
}
