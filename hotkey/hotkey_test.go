package hotkey

import (
	"fmt"
	"testing"
)

func TestHotKeyASCII(t *testing.T) {
	for i := range HotKeyBuy {
		fmt.Printf("%c: %d\n", HotKeyBuy[i], HotKeyBuy[i])
	}

	for i := range HotKeySale {
		fmt.Printf("%c: %d\n", HotKeySale[i], HotKeySale[i])
	}
}
