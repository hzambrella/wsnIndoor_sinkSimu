package money

import (
	"testing"
)

var testCmp = []struct {
	a    float64
	b    float64
	prec uint
	out  int
}{
	{1, 1, 0, 0},
	{1, 1, 1, 0},
	{1, 1, 10, 0},

	{1.0, 1.0, 0, 0},
	{1.0, 1.0, 1, 0},
	{1.0, 1.0, 10, 0},

	{1.01, 1.00, 0, 0},
	{1.01, 1.00, 1, 0},
	{1.01, 1.00, 10, 1},

	{1.0000001, 1, 0, 0},
	{1.0000001, 1, 1, 0},
	{1.0000001, 1, 10, 1},

	{1.2999999, 1.3, 0, 0},
	{1.2999999, 1.3, 1, 0},
	{1.2999999, 1.3, 10, -1},

	{1.2000009, 1.200001, 6, 0},
	{1.2000009, 1.2000005, 6, 0},
	{1.2000009, 1.2000004, 6, 1},
	{1.2000009, 1.20000008, 7, 1},
}

var testRound = []struct {
	in   float64
	prec uint
	out  float64
}{
	{2, 10, 2},
	{1.05, 10, 1.05},
	{1.00, 10, 1.00},
	{1.05, 1, 1.1},
	{1.04, 1, 1.0},
	{1.05999, 1, 1.1},
	{1.04999, 1, 1.0},
	{1.04999, 1, 1.0},
	{1.04999, 0, 1.0},
	{1.54999, 0, 2.0},
}

func TestRound(t *testing.T) {
	for index, val := range testRound {
		result := Round(val.in, val.prec)
		if result != val.out {
			t.Fatal(index, result, val.out)
		}
	}
}

var testRoundUp = []struct {
	in   float64
	prec uint
	out  float64
}{
	{2, 10, 2},
	{1.05, 10, 1.05},

	{1.00, 10, 1.00},
	{1.05, 1, 1.1},
	{1.04, 1, 1.1},

	{1.00, 1, 1.0},
	{1.05999, 1, 1.1},
	{1.04999, 1, 1.1},
	{1.04999, 1, 1.1},

	{1.00, 0, 1.0},
	{1.40, 0, 2.0},
	{1.50, 0, 2.0},
}

func TestRoundUp(t *testing.T) {
	for index, val := range testRoundUp {
		result := RoundUp(val.in, val.prec)
		if result != val.out {
			t.Fatal(index, result, val.out)
		}
	}
}

var testRoundDown = []struct {
	in   float64
	prec uint
	out  float64
}{
	{2, 10, 2},
	{1.05, 10, 1.05},
	{1.00, 10, 1.00},
	{1.05, 1, 1.0},
	{1.04, 1, 1.0},
	{1.05999, 1, 1.0},
	{1.04999, 1, 1.0},
	{1.04999, 1, 1.0},
	{1.00, 1, 1.0},

	{1.00, 0, 1.0},
	{1.40, 0, 1.0},
	{1.90, 0, 1.0},
}

func TestRoundDown(t *testing.T) {
	for index, val := range testRoundDown {
		result := RoundDown(val.in, val.prec)
		if result != val.out {
			t.Fatal(index, result, val.out)
		}
	}
}
