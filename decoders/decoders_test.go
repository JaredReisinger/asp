package decoders

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBetterStringToTime(t *testing.T) {
	t.Parallel()

	fn := BetterStringToTime()

	isoDatetime := "2022-01-02T03:04:05.666-07:00"

	expected, err := time.Parse(time.RFC3339Nano, isoDatetime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	cases := map[string]struct {
		input    interface{}
		outType  reflect.Type
		err      error
		expected interface{}
	}{
		"ISO datetime":          {isoDatetime, reflect.TypeOf(expected), nil, expected},
		"to string passthrough": {isoDatetime, reflect.TypeOf(isoDatetime), nil, isoDatetime},
		"from int passthrough":  {1, reflect.TypeOf(time.Time{}), nil, 1},
	}

	for k, v := range cases {
		input := v.input
		typ := v.outType
		expectedErr := v.err
		expectedResult := v.expected
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			actualIntf, err := fn(reflect.TypeOf(input), typ, input)
			if expectedErr == assert.AnError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, expectedErr, err)
			}

			if err == nil {
				assert.Equal(t, expectedResult, actualIntf)
			}
		})
	}
}

func TestStringToByteSlice(t *testing.T) {
	t.Parallel()

	fn := StringToByteSlice()

	cases := map[string]struct {
		input    interface{}
		outValue reflect.Value
		err      error
		expected interface{}
	}{
		"deadbeef":              {"deadbeef", reflect.ValueOf([]byte{}), nil, []byte{0xde, 0xad, 0xbe, 0xef}},
		"to string passthrough": {"deadbeef", reflect.ValueOf(""), nil, "deadbeef"},
		"from int passthrough":  {1, reflect.ValueOf([]byte{}), nil, 1},
	}

	for k, v := range cases {
		input := v.input
		outValue := v.outValue
		expectedErr := v.err
		expectedResult := v.expected
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			actualIntf, err := fn(reflect.ValueOf(input), outValue)
			if expectedErr == assert.AnError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, expectedErr, err)
			}

			if err == nil {
				assert.Equal(t, expectedResult, actualIntf)
			}
		})
	}
}

func TestStringToMapStringInt(t *testing.T) {
	t.Parallel()

	fn := StringToMapStringInt()

	expected := map[string]int{
		"key1": 1,
		"key2": 2,
	}

	cases := map[string]struct {
		err    error
		result map[string]int
	}{
		"key1=1,key2=2":   {nil, expected},
		"[key1=1,key2=2]": {nil, expected},
		"key1":            {assert.AnError, nil},
		"key1=FAIL":       {assert.AnError, nil},
	}

	for k, v := range cases {
		input := k
		expectedErr := v.err
		expectedResult := v.result
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			actualIntf, err := fn(reflect.ValueOf(input), reflect.ValueOf(map[string]int{}))
			if expectedErr == assert.AnError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, expectedErr, err)
			}

			if expectedResult != nil {
				actual := actualIntf.(map[string]int)
				assert.Equal(t, expectedResult, actual)
			}
		})
	}
}
func TestStringToMapStringString(t *testing.T) {
	fn := StringToMapStringString()

	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	cases := map[string]struct {
		err    error
		result map[string]string
	}{
		"key1=value1,key2=value2": {nil, expected},
		// "[key1=value1,key2=value2]": {assert.AnError, nil},
		"key1": {assert.AnError, nil},
		// "key1=FAIL":                 {assert.AnError, nil},
	}

	for k, v := range cases {
		input := k
		expectedErr := v.err
		expectedResult := v.result
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			actualIntf, err := fn(reflect.ValueOf(input), reflect.ValueOf(map[string]string{}))
			if expectedErr == assert.AnError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, expectedErr, err)
			}

			if expectedResult != nil {
				actual := actualIntf.(map[string]string)
				assert.Equal(t, expectedResult, actual)
			}
		})
	}
}

func TestBetterStringToSlice(t *testing.T) {
	t.Parallel()

	fn := BetterStringToSlice(",")

	expected := []string{"one", "two"}

	cases := map[string]struct {
		err    error
		result []string
	}{
		"one,two":   {nil, expected},
		"[one,two]": {nil, expected},
		"":          {nil, []string{}},
	}

	for k, v := range cases {
		input := k
		expectedErr := v.err
		expectedResult := v.result
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			actualIntf, err := fn(reflect.String, reflect.Slice, input)
			if expectedErr == assert.AnError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, expectedErr, err)
			}

			if expectedResult != nil {
				actual := actualIntf.([]string)
				assert.Equal(t, expectedResult, actual)
			}
		})
	}
}
