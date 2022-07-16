package test

import "testing"

func TestSize(t *testing.T) {
	cases := []struct {
		name     string
		val      string
		stepSize string
	}{
		{"1", "0.00100000", "1.10000000"},
		{"2", "0.00100000", "1.01000000"},
		{"3", "0.00100000", "1.00100000"},
		{"4", "0.00100000", "1.00010000"},
		{"5", "0.10000000", "1.10000000"},
		{"6", "1.00000000", "1.10000000"},
		{"7", "0.00000001", "0.10000000"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var (
				start  bool
				length = 0
			)

			for i := len(c.stepSize) - 1; i > 0; i-- {
				if c.stepSize[i] == '1' {
					start = true
				}
				if start {
					length++
				}
			}

			for i := range c.val {
				if c.val[i] == '.' {
					c.val = c.val[:i+length]
					break
				}
			}

			t.Logf("name: %s length: %v quantity: %v", c.name, length, c.val)
		})
	}
}
