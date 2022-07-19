package global

import "image/color"

var (
	Ping  string
	Debug = false

	// 分仓参数
	Average              int32 = 7
	AverageSymbol1Amount float64
	AverageSymbol2Amount float64

	// 切换交易对channel
	FreshC = make(chan string)

	// 交易模式参数
	TradeMode            TransactionMode = FiveAfterMulPoint
	VolatilityRatiosBuy  [20]float32
	VolatilityRatiosSale [20]float32
	VolatilityRatiosF1   float32
	VolatilityRatiosF2   float32
	VolatilityRatiosF5   float32
	VolatilityRatiosF6   float32
	VolatilityRatiosF12  float32

	// 大单提醒档位
	// 成交明细
	AggTradeBigOrderReminder [5]int32
	IsShieldSmallOrder       = true
	// 订单簿2
	Order2BigOrderReminder [5]int32
	// 订单簿1
	Order1BigOrderReminder [5]int32

	// ws参数
	Levels = 20
	Limit  = 500

	// 热键设置
	HotKeyRun   = false
	HotKeyColor = BLUE

	Order2FontSize float32 = 16
)

var (
	RED      = color.RGBA{0xFF, 0x33, 0x33, 0xFF}
	YELLOW   = color.RGBA{0xFF, 0x99, 0x00, 0xFF}
	Order2Bg = color.RGBA{0xFF, 0x66, 0x00, 0xFF}
	BLUE     = color.RGBA{0x00, 0x66, 0xCC, 0xFF}
	BLUE2    = color.RGBA{0x33, 0x66, 0xFF, 0xFF}
	// BLACK  = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	BLACK = color.RGBA{0x7F, 0xFF, 0x00, 0xFF}
	WHITE = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

	GREEN  = color.RGBA{0x66, 0xCC, 0x00, 0xFF}
	PURPLE = color.RGBA{0x33, 0x33, 0xFF, 0xFF}
)

type TransactionMode int

const (
	_ = iota
	AllPlusOneSize
	FiveAfterMulPoint
)

func init() {
	VolatilityRatiosSale[5] = 1.003
	VolatilityRatiosSale[6] = 1.006
	VolatilityRatiosSale[7] = 1.009
	VolatilityRatiosSale[8] = 1.012
	VolatilityRatiosSale[9] = 1.015

	VolatilityRatiosSale[10] = 1.02
	VolatilityRatiosSale[11] = 1.025
	VolatilityRatiosSale[12] = 1.03
	VolatilityRatiosSale[13] = 1.035
	VolatilityRatiosSale[14] = 1.04
	VolatilityRatiosSale[15] = 1.05
	VolatilityRatiosSale[16] = 1.06
	VolatilityRatiosSale[17] = 1.07
	VolatilityRatiosSale[18] = 1.08
	VolatilityRatiosSale[19] = 1.09

	VolatilityRatiosBuy[5] = 0.997
	VolatilityRatiosBuy[6] = 0.994
	VolatilityRatiosBuy[7] = 0.991
	VolatilityRatiosBuy[8] = 0.988
	VolatilityRatiosBuy[9] = 0.985

	VolatilityRatiosBuy[10] = 0.980
	VolatilityRatiosBuy[11] = 0.975
	VolatilityRatiosBuy[12] = 0.970
	VolatilityRatiosBuy[13] = 0.965
	VolatilityRatiosBuy[14] = 0.96
	VolatilityRatiosBuy[15] = 0.95
	VolatilityRatiosBuy[16] = 0.94
	VolatilityRatiosBuy[17] = 0.93
	VolatilityRatiosBuy[18] = 0.92
	VolatilityRatiosBuy[19] = 0.91

	VolatilityRatiosF1 = 1.005
	VolatilityRatiosF2 = 1.005
	VolatilityRatiosF5 = 0.995
	VolatilityRatiosF6 = 0.995
	VolatilityRatiosF12 = 0.98

	AggTradeBigOrderReminder[1] = 1000
	AggTradeBigOrderReminder[2] = 3000
	AggTradeBigOrderReminder[3] = 6000
	AggTradeBigOrderReminder[4] = 10000
	Order2BigOrderReminder[1] = 2000
	Order2BigOrderReminder[2] = 10000
	Order2BigOrderReminder[3] = 35000
	Order2BigOrderReminder[4] = 100000
	Order1BigOrderReminder[1] = 10000
	Order1BigOrderReminder[2] = 30000
	Order1BigOrderReminder[3] = 50000
	Order1BigOrderReminder[4] = 100000
}

func GetHotKeyStatus() string {
	if HotKeyRun {
		return "开启"
	}
	return "关闭"
}

func GetTradeMode() string {
	if TradeMode == AllPlusOneSize {
		return "模式一"
	}
	return "模式二"
}

func ReverseHotKeyStatus() {
	HotKeyRun = !HotKeyRun
	if HotKeyColor == BLUE {
		HotKeyColor = RED
		return
	}
	HotKeyColor = BLUE
}
