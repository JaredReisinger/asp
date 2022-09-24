package asp

import "time"

// helper for time.Time values...
const timeLayout = time.RFC3339Nano

type timeValue time.Time

var (
	timeNow = time.Now // alias to enable testing
)

func newTimeValue(val time.Time, p *time.Time) *timeValue {
	*p = val
	return (*timeValue)(p)
}

func (d *timeValue) Set(s string) error {
	v, err := timeConv(s)

	*d = timeValue(v.(time.Time))
	return err
}

func (d *timeValue) Type() string {
	return "time"
}

func (d *timeValue) String() string {
	// return (*time.Time)(d).String()
	t := (*time.Time)(d)
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func timeConv(sval string) (interface{}, error) {
	// Handle some special time constants...
	switch sval {
	case "":
		return time.Time{}, nil
	case "now":
		return timeNow(), nil
	case "utc":
		return timeNow().UTC(), nil
	case "local":
		return timeNow().Local(), nil
	}

	return time.Parse(timeLayout, sval)
}
