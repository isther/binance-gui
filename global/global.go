package global

var (
	Ping string

	Average              int32 = 10
	AverageSymbol1Amount float64
	AverageSymbol2Amount float64

	FreshC = make(chan string)

	TradeMode TransactionMode = AllPlusOneSize

	Levels = 20
	Limit  = 500

	HotKeyRun = false

	Order2FontSize float32 = 16
)

type TransactionMode int

const (
	_ = iota
	AllPlusOneSize
	FiveAfterMulPoint
)

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

func ReverseTradeMode() {
	if TradeMode == AllPlusOneSize {
		TradeMode = FiveAfterMulPoint
		return
	}
	TradeMode = AllPlusOneSize
}
