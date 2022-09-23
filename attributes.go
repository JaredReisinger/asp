package asp

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// There's a precedence list for which values are used: first, the
// attribute-specific tag (`asp.long`, asp.desc`) is used if present--and an
// explicit empty string (`asp.long:""`) can be used to cancel/disable that
// attribute. Next, that component of the general comma-separated `asp` tag is
// used, but note that an empty or missing component (`asp:"long,,"`) does *not*
// cancel/disable the attribute.  Finally, the default calculated value is used
// as a fallback if the attribute hasn't been canceled.
type tagKind string

const (
	tagAll   = tagKind("asp")
	tagLong  = tagKind("asp.long")
	tagShort = tagKind("asp.short")
	tagEnv   = tagKind("asp.env")
	tagDesc  = tagKind("asp.desc")
)

type allIndex int

const (
	allLong allIndex = iota
	allShort
	allEnv
	allDesc

	allMax
)

// parentCanonical and parentEnv come "pre-suffixed"
func getAttributes(f reflect.StructField,
	parentCanonical string, parentEnv string) (
	canonicalName string,
	attrLong string, attrShort string, attrEnv string, attrDesc string) {

	canonicalName = fmt.Sprintf("%s%s", parentCanonical, f.Name)

	// note that we use "just uppercase" for the name, not
	// strcase.ToScreamingSnake(), because convention for environment variables
	// is to use '_' delimeters only at hierarchy boundaries
	envName := fmt.Sprintf("%s%s", parentEnv, strings.ToUpper(f.Name))

	// get attribute values...
	attrLong = getAttribute(f, tagLong, allLong, strcase.ToKebab(canonicalName))
	attrShort = getAttribute(f, tagShort, allShort, "")
	attrEnv = getAttribute(f, tagEnv, allEnv, envName)
	attrDesc = getAttribute(f, tagDesc, allDesc, fmt.Sprintf("sets the %s value", canonicalName))

	return
}

func getAttribute(f reflect.StructField, k tagKind, i allIndex, fallback string) string {
	// we end up calling Tag.Get(tagAll) multiple times... a little overhead,
	// but calling before-hand requires the caller to understand the internals,
	// *and* we need to call Tag.Get() for the specific tag anyway
	attr, ok := f.Tag.Lookup(string(k))
	if !ok {
		all := strings.SplitN(f.Tag.Get(string(tagAll)), ",", int(allMax))
		if len(all) > int(i) {
			attr = strings.TrimSpace(all[i])
		}

		if attr == "" {
			attr = fallback
		}
	}

	return attr
}
