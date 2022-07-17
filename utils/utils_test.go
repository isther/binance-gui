package utils

import "testing"

func TestCorrection(t *testing.T) {
	cases := []struct {
		name     string
		val      float64
		stepSize string
	}{
		{"1", 2.22222222, "0.00000001"},
		{"2", 2.22222220, "0.00000001"},
		{"3", 22.2222222, "1.0"},
		{"4", 222.222, "10.0"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := correction(c.val, c.stepSize)
			t.Log(res)
		})
	}
}

func TestRound(t *testing.T) {
	cases := []struct {
		name      string
		val       float64
		precision int
	}{
		{"1", 555.555, 1},
		{"2", 555.555, -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Log(c.name, ": ", RoundLower(c.val, c.precision))
			t.Log(c.name, ": ", RoundUpper(c.val, c.precision))
		})
	}
}
