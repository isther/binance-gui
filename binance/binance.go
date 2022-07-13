package binance

import (
	"time"

	libBinance "github.com/adshao/go-binance/v2"
)

func init() {
	//order
	// libBinance.UseTestnet = true

	// websocket
	libBinance.WebsocketKeepalive = true
	libBinance.WebsocketTimeout = 4 * time.Minute
}
