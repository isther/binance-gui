package binance

import "sync"

var CostInstance *Cost

type Cost struct {
	Quantity   float64
	GlobalCost float64

	mu sync.Mutex
}

func init() {
	CostInstance = NewCost()
}

func ResetCostInstance() {
	CostInstance = NewCost()
}

func NewCost() *Cost {
	return &Cost{
		Quantity:   0,
		GlobalCost: 0,
	}
}

func (c *Cost) Buy(quantity float64, price float64) *Cost {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Quantity += quantity
	c.GlobalCost += price * quantity
	return c
}

func (c *Cost) Sale(quantity float64, price float64) *Cost {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Quantity == 0 {
		return NewCost()
	}

	c.Quantity -= quantity
	c.GlobalCost -= price * quantity
	return c
}

func (c *Cost) UpdateAverageCode() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Quantity == 0 {
		return "0.0"
	}

	return correction(c.GlobalCost/c.Quantity, AccountInstance.PriceFilter.tickSize)
}
