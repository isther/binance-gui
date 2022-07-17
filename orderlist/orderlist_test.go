package orderlist

import (
	"fmt"
	"sync"
	"testing"
	"time"

	libBinance "github.com/adshao/go-binance/v2"
)

func TestOrderListToAddOrders(t *testing.T) {
	var (
		orderList = NewOrderList()
		wg        sync.WaitGroup
	)
	orderList.print()
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			orderList.AddOrders(newASaleOrder(i))
		}(i)
	}
	wg.Wait()
	orderList.print()
}

func TestOrderListToGetOrders(t *testing.T) {
	var (
		orderList = NewOrderList()
		wg        sync.WaitGroup
	)
	orderList.print()
	for i := 1; i <= 1; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))

			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
			orderList.AddOrders(newABuyOrder(i))
		}(i)
	}
	wg.Wait()
	orderList.print()

	orders := orderList.GetOrders("A")
	if orders == nil {
		t.Log("No orders")
		return
	}

	for i, order := range orders {
		t.Log("order: ", i, " ", order)
	}
}

func TestOrderListAddOrdersAndGetOrders(t *testing.T) {
	var (
		orderList = NewOrderList()
		wg        sync.WaitGroup
	)
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			orders := orderList.GetOrders("A")
			for i, order := range orders {
				t.Log("order: ", i, " ", order)
			}
			<-ticker.C
		}
	}()

	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			orderList.AddOrders(newABuyOrder(i))
		}(i)
	}
	wg.Wait()
	orderList.print()
}

func TestOrderListCancelOrdersByKey(t *testing.T) {
	var (
		orderList = NewOrderList()
		wg        sync.WaitGroup
	)
	orderList.print()
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			orderList.AddOrders(newASaleOrder(i))
			orderList.AddOrders(newASaleOrder(i))
		}(i)
	}
	wg.Wait()
	orderList.print()

	orderList.CancelOrdersByKey("A")

	orderList.print()
}

func TestOrderListCancelOrdersByOrderID(t *testing.T) {
	var (
		orderList = NewOrderList()
		wg        sync.WaitGroup
	)
	orderList.print()

	newASaleOrderWithKey := func(key string, i int) *libBinance.Order {
		return &libBinance.Order{
			ClientOrderID: fmt.Sprintf("%d%s", time.Now().UnixNano(), key),
			Side:          libBinance.SideTypeSell,
			Price:         "1",
			OrigQuantity:  "1",
		}
	}
	for i := 1; i <= 10; i++ {
		orderList.AddOrders(newASaleOrderWithKey("A", i))
		time.Sleep(100 * time.Millisecond)
	}
	wg.Wait()
	orderList.print()

	orders := orderList.GetOrders("A")
	if orders == nil {
		t.Log("No orders")
		return
	}

	for i, order := range orders {
		if i < 5 {
			orderList.CancelOrdersByID(order.ClientOrderID)
		}
	}

	orderList.print()
}

func newABuyOrder(i int) *libBinance.Order {
	return &libBinance.Order{
		ClientOrderID: fmt.Sprintf("%d-%d", i, time.Now().UnixNano()),
		Side:          libBinance.SideTypeBuy,
		Price:         "1",
		OrigQuantity:  "1",
	}
}

func newASaleOrder(i int) *libBinance.Order {
	return &libBinance.Order{
		ClientOrderID: fmt.Sprintf("%d-%d", i, time.Now().UnixNano()),
		Side:          libBinance.SideTypeSell,
		Price:         "1",
		OrigQuantity:  "1",
	}
}
