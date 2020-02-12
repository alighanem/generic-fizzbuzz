package main

import "testing"

func TestCounter_Inc(t *testing.T) {
	type testCase struct {
		Converter    *ValueConverter
		ExpectedHits int
	}

	cases := map[string]testCase{
		"new converter": {
			Converter: &ValueConverter{
				Value1: 2,
				Word1:  "foo",
				Value2: 5,
				Word2:  "bar",
				Limit:  100,
			},
			ExpectedHits: 1,
		},
		"converter already added": {
			Converter: &ValueConverter{
				Value1: 3,
				Word1:  "fizz",
				Value2: 5,
				Word2:  "buzz",
				Limit:  100,
			},
			ExpectedHits: 2,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			counter := setupCounter()
			counter.Inc(c.Converter)

			metric, ok := counter.Metrics[c.Converter.Key()]
			if !ok {
				t.Fatal("converter not found", "key", c.Converter.Key())
			}

			if metric.Hits != c.ExpectedHits {
				t.Fatal("unexpected number of hits", "expected", c.ExpectedHits, "actual", metric.Hits)
			}
		})
	}
}

func TestCounter_GetMaxHits(t *testing.T) {
	type testCase struct {
		Counter      *Counter
		ExpectedKey  string
		ExpectedHits int
	}

	cases := map[string]testCase{
		"get max hits": {
			Counter:      setupCounter(),
			ExpectedKey:  "2_7_100_foo_buzz",
			ExpectedHits: 3,
		},
		"no hits saved": {
			Counter:      &Counter{Metrics: map[string]*Metric{}},
			ExpectedKey:  "",
			ExpectedHits: 0,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			max := c.Counter.GetMaxHits()

			var key string
			var hits int
			if max != nil {
				key = max.Converter.Key()
				hits = max.Hits
			}

			if c.ExpectedKey != key {
				t.Fatal("unexpected converter found", "expected", c.ExpectedKey, "actual", key)
			}

			if c.ExpectedHits != hits {
				t.Fatal("unexpected number of hits", "expected", c.ExpectedHits, "actual", hits)
			}
		})
	}
}

func setupCounter() *Counter {
	converter1 := ValueConverter{
		Value1: 3,
		Word1:  "fizz",
		Value2: 5,
		Word2:  "buzz",
		Limit:  100,
	}

	converter2 := ValueConverter{
		Value1: 2,
		Word1:  "foo",
		Value2: 7,
		Word2:  "buzz",
		Limit:  100,
	}

	c := Counter{Metrics: map[string]*Metric{
		converter1.Key(): {
			Converter: &converter1,
			Hits:      1,
		},
		converter2.Key(): {
			Converter: &converter2,
			Hits:      3,
		},
	}}

	return &c
}
