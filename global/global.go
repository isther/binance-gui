package global

import "time"

var (
	Ticker = time.NewTicker(time.Millisecond * 1000)

	Symbol = "BTCUSDT"
	Levels = 20
	Limit  = 500
)
