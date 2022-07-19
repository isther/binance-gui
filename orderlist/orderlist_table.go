package orderlist

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/global"

	libBinance "github.com/adshao/go-binance/v2"
)

var (
	orderListBuyTable  []*giu.TableRowWidget
	orderListSaleTable []*giu.TableRowWidget

	// buildOrderListTableCh = make(chan struct{}, 100)
)

func StartBuildingOrderListTable() {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			orderListBuyTable = buildOpenBuyOrderTable()
			orderListSaleTable = buildOpenSaleOrderTable()
			<-ticker.C
		}
	}()
}

func GetOpenBuyOrdersTable() []*giu.TableRowWidget {
	return orderListBuyTable
}

func GetOpenSaleOrdersTable() []*giu.TableRowWidget {
	return orderListSaleTable
}

func buildOpenBuyOrderTable() []*giu.TableRowWidget {
	var (
		rows      []*giu.TableRowWidget
		buyOrders = OrderListInstance.OutPutBuyOrders()
	)

	rows = append(rows, giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE))

	sortOrder(buyOrders)
	for i := range buyOrders {
		order := buyOrders[i]
		price, _ := strconv.ParseFloat(order.Price, 64)
		quantity, _ := strconv.ParseFloat(order.OrigQuantity, 64)
		rows = append(rows, giu.TableRow(
			giu.Label(fmt.Sprintf("%s", parseOrderIDToKey(order.ClientOrderID))),
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
	var (
		rows       []*giu.TableRowWidget
		saleOrders = OrderListInstance.OutPutSaleOrders()
	)

	rows = append(rows, giu.TableRow(
		giu.Label("快捷键"),
		giu.Label("价格"),
		giu.Label("成交额"),
	).BgColor(global.PURPLE))

	sortOrder(saleOrders)
	for i := range saleOrders {
		order := saleOrders[i]
		price, _ := strconv.ParseFloat(order.Price, 64)
		quantity, _ := strconv.ParseFloat(order.OrigQuantity, 64)
		rows = append(rows, giu.TableRow(
			giu.Label(fmt.Sprintf("%s", parseOrderIDToKey(order.ClientOrderID))),
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

func priceFloat648Point(s string) string {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] == '.' {
			return string(append([]byte(s[:]), '0'))
		}
		if s[i] != '0' {
			return s
		}
		s = s[:i]
	}
	return s
}
