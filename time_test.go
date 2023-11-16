package asp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeValue(t *testing.T) {
	v := &timeValue{}

	input := "2022-01-02T03:04:05.666-07:00"

	err := v.Set(input)
	assert.NoError(t, err)

	expected, err := time.Parse(time.RFC3339Nano, input)
	assert.NoError(t, err)

	assert.Equal(t, expected, time.Time(*v))
	assert.Equal(t, "time", v.Type())
	assert.Equal(t, input, v.String())
}

func TestTimeConv(t *testing.T) {
	// t.Parallel()
	origTimeNow := timeNow
	defer func() { timeNow = origTimeNow }()

	now, _ := time.Parse(time.RFC3339Nano, "2022-01-02T03:04:05.6666-07:00")

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
		})
	}
}
