package decoders

import (
	"reflect"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

type decoderCase struct {
	input    interface{}
	output   interface{}
	err      error
	expected interface{}
}

type decoderCases map[string]decoderCase

func runCases(t *testing.T, parallel bool, fn mapstructure.DecodeHookFuncValue, cases decoderCases) {
	t.Helper()
	if parallel {
		t.Parallel()
	}

	for k, v := range cases {
		input := v.input
		output := v.output
		expectedErr := v.err
		expectedResult := v.expected
		t.Run(k, func(t *testing.T) {
			if parallel {
				t.Parallel()
			}
			actualIntf, err := fn(reflect.ValueOf(input), reflect.ValueOf(output))
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

func TestStringToTime(t *testing.T) {
	isoDatetime := "2022-01-02T03:04:05.666-07:00"
	expected, err := time.Parse(time.RFC3339Nano, isoDatetime)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	runCases(t, true, StringToTime(), decoderCases{
		"ISO datetime":          {isoDatetime, time.Time{}, nil, expected},
		"to string passthrough": {isoDatetime, "", nil, isoDatetime},
		"from int passthrough":  {1, time.Time{}, nil, 1},
	})
}

func TestStringToByteSlice(t *testing.T) {
	runCases(t, true, StringToByteSlice(), decoderCases{
		"deadbeef":              {"deadbeef", []byte{}, nil, []byte{0xde, 0xad, 0xbe, 0xef}},
		"to string passthrough": {"deadbeef", "", nil, "deadbeef"},
		"from int passthrough":  {1, []byte{}, nil, 1},
	})
}

func TestStringToMapStringInt(t *testing.T) {
	expected := map[string]int{
		"key1": 1,
		"key2": 2,
	}

	runCases(t, true, StringToMapStringInt(), decoderCases{
		"simple":                {"key1=1,key2=2", map[string]int{}, nil, expected},
		"wrapped":               {"[key1=1,key2=2]", map[string]int{}, nil, expected},
		"syntax err":            {"key1", map[string]int{}, assert.AnError, nil},
		"parse err":             {"key1=FAIL", map[string]int{}, assert.AnError, nil},
		"to string passthrough": {"key1=1,key2=2", "", nil, "key1=1,key2=2"},
		"from int passthrough":  {1, map[string]int{}, nil, 1},
	})
}
func TestStringToMapStringString(t *testing.T) {
	expected := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	runCases(t, true, StringToMapStringString(), decoderCases{
		"simple":                {"key1=value1,key2=value2", map[string]string{}, nil, expected},
		"wrapped":               {"[key1=value1,key2=value2]", map[string]string{}, nil, expected},
		"syntax err":            {"key1", map[string]string{}, assert.AnError, nil},
		"to string passthrough": {"key1=value1,key2=value2", "", nil, "key1=value1,key2=value2"},
		"from int passthrough":  {1, map[string]string{}, nil, 1},
	})
}

func TestStringToSlice(t *testing.T) {
	expected := []string{"one", "two"}

	runCases(t, true, StringToSlice(","), decoderCases{
		"simple":                {"one,two", []string{}, nil, expected},
		"wrapped":               {"[one,two]", []string{}, nil, expected},
		"empty":                 {"", []string{}, nil, []string{}},
		"to string passthrough": {"one,two", "", nil, "one,two"},
		"from int passthrough":  {1, []string{}, nil, 1},
	})

	// t.Parallel()

	// fn := StringToSlice(",")

	// cases := map[string]struct {
	// 	err    error
	// 	result []string
	// }{
	// 	"one,two":   {nil, expected},
	// 	"[one,two]": {nil, expected},
	// 	"":          {nil, []string{}},
	// }

	// for k, v := range cases {
	// 	input := k
	// 	expectedErr := v.err
	// 	expectedResult := v.result
	// 	t.Run(input, func(t *testing.T) {
	// 		t.Parallel()
	// 		actualIntf, err := fn(reflect.String, reflect.Slice, input)
	// 		if expectedErr == assert.AnError {
	// 			assert.Error(t, err)
	// 		} else {
	// 			assert.Equal(t, expectedErr, err)
	// 		}

	// 		if expectedResult != nil {
	// 			actual := actualIntf.([]string)
	// 			assert.Equal(t, expectedResult, actual)
	// 		}
	// 	})
	// }
}
