package jsonquery

import (
	"encoding/json"
	"strings"
	"testing"
)

const TestData = `{
	"foo": 1,
	"bar": 2,
	"test": "Hello, world!",
	"baz": 123.1,
	"numstring": "42",
	"floatstring": "42.1",
	"array": [
		{"foo": 1},
		{"bar": 2},
		{"baz": 3}
	],
	"subobj": {
		"foo": 1,
		"subarray": [1,2,3],
		"subsubobj": {
			"bar": 2,
			"baz": 3,
			"array": ["hello", "world"]
		}
	},
	"collections": {
		"bools": [false, true, false],
		"strings": ["hello", "strings"],
		"numbers": [1,2,3,4],
		"arrays": [[1.0,2.0],[2.0,3.0],[4.0,3.0]],
		"objects": [
			{"obj1": 1},
			{"obj2": 2}
		]
	},
	"bool": true
}`

func tErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
}

func TestQuery(t *testing.T) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(strings.NewReader(TestData))
	err := dec.Decode(&data)
	tErr(t, err)
	q := NewQuery(data)

	ival, err := q.Int("foo")
	if ival != 1 {
		t.Errorf("Expecting 1, got %v\n", ival)
	}

	ival, err = q.Int("bar")
	if ival != 2 {
		t.Errorf("Expecting 2, got %v\n", ival)
	}
	tErr(t, err)

	ival, err = q.Int("subobj", "foo")
	if ival != 1 {
		t.Errorf("Expecting 1, got %v\n", ival)
	}
	tErr(t, err)

	sval, err := q.String("test")
	if sval != "Hello, world!" {
		t.Errorf("Expecting \"Hello, World!\", got \"%v\"\n", sval)
	}

	astrings, err := q.ArrayOfStrings("collections", "strings")
	tErr(t, err)
	if astrings[0] != "hello" {
		t.Errorf("Expecting hello, got %v\n", astrings[0])
	}

	aa, err := q.ArrayOfArrays("collections", "arrays")
	tErr(t, err)
	if aa[0][0].(float64) != 1 {
		t.Errorf("Expecting 1, got %v\n", aa[0][0])
	}

	aobjs, err := q.ArrayOfObjects("collections", "objects")
	tErr(t, err)
	if aobjs[0]["obj1"].(float64) != 1 {
		t.Errorf("Expecting 1, got %v\n", aobjs[0]["obj1"])
	}
}
