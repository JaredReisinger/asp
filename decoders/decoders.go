package decoders

// We want to encourage the best behavior by making it as easy as possible to
// provide all of the setting variants:
//   * config file (the config struct itself)
//   * command line - long and short versions
//   * environment variable(s)
// ... and also description info.
//
// for example:
//   type Config struct {
//     Host string `asp:"host,h,HOST,The host to use."`
//   }

import (
	// "builtin"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
)

// Note that even for the implementations that extend from the default
// mapstructure one, we now *always* use the [mapstructure.DecodeHookFuncValue]
// interface, rather than only using type or kind.  This has the advantage of
// being the "newer" style, and also allows us to refactor out some common test
// cases.

// StringToTime is similar to [mapstructure.StringToTimeHookFunc], but it
// supports a few magical string constants ("" is zero time, and "now", "utc",
// and "local" yield the current time), and assumes RFC3399Nano layout.
func StringToTime() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to time? %v --> %v", from, to)
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf(time.Time{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to time!! %v --> %v (%q)", from, to, from.String())
		return timeConv(from.String())
	}
}

// StringToByteSlice decodes a hex-encoded string as a series of bytes, useful
// for binary tokens, etc.
func StringToByteSlice() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to byte slice? %v", from.Interface())
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf([]byte{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to byte slice! %v, %v", from.Interface(), to.Interface())

		return hex.DecodeString(from.String())
	}
}

// StringToMapStringInt decodes a "keyOne=1,keyTwo=2"-style string into a map.
// The input string may optionally be surrounded by square brackets
// ("[keyOne=1,keyTwo=2]").
func StringToMapStringInt() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to map[string]int? %v", from.Interface())
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf(map[string]int{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to map[string]int! %v, %v", from.Interface(), to.Interface())
		// dest := to.Interface().(map[string]int)
		dest := make(map[string]int)
		entries := getListEntries(from.String(), ",")
		for _, entry := range entries {
			// log.Printf("converting %q", entry)
			keyVal := strings.SplitN(entry, "=", 2)
			if len(keyVal) != 2 {
				return nil, fmt.Errorf("unexpected map entry %q", entry)
			}
			val, err := strconv.Atoi(keyVal[1])
			if err != nil {
				return nil, err
			}
			dest[keyVal[0]] = val
		}

		return dest, nil
	}
}

// StringToMapStringString decodes a "keyOne=val1,keyTwo=val2"-style string into
// a map. The input string may optionally be surrounded by square brackets
// ("[keyOne=val1,keyTwo=val2]").
func StringToMapStringString() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to map[string]string? %v", from.Interface())
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf(map[string]string{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to map[string]string! %v, %v", from.Interface(), to.Interface())
		// dest := to.Interface().(map[string]string)
		dest := make(map[string]string)
		entries := getListEntries(from.String(), ",")
		for _, entry := range entries {
			// log.Printf("converting %q", entry)
			keyVal := strings.SplitN(entry, "=", 2)
			if len(keyVal) != 2 {
				return nil, fmt.Errorf("unexpected map entry %q", entry)
			}
			dest[keyVal[0]] = keyVal[1]
		}

		return dest, nil
	}
}

// StringToSlice is similar to [mapstructure.StringToSliceHookFunc], but can
// also handle input wrapped in square brackets, which sometimes happens during
// CLI flag serialization.
func StringToSlice(sep string) mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf([]string{}) {
			return from.Interface(), nil
		}

		return getListEntries(from.String(), sep), nil
	}
}

// getListEntries is a helper for the string to map/slice decoders, which all
// need to check for enclosing "[]" and then split.
func getListEntries(s string, sep string) []string {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		s = strings.TrimSuffix(strings.TrimPrefix(s, "["), "]")
	}

	// strings.Split() will result in a 1-element slice containing the empty
	// string.  This isn't really what we want... we want an empty slice!
	if s == "" {
		return []string{}
	}

	return strings.Split(s, sep)
}
