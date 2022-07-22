package binance

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/AllenDang/giu"
	"github.com/isther/binanceGui/conf"
	"github.com/isther/binanceGui/console"
	"github.com/isther/binanceGui/global"
	"github.com/isther/binanceGui/utils"
)

var (
	earlyWarningTable1m []*giu.TableRowWidget
	earlyWarningTable3m []*giu.TableRowWidget
)

func GetEarlyWaringTable1m() []*giu.TableRowWidget {
	return earlyWarningTable1m
}

func GetEarlyWaringTable3m() []*giu.TableRowWidget {
	return earlyWarningTable3m
}

func ListenEarlyWarning() {
	go func() {
		for {
			if global.EarlyWarning {
				filterAndBuild()
				<-time.Tick(1 * time.Minute)
			}
		}
	}()
}

func filterAndBuild() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		earlyWarningTable1m = NewEarlyWarner("1m").setThreshold(float64(global.EarlyWarning1mAmplitude), float64(global.EarlyWarning1mTurnOver)).updateKlines().buildTable()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		earlyWarningTable3m = NewEarlyWarner("3m").setThreshold(float64(global.EarlyWarning3mAmplitude), float64(global.EarlyWarning3mTurnOver)).updateKlines().buildTable()
	}()
	wg.Wait()
}

type Indicator struct {
	Symbol    string
	Amplitude float64
	Turnover  float64
}

type EarlyWarner struct {
	threshold Indicator
	set       (map[string]*Indicator)
	interval  string
	mu        sync.Mutex
}

func NewEarlyWarner(interval string) *EarlyWarner {
	return &EarlyWarner{
		set:      make(map[string]*Indicator),
		interval: interval,
	}
}

func (this *EarlyWarner) buildTable() []*giu.TableRowWidget {
	var (
		keySlice []string
		rows     []*giu.TableRowWidget
	)

	for k := range this.set {
		keySlice = append(keySlice, k)
	}

	sort.SliceStable(keySlice, func(i, j int) bool {
		return this.set[keySlice[i]].Amplitude > this.set[keySlice[j]].Amplitude
	})

	rows = append(rows, giu.TableRow(
		giu.Label("BaseAsset"),
		giu.Label("Amplitude"),
		giu.Label("TurnOver"),
	).BgColor(global.PURPLE))

	for i := range keySlice {
		var key = keySlice[i]

		rows = append(rows, giu.TableRow(
			giu.Label(this.set[key].Symbol),
			giu.Label(fmt.Sprintf("%.2f%%", this.set[key].Amplitude)),
			giu.Label(fmt.Sprintf("%.2fK", this.set[key].Turnover/1000)),
		))
	}
	if len(rows) > 1 {
		utils.WinSound()
	}
	return rows
}

func (this *EarlyWarner) setThreshold(amplitude, turnover float64) *EarlyWarner {
	this.threshold.Amplitude = amplitude
	this.threshold.Turnover = turnover

	return this
}

func (this *EarlyWarner) updateKlines() *EarlyWarner {
	var wg sync.WaitGroup
	for i := range conf.Conf.Symbols {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			this.updateAKLine(conf.Conf.Symbols[i])
		}(i)
	}
	wg.Wait()

	return this
}

func (this *EarlyWarner) updateAKLine(symbol string) {
	var err error

	res, err := GetClient().NewKlinesService().Symbol(symbol).Interval(this.interval).Limit(2).Do(context.Background())
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}

	this.mu.Lock()
	defer this.mu.Unlock()
	var kline = res[0]
	volume, err := strconv.ParseFloat(kline.QuoteAssetVolume, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return
	}
	amplitude := getAmplitude(kline.High, kline.Low, kline.Close) * 100

	if amplitude >= this.threshold.Amplitude && volume >= this.threshold.Turnover {
		this.set[symbol] = &Indicator{
			Symbol:    symbol,
			Amplitude: amplitude,
			Turnover:  volume,
		}
	}
}

func getAmplitude(highStr, lowerStr, closeStr string) float64 {
	var err error
	high, err := strconv.ParseFloat(highStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return 0.0
	}
	lower, err := strconv.ParseFloat(lowerStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return 0.0
	}
	close, err := strconv.ParseFloat(closeStr, 64)
	if err != nil {
		console.ConsoleInstance.Write(fmt.Sprintf("Error: %v", err))
		return 0.0
	}
	return (high - lower) / close
}
