package binance

import (
	libBinance "github.com/adshao/go-binance/v2"
)

func init() {

	// websocket
	libBinance.WebsocketKeepalive = true
}
