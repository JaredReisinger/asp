package asp

import (
	"testing"
	"time"
)

func TestTimeValue(t *testing.T) {
	v := &timeValue{}

	input := "2022-01-02T03:04:05.666-07:00"

	err := v.Set(input)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	expected, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if !time.Time(*v).Equal(expected) {
		t.Logf(
			"expected %q (%d), got %q (%d)",
			expected, expected.UnixNano(),
			v, time.Time(*v).UnixNano())
		t.Fail()
	}

	if v.Type() != "time" {
		t.Fail()
	}

	if v.String() != input {
		t.Fail()
	}
}

func TestTimeConv(t *testing.T) {
	origTimeNow := timeNow
	defer func() { timeNow = origTimeNow }()

	now, _ := time.Parse(time.RFC3339Nano, "2022-01-02T03:04:05.6666-07:00")

	timeNow = func() time.Time {
		return now
	}

	cases := []struct {
		input    string
		expected time.Time
	}{
		{input: "", expected: time.Time{}},
		{input: "now", expected: now},
		{input: "utc", expected: now.UTC()},
		{input: "local", expected: now.Local()},
	}

	for i, c := range cases {
		v, err := timeConv(c.input)
		if err != nil {
			t.Logf("%d unexpected error: %v", i, err)
			t.Fail()
		}

		if !c.expected.Equal(v.(time.Time)) {
			t.Logf("case %d expected %q, got %q", i, c.expected, v)
			t.Fail()
		}
	}
}
