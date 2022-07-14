package hotkey

var HotKeySale = []byte{
	'A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L', ';',
	'Z', 'X', 'C', 'V', 'B', 'N', 'M', ',', '.', '/',
}

var HotKeyBuy = []byte{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
	'Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P',
}

func GetSaleKeyIndex(key byte) int {
	for i := range HotKeySale {
		if key == HotKeySale[i] {
			return i
		}
	}
	return -1
}

func GetBuyKeyIndex(key byte) int {
	for i := range HotKeyBuy {
		if key == HotKeyBuy[i] {
			return i
		}
	}
	return -1
}
