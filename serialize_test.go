package asp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type withSensitive struct {
	Public  string
	Private string `asp.sensitive:"true"`
	Ignored bool   `asp:"-"` // "-" fields are ignored!
	ignored bool   // unexported fields are ignored!
}

func TestSerializeFlags(t *testing.T) {
	empty := processTestConfig{}
	firstHalf := processTestConfig{
		Time:          time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC),
		Duration:      30 * time.Second,
		DurationSlice: []time.Duration{30 * time.Second, 5 * time.Minute},
		Bool:          true,
		Int:           -7,
		Uint:          6,
		String:        "something",
		BoolSlice:     []bool{false, true},
	}
	lastHalf := processTestConfig{
		IntSlice:        []int{1, -2, 3},
		UintSlice:       []uint{1, 2, 3},
		ByteSlice:       []byte{1, 2, 3},
		StringSlice:     []string{"one", "two", "three"},
		MapStringInt:    map[string]int{"one": 1, "two": 2},
		MapStringString: map[string]string{"key1": "value1", "key2": "value2"},

		Nested: Nested{
			Dummy: 2,
		},

		AnonymousEmbedded: AnonymousEmbedded{
			EmbeddedInt: 8,
		},
	}

	sensitiveEmpty := withSensitive{"public", "", false, false}
	sensitiveSet := withSensitive{"", "private", false, false}

	cases := map[string]struct {
		input     any
		omitEmpty bool
		expected  string
	}{
		"empty omit": {empty, true, ``},

		"empty all": {empty, false, `--time "" --duration "" --duration-slice "" --bool "false" --int "" --uint "" --string "" --bool-slice "" --int-slice "" --uint-slice "" --byte-slice "" --string-slice "" --map-string-int "" --map-string-string "" --nested-dummy "" --embedded-int ""`},

		"first omit": {firstHalf, true, `--time "2000-01-02T03:04:05.000000006Z" --duration "30s" --duration-slice "30s,5m0s" --bool "true" --int "-7" --uint "6" --string "something" --bool-slice "false,true"`},

		"first all": {firstHalf, false, `--time "2000-01-02T03:04:05.000000006Z" --duration "30s" --duration-slice "30s,5m0s" --bool "true" --int "-7" --uint "6" --string "something" --bool-slice "false,true" --int-slice "" --uint-slice "" --byte-slice "" --string-slice "" --map-string-int "" --map-string-string "" --nested-dummy "" --embedded-int ""`},

		"last omit": {lastHalf, true, `--int-slice "1,-2,3" --uint-slice "1,2,3" --byte-slice "010203" --string-slice "one,two,three" --map-string-int "one=1,two=2" --map-string-string "key1=value1,key2=value2" --nested-dummy "2" --embedded-int "8"`},

		"last all": {lastHalf, false, `--time "" --duration "" --duration-slice "" --bool "false" --int "" --uint "" --string "" --bool-slice "" --int-slice "1,-2,3" --uint-slice "1,2,3" --byte-slice "010203" --string-slice "one,two,three" --map-string-int "one=1,two=2" --map-string-string "key1=value1,key2=value2" --nested-dummy "2" --embedded-int "8"`},

		"sensitive-empty omit": {sensitiveEmpty, true, `--public "public"`},

		"sensitive-empty all": {sensitiveEmpty, false, `--public "public" --private ""`},

		"sensitive-set omit": {sensitiveSet, true, `--private [REDACTED]`},

		"sensitive-set all": {sensitiveSet, false, `--public "" --private [REDACTED]`},
	}

	for k, v := range cases {
		input, omitEmpty, expected := v.input, v.omitEmpty, v.expected
		t.Run(k, func(t *testing.T) {
			t.Parallel()

			s, err := SerializeFlags(input, omitEmpty)
			assert.NoError(t, err)
			assert.Equal(t, expected, s)
		})
	}
}

func TestSerializeFlagsNonStruct(t *testing.T) {
	s, err := SerializeFlags(1, false)
	assert.ErrorIs(t, err, ErrConfigMustBeStruct)
	assert.Equal(t, "", s)
}

func TestSerializeFlagsUnsupportedType(t *testing.T) {
	type unsupported struct {
		Rune rune
	}

	s, err := SerializeFlags(unsupported{}, false)
	assert.ErrorIs(t, err, ErrConfigFieldUnsupported)
	assert.Equal(t, "", s)

	s, err = SerializeFlags(struct{ Inner unsupported }{}, false)
	assert.ErrorIs(t, err, ErrConfigFieldUnsupported)
	assert.Equal(t, "", s)
}
