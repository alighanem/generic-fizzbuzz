package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/ali.ghanem/leboncoin-exo/golden"
)

func TestAPI_Generate(t *testing.T) {
	type testCase struct {
		Query          string
		ExpectedStatus int
	}

	cases := map[string]testCase{
		"nominal case generate succeed": {
			Query:          "int1=3&int2=5&limit=20&str1=test1&str2=test2",
			ExpectedStatus: http.StatusOK,
		},
		"no param set return errors": {
			Query:          "",
			ExpectedStatus: http.StatusBadRequest,
		},
		"some params not set return errors": {
			Query:          "int1=3&int2=5&limit=&str1=fizz&tsr2=",
			ExpectedStatus: http.StatusBadRequest,
		},
		"int params are invalid return errors": {
			Query:          "int1=fail&int2=not&limit=10&str1=test1&str2=test2",
			ExpectedStatus: http.StatusBadRequest,
		},
		"limit too low return errors": {
			Query:          "int1=3&int2=5&limit=1&str1=test1&str2=test2",
			ExpectedStatus: http.StatusBadRequest,
		},
		"int params are negative return errors": {
			Query:          "int1=-3&int2=-5&limit=10&str1=test1&str2=test2",
			ExpectedStatus: http.StatusBadRequest,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// setup
			r, err := http.NewRequest("GET", fmt.Sprintf("/generate?%s", c.Query), nil)
			if err != nil {
				t.Fatal(err)
			}

			var a = API{
				Counter: &Counter{
					RWMutex: sync.RWMutex{},
					Metrics: make(map[string]*Metric),
				},
				converters: make(chan *ValueConverter),
			}
			go a.consumeConverters()
			handler := http.HandlerFunc(a.generate)
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			// results
			resp := w.Result()

			// assert
			if resp.StatusCode != c.ExpectedStatus {
				t.Errorf("Unexpected status code %d", resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			golden.JSONEq(t, t.Name()+".json", body)
			close(a.converters)
		})
	}
}

func TestAPI_Statistics(t *testing.T) {
	t.Run("get statistics", func(t *testing.T) {
		// setup
		r, err := http.NewRequest("GET", "/statistics", nil)
		if err != nil {
			t.Fatal(err)
		}

		var a = API{
			Counter:    setupCounter(),
			converters: make(chan *ValueConverter),
		}
		go a.consumeConverters()
		handler := http.HandlerFunc(a.statistics)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		// results
		resp := w.Result()

		// assert
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Unexpected status code %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		golden.JSONEq(t, t.Name()+".json", body)
		close(a.converters)
	})
}
