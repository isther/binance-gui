package global

import "time"

var (
	Ticker = time.NewTicker(time.Millisecond * 100)

	Symbol        = "BUSDUSDT"
	Symbol1       = "BUSD"
	Symbol1Free   string
	Symbol1Locked string
	Symbol2       = "USDT"
	Symbol2Free   string
	Symbol2Locked string

	Levels = 20
	Limit  = 500

	FreshC = make(chan string)

	Order2FontSize float32 = 16
)
