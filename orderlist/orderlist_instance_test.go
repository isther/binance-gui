package orderlist

import (
	"fmt"
	"sync"
	"testing"
	"time"

	libBinance "github.com/adshao/go-binance/v2"
)

func TestOrderListInstance(t *testing.T) {
	newABuyOrderWithKey := func(key string) *libBinance.Order {
		return &libBinance.Order{
			ClientOrderID: fmt.Sprintf("%d%s", time.Now().UnixNano(), key),
			Side:          libBinance.SideTypeBuy,
			Price:         "1",
			OrigQuantity:  "1",
		}
	}
	for i := 0; i < 10; i++ {
		OrderListInstance.AddOrders(newABuyOrderWithKey("1"))
	}
	OrderListInstance.print()

	orders := OrderListInstance.GetOrders("1")
	if orders == nil {
		t.Log("No orders")
		return
	}

	var wg sync.WaitGroup
	for i, order := range orders {
		if i < 5 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				OrderListInstance.cancelOrderByID(order.ClientOrderID)
			}()
		}
	}
	wg.Wait()
	OrderListInstance.print()

	OrderListInstance.Clear()
	OrderListInstance.print()
}
