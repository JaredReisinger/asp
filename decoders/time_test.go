package decoders

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeValue(t *testing.T) {
	origTimeNow := timeNow
	defer func() { timeNow = origTimeNow }()

	nowStr := "2022-01-02T03:04:05.666-07:00"
	now, _ := time.Parse(time.RFC3339Nano, nowStr)

	timeNow = func() time.Time {
		return now
	}

	cases := map[string]time.Time{
		"":      {},
		"now":   now,
		"utc":   now.UTC(),
		"local": now.Local(),
	}

	for k, v := range cases {
		input := k
		expected := v
		t.Run(input, func(t *testing.T) {
			// t.Parallel()
			v := NewTimeValue()

			err := v.Set(input)
			assert.NoError(t, err)
			assert.Equal(t, "time", v.Type())
			assert.Equal(t, expected, time.Time(*v))
			if input == "" {
				assert.Equal(t, "", v.String())
			} else {
				assert.Equal(t, expected.Format(time.RFC3339Nano), v.String())
			}
		})
	}

	// v := NewTimeValue()

	// input := "2022-01-02T03:04:05.666-07:00"

	// err := v.Set(input)
	// assert.NoError(t, err)

	// expected, err := time.Parse(time.RFC3339Nano, input)
	// assert.NoError(t, err)

	// assert.Equal(t, expected, time.Time(*v))
	// assert.Equal(t, "time", v.Type())
	// assert.Equal(t, input, v.String())
}

func TestTimeConv(t *testing.T) {
	// t.Parallel()
	origTimeNow := timeNow
	defer func() { timeNow = origTimeNow }()

	nowStr := "2022-01-02T03:04:05.666-07:00"
	now, _ := time.Parse(time.RFC3339Nano, nowStr)

	timeNow = func() time.Time {
		return now
	}

	cases := map[string]time.Time{
		"":      {},
		"now":   now,
		"utc":   now.UTC(),
		"local": now.Local(),
	}

	for k, v := range cases {
		input := k
		expected := v
		t.Run(input, func(t *testing.T) {
			// t.Parallel()
			intf, err := timeConv(input)
			assert.NoError(t, err)
			assert.Equal(t, expected, intf.(time.Time))
			// round-trip to string...

		})
	}
}
