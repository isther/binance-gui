package binance

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	libBinance "github.com/adshao/go-binance/v2"
	"github.com/isther/binance/conf"
	"github.com/tidwall/gjson"
)

var (
	ClientTimeOut = 5 * time.Second
)

func GetClient() *libBinance.Client {
	client := libBinance.NewClient(conf.Conf.ApiKey, conf.Conf.SecretKey)
	return client
}

func SystemTime() int64 {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v3/time", nil)
	if err != nil {
		panic(err)
	}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	timestamp := gjson.Get(string(body), "serverTime").Int()
	fmt.Println(timestamp)
	return timestamp - currentTimestamp()
}

func currentTimestamp() int64 {
	return formatTimestamp(time.Now())
}

// FormatTimestamp formats a time into Unix timestamp in milliseconds, as requested by Binance.
func formatTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
