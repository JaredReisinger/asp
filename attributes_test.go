package asp

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAttributes(t *testing.T) {
	t.Parallel()

	defaultDesc := "sets the {{delimited .Name ' '}} value"
	// fields: fieldName, tag, name, long, short, env, desc
	cases := map[string][7]string{
		// "Ex" for "Example" is enough to catch case changes
		"no tag":         {"Ex", ``, "Ex", "ex", "", "EX", defaultDesc},
		"long only":      {"Ex", `asp.long:"replaced"`, "Ex", "replaced", "", "EX", defaultDesc},
		"short only":     {"Ex", `asp.short:"r"`, "Ex", "ex", "r", "EX", defaultDesc},
		"env only":       {"Ex", `asp.env:"REPLACED"`, "Ex", "ex", "", "REPLACED", defaultDesc},
		"desc only":      {"Ex", `asp.desc:"Replaced"`, "Ex", "ex", "", "EX", "Replaced"},
		"all long only":  {"Ex", `asp:"replaced"`, "Ex", "replaced", "", "EX", defaultDesc},
		"all short only": {"Ex", `asp:",r"`, "Ex", "ex", "r", "EX", defaultDesc},
		"all env only":   {"Ex", `asp:",,REPLACED"`, "Ex", "ex", "", "REPLACED", defaultDesc},
		"all desc only":  {"Ex", `asp:",,,Replaced"`, "Ex", "ex", "", "EX", "Replaced"},
		"all all":        {"Ex", `asp:"replaced,r,REPLACED,Replaced"`, "Ex", "replaced", "r", "REPLACED", "Replaced"},

		// multi-word name
		"multi no tag":         {"ExAm", ``, "ExAm", "ex-am", "", "EXAM", defaultDesc},
		"multi long only":      {"ExAm", `asp.long:"replaced"`, "ExAm", "replaced", "", "EXAM", defaultDesc},
		"multi short only":     {"ExAm", `asp.short:"r"`, "ExAm", "ex-am", "r", "EXAM", defaultDesc},
		"multi env only":       {"ExAm", `asp.env:"REPLACED"`, "ExAm", "ex-am", "", "REPLACED", defaultDesc},
		"multi desc only":      {"ExAm", `asp.desc:"Replaced"`, "ExAm", "ex-am", "", "EXAM", "Replaced"},
		"multi all long only":  {"ExAm", `asp:"replaced"`, "ExAm", "replaced", "", "EXAM", defaultDesc},
		"multi all short only": {"ExAm", `asp:",r"`, "ExAm", "ex-am", "r", "EXAM", defaultDesc},
		"multi all env only":   {"ExAm", `asp:",,REPLACED"`, "ExAm", "ex-am", "", "REPLACED", defaultDesc},
		"multi all desc only":  {"ExAm", `asp:",,,Replaced"`, "ExAm", "ex-am", "", "EXAM", "Replaced"},

		// combinations of tags
		"all and long":  {"Ex", `asp:"r,r,R,R" asp.long:"override"`, "Ex", "override", "r", "R", "R"},
		"all and short": {"Ex", `asp:"r,r,R,R" asp.short:"o"`, "Ex", "r", "o", "R", "R"},
		"all and env":   {"Ex", `asp:"r,r,R,R" asp.env:"OVERRIDE"`, "Ex", "r", "r", "OVERRIDE", "R"},
		"all and desc":  {"Ex", `asp:"r,r,R,R" asp.desc:"Override"`, "Ex", "r", "r", "R", "Override"},
	}

	for k, v := range cases {
		field, tag, name, long, short, env, desc := v[0], v[1], v[2], v[3], v[4], v[5], v[6]

		t.Run(k, func(t *testing.T) {
			t.Parallel()

			f := reflect.StructField{
				Name: field,
				Tag:  reflect.StructTag(tag),
			}

			attrs := getAttributes(f)
			// t.Errorf("%q", k)
			assert.Equal(t, name, attrs.name)
			assert.Equal(t, long, attrs.long)
			assert.Equal(t, short, attrs.short)
			assert.Equal(t, env, attrs.env)
			assert.Equal(t, desc, attrs.desc)
		})
	}
}

func TestJoinField(t *testing.T) {
	t.Parallel()

	// fields: prefix, suffix, sep, expected
	cases := map[string][4]string{
		"all empty":     {"", "", "", ""},
		"prefix only":   {"P", "", "", "P"},
		"suffix only":   {"", "S", "", "S"},
		"sep only":      {"", "", "_", ""},
		"prefix suffix": {"P", "S", "", "PS"},
		"prefix sep":    {"P", "", "_", "P"},
		"suffix sep":    {"", "S", "_", "S"},
		"everything":    {"P", "S", "_", "P_S"},
	}

	for k, v := range cases {
		prefix, suffix, sep, expected := v[0], v[1], v[2], v[3]

		t.Run(k, func(t *testing.T) {
			t.Parallel()

			actual := joinField(prefix, suffix, sep)
			assert.Equal(t, expected, actual)
		})
	}
}

func TestAttrsJoin(t *testing.T) {
	t.Parallel()

	attrsNone := attrs{"", "", "", "", ""}
	attrsAll := attrs{"Name", "long", "s", "ENV", "desc"}

	cases := map[string][3]attrs{
		"none none": {attrsNone, attrsNone, attrs{"", "", "", "", ""}},
		"none all":  {attrsNone, attrsAll, attrs{"Name", "long", "s", "ENV", "desc"}},
		"all none":  {attrsAll, attrsNone, attrs{"Name", "long", "", "ENV", ""}},
		"all all":   {attrsAll, attrsAll, attrs{"Name.Name", "long-long", "s", "ENV_ENV", "desc"}},
	}

	for k, v := range cases {
		parent, child, expected := v[0], v[1], v[2]

		t.Run(k, func(t *testing.T) {
			t.Parallel()

			actual := parent.join(child)
			assert.Equal(t, expected.name, actual.name)
			assert.Equal(t, expected.long, actual.long)
			assert.Equal(t, expected.short, actual.short)
			assert.Equal(t, expected.env, actual.env)
			assert.Equal(t, expected.desc, actual.desc)
		})
	}

}
