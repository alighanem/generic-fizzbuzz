package main

import "testing"

func TestValueConverter_Convert(t *testing.T) {
	type TestCase struct {
		Values         []int
		ExpectedResult string
	}

	cases := map[string]TestCase{
		"only multiples of 3 return fizz": {
			Values:         []int{3, 9, 18},
			ExpectedResult: "fizz",
		},
		"only multiples of 5 return buzz": {
			Values:         []int{5, 10, 20, 40},
			ExpectedResult: "buzz",
		},
		"multiples of 15 return fizzbuzz": {
			Values:         []int{15, 30, 75},
			ExpectedResult: "fizzbuzz",
		},
		"nominal values return string representation": {
			Values:         []int{11},
			ExpectedResult: "11",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			converter := ValueConverter{
				Value1: 3,
				Word1:  "fizz",
				Value2: 5,
				Word2:  "buzz",
			}
			for _, value := range c.Values {
				actual, err := converter.convert(value)
				if err != nil {
					t.Fatal("unexpected error", "err", err)
				}

				if actual != c.ExpectedResult {
					t.Fatal("unexpected result", "expected", c.ExpectedResult, "actual", actual)
				}
			}
		})
	}
}

func TestValueConverter_Key(t *testing.T) {
	type testCase struct {
		Converter   *ValueConverter
		ExpectedKey string
	}

	cases := map[string]testCase{
		"generate key": {
			Converter: &ValueConverter{
				Value1: 3,
				Word1:  "fizz",
				Value2: 5,
				Word2:  "buzz",
				Limit:  100,
			},
			ExpectedKey: "3_5_100_fizz_buzz",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			key := c.Converter.Key()
			if key != c.ExpectedKey {
				t.Fatal("wrong key generated", "expected", c.ExpectedKey, "actual", key)
			}
		})
	}
}
