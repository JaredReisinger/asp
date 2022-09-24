package asp

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestBetterStringToTime(t *testing.T) {
	fn := betterStringToTime()

	input := "2022-01-02T03:04:05.666-07:00"
	actual, err := fn(reflect.TypeOf(""), reflect.TypeOf(time.Time{}), input)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}
	t.Logf("converted %q to %#v", input, actual)

	expected, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if actual != expected {
		t.Logf("expected %q, got %q", expected, actual)
		t.Fail()
	}

	// other types cause input to just pass through...
	actual, err = fn(reflect.TypeOf(""), reflect.TypeOf(""), input)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if actual != input {
		t.Logf("expected %q, got %q", input, actual)
		t.Fail()
	}

	actual, err = fn(reflect.TypeOf(1), reflect.TypeOf(time.Time{}), input)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if actual != input {
		t.Logf("expected %q, got %q", input, actual)
		t.Fail()
	}
}

func TestStringToByteSlice(t *testing.T) {
	fn := stringToByteSlice()

	expected := []byte{0xde, 0xad, 0xbe, 0xef}
	actual, err := fn(reflect.ValueOf("deadbeef"), reflect.ValueOf([]byte{}))
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	// expectBytes(t, expected, actual.([]byte))
	if !bytes.Equal(expected, actual.([]byte)) {
		t.Logf("expected %q, got %q", expected, actual)
		t.Fail()
	}

	// different values pass through the 'from' unchanged...
	input := reflect.ValueOf("deadbeef")
	actual, err = fn(input, reflect.ValueOf(""))
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if actual != input.Interface() {
		t.Log("expected inupt to pass through unchanged")
		t.Fail()
	}

	input = reflect.ValueOf(1)
	actual, err = fn(input, reflect.ValueOf([]byte{}))
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if actual != input.Interface() {
		t.Log("expected input to pass through unchanged")
		t.Fail()
	}
}

func TestStringToMapStringString(t *testing.T) {
	fn := stringToMapStringString()

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	actualIntf, err := fn(reflect.ValueOf("key1=value1,key2=value2"), reflect.ValueOf(map[string]string{}))
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	actual := actualIntf.(map[string]string)

	if len(actual) != len(expected) {
		t.Logf("expected len %d, got len %d", len(expected), len(actual))
		t.Fail()
	}

	for k, v := range expected {
		va, ok := actual[k]
		if !ok {
			t.Logf("expected key %q, missing", k)
			t.Fail()
		}
		if va != v {
			t.Logf("expected %q value %q, got %q", k, v, va)
			t.Fail()
		}
	}
}

func TestBetterStringToSlice(t *testing.T) {
	fn := betterStringToSlice(",")

	cases := []struct {
		input    string
		expected []string
	}{
		{input: "one,two", expected: []string{"one", "two"}},
		{input: "[one,two]", expected: []string{"one", "two"}},
		// {input: "", expected: []string{"", ""}},
	}

	for i, c := range cases {
		actualIntf, err := fn(reflect.String, reflect.Slice, c.input)
		if err != nil {
			t.Logf("unexpected error (%d): %v", i, err)
			t.Fail()
		}

		actual := actualIntf.([]string)
		expected := c.expected

		if len(actual) != len(expected) {
			t.Logf("expected (%d) len %d, got len %d", i, len(expected), len(actual))
			t.Fail()
		}

		for k, v := range expected {
			va := actual[k]
			if va != v {
				t.Logf("expected (%d) %q value %q, got %q", i, k, v, va)
				t.Fail()
			}
		}
	}
}
