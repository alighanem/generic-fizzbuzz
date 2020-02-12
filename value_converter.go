package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ValueConverter struct {
	Value1 int    `json:"value1"`
	Word1  string `json:"word1"`
	Value2 int    `json:"value2"`
	Word2  string `json:"word2"`
	Limit  int    `json:"limit"`
}

func (v *ValueConverter) convert(value int) (string, error) {
	var result strings.Builder
	var printed bool

	if value%v.Value1 == 0 {
		_, err := result.WriteString(v.Word1)
		if err != nil {
			return "", err
		}
		printed = true
	}

	if value%v.Value2 == 0 {
		_, err := result.WriteString(v.Word2)
		if err != nil {
			return "", err
		}
		printed = true
	}

	if !printed {
		return strconv.Itoa(value), nil
	}

	return result.String(), nil
}

// Key returns the unique key of the converter
func (v *ValueConverter) Key() string {
	return fmt.Sprintf("%v_%v_%v_%v_%v", v.Value1, v.Value2, v.Limit, v.Word1, v.Word2)
}

// Generate generates integers and returns their corresponding string representation
func (v *ValueConverter) Generate() ([]string, error) {
	var results []string
	for i := 1; i <= v.Limit; i++ {
		v, err := v.convert(i)
		if err != nil {
			return nil, fmt.Errorf("cannot convert value: %v - err: %w", i, err)
		}
		results = append(results, v)
	}

	return results, nil
}
