package binance

import (
	"context"
	"fmt"
	"time"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binanceGui/conf"
)

var (
	ClientTimeOut = 5 * time.Second
)

func GetClient() *libBinance.Client {
	return libBinance.NewClient(conf.Conf.ApiKey, conf.Conf.SecretKey)
}

func GetAll() {
	info, err := GetClient().NewGetAllCoinsInfoService().Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(info[0])
}

func GetApiPermission() {
	res, err := GetClient().NewGetAPIKeyPermission().Do(context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Println("ipRestrict: ", res.IPRestrict)
	fmt.Println("createTime: ", res.CreateTime)
	fmt.Println("enableWithDrawals: ", res.EnableWithdrawals)
	fmt.Println("enableInternalTransfer: ", res.EnableInternalTransfer)
	fmt.Println("permitsUniversalTransfer: ", res.PermitsUniversalTransfer)
	fmt.Println("enableVanillaOptions: ", res.EnableVanillaOptions)
	fmt.Println("enableReading: ", res.EnableReading)
	fmt.Println("enableFutures: ", res.EnableFutures)
	fmt.Println("enableMargin: ", res.EnableMargin)
	fmt.Println("enableSpotAndMarginTrading: ", res.EnableSpotAndMarginTrading)
	fmt.Println("tradingAuthorityExpirationTime: ", res.TradingAuthorityExpirationTime)
}
