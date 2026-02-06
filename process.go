package asp

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/jaredreisinger/asp/decoders"
)

var (
	// ErrConfigMustBeStruct indicates that the (default) config value passed to
	// [asp.Attach] must be a struct (or pointer to struct) type.
	ErrConfigMustBeStruct = errors.New("config must be a struct or pointer to struct")

	// ErrConfigFieldUnsupported is returned when a member of the config struct
	// is an unsupported type.
	ErrConfigFieldUnsupported = errors.New("config struct field is of an unsupported type (pointer, array, channel or size-specific number)")
)

// TODO: with the introduction of sprig, we may not need these...
var templateFuncs = template.FuncMap{
	"camel":              strcase.ToCamel,
	"delimited":          strcase.ToDelimited,
	"kebab":              strcase.ToKebab,
	"lowerCamel":         strcase.ToLowerCamel,
	"screamingDelimited": strcase.ToScreamingDelimited,
	"screamingKebab":     strcase.ToScreamingKebab,
	"screamingSnake":     strcase.ToScreamingSnake,
	"snake":              strcase.ToSnake,
	"snakeWithIgnore":    strcase.ToSnakeWithIgnore,
}

// processStruct is the main entrypoint for creating CLI flags and environment
// variable mappings for a configuration struct type.
func (a *aspBase) processStruct(s interface{}) error {
	// Now that we've separated the entrypoint and recursive handler, we can be
	// slightly more specific about the requirements on the incoming
	// type/defaults. (We could insist on a struct value, and not a
	// point-to-struct.) But I don't think there's any *particular* reason to
	// force this.
	err := a.processStructInner(s, attrs{env: strings.TrimRight(a.envPrefix, "_")})
	if err != nil {
		return err
	}

	// We don't have to check the kind of `s`... if it got through
	// processStructInner, the kind is acceptable!
	a.baseType = reflect.Indirect(reflect.ValueOf(s)).Type()
	return nil
}

// processStructInner is the (recursive) workhorse that adds a (sub-)struct config
// into the viper config and cobra command.
func (a *aspBase) processStructInner(s interface{}, parentAttrs attrs) error {
	vip, flags := a.vip, a.cmd.PersistentFlags()

	// log.Printf("initializing struct for: %#v", s)

	// We expect the incoming value to be a struct or a pointer to a struct.
	// Anything else is invalid.
	structVal := reflect.Indirect(reflect.ValueOf(s))
	if structVal.Kind() != reflect.Struct {
		return ErrConfigMustBeStruct
	}

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

		// Special handling for the description: if neither {{.Env}} or
		// {{.NoEnv}} appears in the string, we append a "(env: {{Env}})"
		// suffix. Unless, of course, the description has been explicitly
		// omitted (but why, oh why, would you do that?)
		desc := joinedAttrs.desc
		if desc != "" {
			re := regexp.MustCompile(`\{\{\w*\.(?:No)?Env\w*`)
			if !re.Match([]byte(desc)) {
				desc = fmt.Sprintf("%s (env: {{.Env}})", desc)
			}

			// We allow the description attribute to include template values
			// that we fill in based on the calculated name, env, etc.
			tmpl, err := template.New("desc").Funcs(sprig.TxtFuncMap()).Funcs(templateFuncs).Parse(desc)
			if err != nil {
				return errors.WithMessagef(err, "on %s (%q)", f.Name, desc)
			}
			descBuilder := &strings.Builder{}
			err = tmpl.Execute(descBuilder, map[string]string{
				"Name":  joinedAttrs.name,
				"Long":  joinedAttrs.long,
				"Short": joinedAttrs.short,
				"Env":   joinedAttrs.env,
				"NoEnv": "",
				// special context naming (from ASE work)
				"ParentName": parentAttrs.name,
			})
			if err != nil {
				return err
			}
			desc = descBuilder.String()
		}

		// use shortened names purely for concision...
		l, s, d := joinedAttrs.long, joinedAttrs.short, desc

		// Rather than setting handled to true in our myriad cases, we default
		// to true, and make sure to set it to false in our default/unhandled
		// cases.
		handled := true
		addBindings := true
		addEnvBinding := true

		// log.Printf("handling field %q : anonymous? %v, index: %v", canonicalName, f.Anonymous, f.Index)

		// This is some very repetitive code!  The flags helpers are type-safe
		// (but not "fluent"), and thus the long, short, and desc values have to
		// get specified again and again.  Perhaps there's a better library for
		// this, or an opportunity for a new one?
		//
		// ``` flag := flags.AsP(l, s, d)
		// flag.IntP(v2.Int()) ```

		// We *could* use flags.VarP() directly, but that doesn't have the
		// flags.Value-creating helpers, and we need those.

		fieldVal := structVal.FieldByIndex(f.Index)
		intf := fieldVal.Interface()

		// switch it := intf.(type) {
		// default:
		// 	log.Printf("interface type is: %v\n\n%v\n\n", it, v2.InterfaceData())
		// }

		// There are special-case types that we handle up-front, falling back to
		// low-level "kinds" only if we need to...
		switch val := intf.(type) {
		case time.Time:
			// create our own time-parsing flag
			flags.VarP(decoders.NewTimeValue(), l, s, d)

		case time.Duration:
			flags.DurationP(l, s, val, d)

		case []time.Duration:
			flags.DurationSliceP(l, s, val, d)

		case bool:
			flags.BoolP(l, s, val, d)

		case int:
			flags.IntP(l, s, val, d)

		case uint:
			flags.UintP(l, s, val, d)

		case string:
			// FUTURE: should we handle "rich" parsing for things like IP
			// addresses, Durations, etc?
			flags.StringP(l, s, val, d)

		case []bool:
			flags.BoolSliceP(l, s, val, d)

		case []int:
			flags.IntSliceP(l, s, val, d)

		case []uint:
			flags.UintSliceP(l, s, val, d)

		// pFlags supports []byte, but the parsing gets confused?
		// maybe that's viper?

		case []byte:
			// This is really []byte!... we'd double-check, but at runtime,
			// all we see is []uint8.
			flags.BytesHexP(l, s, val, d)

		case []string:
			flags.StringSliceP(l, s, val, d)

		case map[string]int:
			flags.StringToIntP(l, s, val, d)

		case map[string]string:
			flags.StringToStringP(l, s, val, d)

		default:
			if f.Type.Kind() == reflect.Struct {
				recursiveAttrs := joinedAttrs

				// need to think about whether
				if f.Anonymous {
					recursiveAttrs = parentAttrs
				}

				err := a.processStructInner(intf, recursiveAttrs)
				if err != nil {
					return err
				}

				addBindings = false // prevent default flag/config additions!
			} else {
				handled = false
			}
		}

		if !handled {
			log.Printf("unsupported type? %q %#v", f.Type.Kind(), f)
			return ErrConfigFieldUnsupported
		}

		if addBindings {
			// log.Printf("%q, %v, CLI: %q / %q, env: %q, desc: %q",
			// 	canonicalName, f.Type.Kind(),
			// 	attrLong, attrShort, attrEnv, attrDesc)

			// Start pushing into viper?  Note that we're going to need to handle
			// parent paths pretty quickly!
			vip.SetDefault(joinedAttrs.name, intf)

			err := vip.BindPFlag(joinedAttrs.name, flags.Lookup(joinedAttrs.long))
			if err != nil {
				return err
			}

			if addEnvBinding {
				err = vip.BindEnv(joinedAttrs.name, joinedAttrs.env)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
