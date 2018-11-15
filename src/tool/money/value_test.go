package money

import (
	"testing"
)

var testFormat = []struct {
	value  Value
	format uint
	out    string
}{
	{
		New(1.44),
		0,
		"1",
	},
	{
		New(1.44),
		1,
		"1.4",
	},
	{
		New(1.44),
		3,
		"1.440",
	},

	{
		New(1.55),
		0,
		"2",
	},
	{
		New(1.55),
		1,
		"1.6",
	},
	{
		New(1.55),
		3,
		"1.550",
	},
	{
		New(1.99),
		0,
		"2",
	},
	{
		New(1.99),
		1,
		"2.0",
	},
	{
		New(1.99),
		3,
		"1.990",
	},
	{
		New(1.625),
		2,
		"1.63",
	},
}

func TestFormat(t *testing.T) {
	for _, val := range testFormat {
		out := val.value.Format(val.format)
		if out != val.out {
			t.Fatal(out, val)
			return
		}
	}
}

var testParse = []struct {
	in    string
	value Value
	err   error
}{
	{
		"1",
		New(1),
		nil,
	},
	{
		"1.4",
		New(1.4),
		nil,
	},
	{
		"1.4401",
		New(1.4401),
		nil,
	},
	{
		"0x1",
		New(0),
		ErrInput,
	},
}

func TestParse(t *testing.T) {
	for _, val := range testParse {
		m, err := Parse(val.in)
		if err != nil {
			if err!=val.err{
				t.Fatal(val, err)
				return
			}
			
		}
		if m.Cmp(val.value) != 0 {
			t.Fatal(val, m)
			return
		}
	}
}

var testString = []struct {
	in  Value
	out int
}{
	{New(30.0), 2},
	{New(1.00), 1},
	{New(0.1), 3},
	{New(0.1000), 3},
	{New(0.11), 4},
	{New(0.111111), 6},
}

func TestString(t *testing.T) {
	for index, s := range testString {
		result := len(s.in.String())
		if result != s.out {
			t.Fatal(index, result, s.out)
		}
	}
}
