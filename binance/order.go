package binance

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/AllenDang/giu"
	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/hotkey"
)

var (
	OpenOrdersInstance *OpenOrders

	openOrdersBuyTable  []*giu.TableRowWidget
	openOrdersSaleTable []*giu.TableRowWidget

	buildC = make(chan struct{})
)

//重置订单列表
func ResetOrders() {
	OpenOrdersInstance = NewOpenOrders()
}

func StartBuildOrderTable() {
	ResetOrders()
	go func() {
		for {
			select {
			case <-buildC:
				openOrdersBuyTable = buildOpenBuyOrderTable()
				openOrdersSaleTable = buildOpenSaleOrderTable()
			}
		}
	}()
}

func GetOpenBuyOrdersTable() []*giu.TableRowWidget {
	return openOrdersBuyTable
}

func GetOpenSaleOrdersTable() []*giu.TableRowWidget {
	return openOrdersSaleTable
}

type OpenOrders struct {
	BuyOrder  []*libBinance.Order
	SaleOrder []*libBinance.Order
	lock      *sync.RWMutex
}

func NewOpenOrders() *OpenOrders {
	return &OpenOrders{
		BuyOrder:  make([]*libBinance.Order, 0),
		SaleOrder: make([]*libBinance.Order, 0),
	}
}

func (openOrders *OpenOrders) AddOrders(order *libBinance.Order) {
	if order.Side == libBinance.SideTypeBuy {
		openOrders.BuyOrder = append(openOrders.BuyOrder, order)
	} else {
		openOrders.SaleOrder = append(openOrders.SaleOrder, order)
	}
	buildC <- struct{}{}
}

func (openOrders *OpenOrders) GetOrders(key byte) []*libBinance.Order {
	var (
		orders []*libBinance.Order
	)
	sideType := hotkey.GetHotKeyType(key)

	if sideType == libBinance.SideTypeBuy {
		for _, openOrder := range openOrders.BuyOrder {
			if parseOrderID(openOrder.ClientOrderID) == key {
				orders = append(orders, openOrder)
			}
		}
	} else {
		for _, openOrder := range openOrders.SaleOrder {
			if parseOrderID(openOrder.ClientOrderID) == key {
				orders = append(orders, openOrder)
			}
		}
	}

	return orders
}
func (openOrders *OpenOrders) CancelOrders(order *libBinance.Order) {
	var (
		orders []*libBinance.Order
	)
	if order.Side == libBinance.SideTypeBuy {
		for _, openOrder := range openOrders.BuyOrder {
			if openOrder.ClientOrderID != order.ClientOrderID {
				orders = append(orders, openOrder)
			}
		}
		openOrders.BuyOrder = orders
	} else {
		for _, openOrder := range openOrders.SaleOrder {
			if openOrder.ClientOrderID != order.ClientOrderID {
				orders = append(orders, openOrder)
			}
		}
		openOrders.SaleOrder = orders
	}
	buildC <- struct{}{}
}

func (openOrders *OpenOrders) UpdateOrders(order *libBinance.Order) {
	if order.Side == libBinance.SideTypeBuy {
		for _, openOrder := range openOrders.BuyOrder {
			if openOrder.ClientOrderID == order.ClientOrderID {
			}
		}
	} else {
		for _, openOrder := range openOrders.SaleOrder {
			if openOrder.ClientOrderID == order.ClientOrderID {
			}
		}
	}
}

func parseOrderID(clientOrderID string) byte {
	var (
		keyD = clientOrderID[19:]
	)

	keyInt, err := strconv.Atoi(keyD)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprint("Error: ", err))
		return ' '
	}

	return fmt.Sprintf("%c", keyInt)[0]
}

func buildOpenBuyOrderTable() []*giu.TableRowWidget {
	var rows []*giu.TableRowWidget

	rows = append(rows, giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE))

	sortOrder(OpenOrdersInstance.BuyOrder)
	for i := range OpenOrdersInstance.BuyOrder {
		order := OpenOrdersInstance.BuyOrder[i]
		price, _ := strconv.ParseFloat(order.Price, 64)
		quantity, _ := strconv.ParseFloat(order.OrigQuantity, 64)
		rows = append(rows, giu.TableRow(
			giu.Label(fmt.Sprintf("%c", parseOrderID(order.ClientOrderID))),
			giu.Label(fmt.Sprintf("%v", priceFloat648Point(order.Price))),
			giu.Style().
				SetColor(giu.StyleColorText, global.GREEN).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		))
	}
	return rows
}

func buildOpenSaleOrderTable() []*giu.TableRowWidget {
	var rows []*giu.TableRowWidget

	rows = append(rows, giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE))

	sortOrder(OpenOrdersInstance.SaleOrder)
	for i := range OpenOrdersInstance.SaleOrder {
		order := OpenOrdersInstance.SaleOrder[i]
		price, _ := strconv.ParseFloat(order.Price, 64)
		quantity, _ := strconv.ParseFloat(order.OrigQuantity, 64)
		rows = append(rows, giu.TableRow(
			giu.Label(fmt.Sprintf("%c", parseOrderID(order.ClientOrderID))),
			giu.Label(fmt.Sprintf("%v", priceFloat648Point(order.Price))),
			giu.Style().
				SetColor(giu.StyleColorText, global.RED).
				To(
					giu.Label(fmt.Sprintf("%.2fK", price*quantity/1000)),
				),
		))
	}

	return rows
}

func sortOrder(orders []*libBinance.Order) {
	sort.Slice(orders, func(i, j int) bool {
		return parseFloat(orders[i].Price) > parseFloat(orders[j].Price)
	})
}

func parseFloat(f string) float64 {
	ff, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return 0
	}
	return ff
}
