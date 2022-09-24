package asp

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, i int, label string, expected string, actual string) {
	if actual != expected {
		t.Logf("case %d (%s): expected %q, got %q", i, label, expected, actual)
		t.Fail()
	}
}

func TestGetAttributes(t *testing.T) {
	// For excessivly concise testing, we don't define a struct... we use a
	// fixed-size array.  (We can get away with this because all of the values are
	// strings!)  The values are: [fieldName, fieldTag, parentCanonical, parentEnv,
	// expectedName, expectedLong, expectedShort, expectedEnv, expectedDesc]
	cases := [][9]string{
		// basic cases
		{"Field", ``, "", "", "Field", "field", "", "FIELD", "sets the Field value"},
		{"Field", `asp.short:"f"`, "", "", "Field", "field", "f", "FIELD", "sets the Field value"},
		{"Field", `asp.long:"other"`, "", "", "Field", "other", "", "FIELD", "sets the Field value"},
		{"Field", `asp.env:"OTHER"`, "", "", "Field", "field", "", "OTHER", "sets the Field value"},
		{"Field", `asp.desc:"OTHER"`, "", "", "Field", "field", "", "FIELD", "OTHER"},
		{"Field", `asp:""`, "", "", "Field", "field", "", "FIELD", "sets the Field value"},
		{"Field", `asp:"LONG,S,ENV,DESC with spaces"`, "", "", "Field", "LONG", "S", "ENV", "DESC with spaces"},

		// parent/nesting
		{"Field", ``, "Parent.", "", "Parent.Field", "parent-field", "", "FIELD", "sets the Parent.Field value"},
		{"Field", `asp.short:"f"`, "Parent.", "", "Parent.Field", "parent-field", "f", "FIELD", "sets the Parent.Field value"},
		{"Field", `asp.long:"other"`, "Parent.", "", "Parent.Field", "other", "", "FIELD", "sets the Parent.Field value"},
		{"Field", `asp.env:"OTHER"`, "Parent.", "", "Parent.Field", "parent-field", "", "OTHER", "sets the Parent.Field value"},
		{"Field", `asp.desc:"OTHER"`, "Parent.", "", "Parent.Field", "parent-field", "", "FIELD", "OTHER"},
		{"Field", `asp:""`, "Parent.", "", "Parent.Field", "parent-field", "", "FIELD", "sets the Parent.Field value"},
		{"Field", `asp:"LONG,S,ENV,DESC with spaces"`, "Parent.", "", "Parent.Field", "LONG", "S", "ENV", "DESC with spaces"},

		// env prefix
		{"Field", ``, "", "ENV_", "Field", "field", "", "ENV_FIELD", "sets the Field value"},
		{"Field", `asp.short:"f"`, "", "ENV_", "Field", "field", "f", "ENV_FIELD", "sets the Field value"},
		{"Field", `asp.long:"other"`, "", "ENV_", "Field", "other", "", "ENV_FIELD", "sets the Field value"},
		{"Field", `asp.env:"OTHER"`, "", "ENV_", "Field", "field", "", "OTHER", "sets the Field value"},
		{"Field", `asp.desc:"OTHER"`, "", "ENV_", "Field", "field", "", "ENV_FIELD", "OTHER"},
		{"Field", `asp:""`, "", "ENV_", "Field", "field", "", "ENV_FIELD", "sets the Field value"},
		{"Field", `asp:"LONG,S,ENV,DESC with spaces"`, "", "ENV_", "Field", "LONG", "S", "ENV", "DESC with spaces"},
	}

	for i, c := range cases {
		// there's no easy "spread the array across variables", so we have to do it by
		// hand
		fieldName, fieldTag, parentCanonical, parentEnv,
			expectedName, expectedLong, expectedShort, expectedEnv, expectedDesc :=
			c[0], c[1], c[2], c[3], c[4], c[5], c[6], c[7], c[8]

		f := reflect.StructField{
			Name: fieldName,
			Tag:  reflect.StructTag(fieldTag),
		}

		canonicalName, attrLong, attrShort, attrEnv, attrDesc := getAttributes(f, parentCanonical, parentEnv)

		expect(t, i, "name", expectedName, canonicalName)
		expect(t, i, "long", expectedLong, attrLong)
		expect(t, i, "short", expectedShort, attrShort)
		expect(t, i, "env", expectedEnv, attrEnv)
		expect(t, i, "desc", expectedDesc, attrDesc)
	}
}
