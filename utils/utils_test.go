package utils

import "testing"

func TestFloat64ToStringLen3(t *testing.T) {
	cases := []struct {
		name string
		f    float64
	}{
		{"1", 314.152},
		{"2", 0.00314152},
		{"3", 0.14152},
		{"4", 0.00000152},
		{"5", 0.000012},
		{"6", 0.001},
		{"7", 1.000234},
		{"8", 1.01},
		{"9", 1.0123},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			res := Float64ToStringLen3(c.f)
			t.Log(res)
		})
	}
}

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
