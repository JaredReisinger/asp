package asp

// We want to encourage the best behavior by making it as easy as possible to
// provide all of the setting variants:
//   * config file (the config struct itself)
//   * command line - long and short versions
//   * environment variable(s)
// ... and also description info.
//
// for example:
//   type Config struct {
//     Host string `asp:"host,h,APP_HOST,The host to use."`
//   }

import (
	// "builtin"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

// betterStringToTime handles empty strings as zero time
func betterStringToTime() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (
		interface{}, error) {
		// log.Printf("attempting to convert string to time? %v --> %v", f, t)
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		// log.Printf("attempting to convert string to time!! %v --> %v (%q)", f, t, data.(string))
		// Convert it by parsing
		return timeConv(data.(string))
	}
}

func stringToByteSlice() mapstructure.DecodeHookFuncValue {
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

func stringToMapStringString() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to map[string]string? %v", from.Interface())
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf(map[string]string{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to map[string]string! %v, %v", from.Interface(), to.Interface())

		// dest := to.Interface().(map[string]string)
		dest := make(map[string]string)
		entries := strings.Split(from.String(), ",")
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

// betterStringToSlice improves on mapstructure's StringToSliceHookFunc
// by checking for a wrapping "[" and "]" which sometimes happens during flag
// serialization.
func betterStringToSlice(sep string) mapstructure.DecodeHookFuncKind {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (
		interface{}, error) {
		if f != reflect.String || t != reflect.Slice {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		// check for "[]" around the string...
		if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
			raw = strings.TrimSuffix(strings.TrimPrefix(raw, "["), "]")
		}

		return strings.Split(raw, sep), nil
	}
}
