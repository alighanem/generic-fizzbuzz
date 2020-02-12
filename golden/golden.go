package golden

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"
)

// JSONEq asserts that the content of a golden file
// is same as json received by the test
func JSONEq(t *testing.T, goldenFile string, actual json.RawMessage) {
	var expected json.RawMessage
	readJSON(t, goldenFile, &expected)
	equal, err := jsonBytesEqual(expected, actual)
	if err != nil {
		t.Fatal("unexpected error", "err", err)
	}
	if !equal {
		t.Fatal("Not equal", "expected", string(expected), "actual", string(actual))
	}
}

func jsonBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

func readJSON(t *testing.T, name string, v interface{}) {
	var (
		raw []byte
		err error
	)

	raw, err = ioutil.ReadFile(fmt.Sprintf("testdata/%s", name))
	if err != nil {
		t.Fatalf("reading test file: %v", err)
	}

	err = json.Unmarshal(raw, &v)
	if err != nil {
		t.Fatalf("unmarshalling file value: %s", err)
	}
}
