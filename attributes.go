package asp

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// There’s a precedence list for which values are used: first, the
// attribute-specific tag (`asp.long`, `asp.desc`) is used if present—and an
// explicit empty string (`asp.long:""`) can be used to cancel/disable that
// attribute. Otherwise, that component of the general comma-separated `asp` tag
// is used, but note that an empty or missing component (`asp:","` is missing
// the “long” value) does *not* cancel/disable the attribute. Finally, the
// default calculated/reflected value is used as a fallback *if* the attribute
// hasn’t been canceled.

// getAttributes returns the various asp-consumed attributes for the given
// field, based on the field name and any asp-specific tags.
func getAttributes(f reflect.StructField) attrs {
	// pre-fill with defaults from field name (canonicalize?)
	a := attrs{
		name:      f.Name,
		long:      strcase.ToKebab(f.Name),
		short:     "",
		env:       strings.ToUpper(f.Name),
		desc:      "sets the {{delimited .Name ' '}} value",
		sensitive: false,
	}

	// now go through the possible tags and allow them to override
	for _, tag := range attrTags {
		val, ok := f.Tag.Lookup(tag.tagName)
		if !ok {
			continue
		}
		tag.setter(&a, val)
	}

	return a
}

type tagInfo struct {
	tagName string
	setter  func(*attrs, string)
}

// Note that attrTags serves double-duty; it's the list of tags to parse, *and*
// the order of attribute-specific tags is their order inside the high-level
// overall tag... so that we can look up the indexed setter by adding 1.
var attrTags []tagInfo

func init() {
	// has to be set in init to avoid circular use inside attrs.setAll()!
	attrTags = []tagInfo{
		{"asp", (*attrs).setAll},
		{"asp.long", (*attrs).setLong},
		{"asp.short", (*attrs).setShort},
		{"asp.env", (*attrs).setEnv},
		{"asp.desc", (*attrs).setDesc},
		{"asp.sensitive", (*attrs).setSensitive},
	}
}

// The attrs struct holds the collection of attrs parsed from the struct field
// tags.
type attrs struct {
	name      string
	long      string
	short     string
	env       string
	desc      string // *template* string to allow full name to be substituted in
	sensitive bool
}

func (a *attrs) setAll(s string) {
	parts := strings.SplitN(s, ",", len(attrTags)-1)

	for i, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			attrTags[i+1].setter(a, p)
		}
	}
}

func (a *attrs) setLong(s string)      { a.long = s }
func (a *attrs) setShort(s string)     { a.short = s }
func (a *attrs) setEnv(s string)       { a.env = s }
func (a *attrs) setDesc(s string)      { a.desc = s }
func (a *attrs) setSensitive(s string) { a.sensitive = (strings.ToLower(s) == "true") }

// combine builds a new attribute set using the aggregation/combination rules
// for each individual field
func (a *attrs) join(child attrs) attrs {
	return attrs{
		name:      joinField(a.name, child.name, "."),
		long:      joinField(a.long, child.long, "-"),
		short:     child.short, // short flags are *never* joined!
		env:       joinField(a.env, child.env, "_"),
		desc:      child.desc, // descriptions are *never* joined!
		sensitive: a.sensitive || child.sensitive,
	}
}

// joinField is a helper for joining attrs fields.
func joinField(prefix string, suffix string, sep string) string {
	if prefix == "" {
		return suffix
	}

	if suffix == "" {
		return prefix
	}

	return prefix + sep + suffix
}
