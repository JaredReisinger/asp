package asp

import (
	"cmp"
	"encoding/hex"
	"fmt"
	"log"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jaredreisinger/asp/decoders"
)

// SerializeFlags returns the CLI flags that would re-create the given
// configuration values, with the exception of any redacted sensitive values
// (which are replaced with [REDACTED] in the returned value).
func SerializeFlags[T Config](cfg T, omitEmpty bool) (string, error) {
	return serializeStruct(cfg, omitEmpty)
}

// serializeStruct is the main entrypoint for serializing CLI flags for logging
// purposes. The omitEmpty flag specifies whether empty/zero/default values
// should be omitted from the serialization.
func serializeStruct(s interface{}, omitEmpty bool) (string, error) {
	// Now that we've separated the entrypoint and recursive handler, we can be
	// slightly more specific about the requirements on the incoming
	// type/defaults. (We could insist on a struct value, and not a
	// point-to-struct.) But I don't think there's any *particular* reason to
	// force this.
	str, err := serializeStructInner(s, omitEmpty, attrs{sensitive: false})
	if err != nil {
		return "", err
	}

	// We don't have to check the kind of `s`... if it got through
	// serializeStructInner, the kind is acceptable!
	// a.baseType = reflect.Indirect(reflect.ValueOf(s)).Type()
	return str, nil
}

// serializeStructInner is the (recursive) workhorse that serializes a
// (sub-)struct config; the logic is very similar to [processStructInner].
func serializeStructInner(s interface{}, omitEmpty bool, parentAttrs attrs) (string, error) {
	// log.Printf("initializing struct for: %#v", s)

	// We expect the incoming value to be a struct or a pointer to a struct.
	// Anything else is invalid.
	structVal := reflect.Indirect(reflect.ValueOf(s))
	if structVal.Kind() != reflect.Struct {
		return "", ErrConfigMustBeStruct
	}

	str := &strings.Builder{}

	fields := reflect.VisibleFields(structVal.Type())
	// log.Printf("fields: %#v", fields)

	for _, f := range fields {
		// Skip any unexported fields!
		if !f.IsExported() {
			continue
		}

		// We deal with anonymous (embedded) structs by *not* updating the
		// parentCanonical/parentEnv strings when recursing.  We also need to
		// *not* attempt to process the mirrored sub-elements directly, because
		// we need the canonical structure to get serialized properly.  We can
		// tell if a field is a mirrored embedded field because its "Index"
		// value isn't a length-1 array, it's length 2+.
		if len(f.Index) > 1 {
			continue
		}

		childAttrs := getAttributes(f)
		joinedAttrs := parentAttrs.join(childAttrs)

		// Rather than setting handled to true in our myriad cases, we default
		// to true, and make sure to set it to false in our default/unhandled
		// cases.
		handled := true

		// log.Printf("handling field %q : anonymous? %v, index: %v", canonicalName, f.Anonymous, f.Index)

		fieldVal := structVal.FieldByIndex(f.Index)
		intf := fieldVal.Interface()

		// switch it := intf.(type) {
		// default:
		// 	log.Printf("interface type is: %v\n\n%v\n\n", it, v2.InterfaceData())
		// }

		fieldStr := ""

		// There are special-case types that we handle up-front, falling back to
		// low-level "kinds" only if we need to...
		switch val := intf.(type) {
		case time.Time:
			if !val.IsZero() {
				fieldStr = val.Format(decoders.TimeLayout)
			}

		case time.Duration:
			if val != 0 {
				fieldStr = val.String()
			}

		case []time.Duration:
			if len(val) > 0 {
				fieldStr = strings.Join(mapSlice(
					val,
					func(d time.Duration) string { return d.String() },
				), ",")
			}

		case bool:
			fieldStr = fmt.Sprintf("%t", val)
			if !val && omitEmpty {
				fieldStr = ""
			}

		case int:
			if val != 0 {
				fieldStr = strconv.FormatInt(int64(val), 10)
			}

		case uint:
			if val != 0 {
				fieldStr = strconv.FormatUint(uint64(val), 10)
			}

		case string:
			if val != "" {
				fieldStr = val
			}

		case []bool:
			if len(val) > 0 {
				fieldStr = strings.Join(mapSlice(
					val,
					func(b bool) string { return fmt.Sprintf("%t", b) },
				), ",")
			}

		case []int:
			if len(val) > 0 {
				fieldStr = strings.Join(mapSlice(
					val,
					func(i int) string { return strconv.FormatInt(int64(i), 10) },
				), ",")
			}

		case []uint:
			if len(val) > 0 {
				fieldStr = strings.Join(mapSlice(
					val,
					func(u uint) string { return strconv.FormatUint(uint64(u), 10) },
				), ",")
			}

		case []byte:
			fieldStr = hex.EncodeToString(val)

		case []string:
			if len(val) > 0 {
				fieldStr = strings.Join(val, ",")
			}

		case map[string]int:
			if len(val) > 0 {
				fieldStr = strings.Join(mapMapToSlice(
					val,
					func(k string, v int) string { return fmt.Sprintf("%s=%d", k, v) },
				), ",")
			}

		case map[string]string:
			if len(val) > 0 {
				fieldStr = strings.Join(mapMapToSlice(
					val,
					func(k string, v string) string { return fmt.Sprintf("%s=%s", k, v) },
				), ",")
			}

		default:
			if f.Type.Kind() == reflect.Struct {
				recursiveAttrs := joinedAttrs

				if f.Anonymous {
					recursiveAttrs = parentAttrs
				}

				childStr, err := serializeStructInner(intf, omitEmpty, recursiveAttrs)
				if err != nil {
					return "", err
				}

				if str.Len() > 0 && len(childStr) > 0 {
					str.WriteString(" ")
				}
				str.WriteString(childStr)

				// unlike with processing, we need to skip the "end of switch"
				// handling
				continue
			} else {
				handled = false
			}
		}

		if !handled {
			log.Printf("unsupported type? %q %#v", f.Type.Kind(), f)
			return ",", ErrConfigFieldUnsupported
		}

		if fieldStr != "" || !omitEmpty {
			if str.Len() > 0 {
				str.WriteString(" ")
			}
			formattedValue := fmt.Sprintf("%q", fieldStr)
			if joinedAttrs.sensitive && fieldStr != "" {
				formattedValue = "[REDACTED]"
			}
			fmt.Fprintf(str, "--%s %s", joinedAttrs.long, formattedValue)
		}

	}

	return str.String(), nil
}

func mapSlice[S ~[]E, E any, X any](s S, fn func(E) X) []X {
	x := make([]X, 0, len(s))

	for _, e := range s {
		x = append(x, fn(e))
	}

	return x
}

func mapMapToSlice[M map[K]V, K cmp.Ordered, V any, X any](m M, fn func(k K, v V) X) []X {
	x := make([]X, 0, len(m))

	for _, k := range slices.Sorted(maps.Keys(m)) {
		x = append(x, fn(k, m[k]))
	}

	return x
}
