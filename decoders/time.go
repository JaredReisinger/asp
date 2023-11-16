package decoders

import (
	"time"
)

// helpers for time.Time values...
type timeValue time.Time

const timeLayout = time.RFC3339Nano

var (
	timeNow = time.Now // alias to enable testing
)

// NewTimeValue returns an interface around [time.Time] that supports [cobra]
// and the [pflag.Value] interface.
func NewTimeValue() *timeValue {
	return &timeValue{}
}

// Set uses our string-to-time conversion/decoding logic to set the [time.Time]
// value.
func (d *timeValue) Set(s string) error {
	v, err := timeConv(s)

	*d = timeValue(v.(time.Time))
	return err
}

// Type returns "time" for the type of the value.
func (d *timeValue) Type() string {
	return "time"
}

// String renders the time value as a string.
func (d *timeValue) String() string {
	// return (*time.Time)(d).String()
	t := (*time.Time)(d)
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func timeConv(s string) (interface{}, error) {
	// handle some special time constants...
	switch s {
	case "":
		return time.Time{}, nil
	case "now":
		return timeNow(), nil
	case "utc":
		return timeNow().UTC(), nil
	case "local":
		return timeNow().Local(), nil
	}

	return time.Parse(timeLayout, s)
}
