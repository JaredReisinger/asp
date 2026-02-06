package asp

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type processTestConfig struct {
	Time            time.Time
	Duration        time.Duration
	DurationSlice   []time.Duration
	Bool            bool
	Int             int
	Uint            uint
	String          string
	BoolSlice       []bool
	IntSlice        []int
	UintSlice       []uint
	ByteSlice       []byte
	StringSlice     []string
	MapStringInt    map[string]int
	MapStringString map[string]string
	Nested          Nested
	AnonymousEmbedded
}

type Nested struct {
	Dummy int
}

type AnonymousEmbedded struct {
	EmbeddedInt int
}

func newBase(t *testing.T) *aspBase {
	t.Helper()
	return &aspBase{
		// config: config,
		// envPrefix:      "APP",
		// withConfigFlag: true,
		cmd: &cobra.Command{},
		vip: viper.New(),
	}
}

func TestProcessStruct(t *testing.T) {
	a := newBase(t)

	err := a.processStruct(processTestConfig{})
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(processTestConfig{}), a.baseType)

	// TODO: check all viper/cobra settings?
	names := []string{
		"time",
		"duration",
		"duration-slice",
		"bool",
		"int",
		"uint",
		"string",
		"bool-slice",
		"int-slice",
		"uint-slice",
		"byte-slice",
		"string-slice",
		"map-string-int",
		"map-string-string",
		"nested-dummy",
		"embedded-int",
	}
	for _, n := range names {
		f := a.cmd.PersistentFlags().Lookup(n)
		assert.NotNil(t, f, n)
	}
}

func TestProcessStructInnerErrors(t *testing.T) {
	a := newBase(t)

	// check for expected errors...
	err := a.processStructInner(nil, attrs{})
	assert.ErrorIs(t, err, ErrConfigMustBeStruct)

	cfgWithBadMember := struct {
		BadMember *int
	}{}

	err = a.processStructInner(cfgWithBadMember, attrs{})
	assert.ErrorIs(t, err, ErrConfigFieldUnsupported)

	cfgWithBadNestedMember := struct {
		BadNest struct {
			BadMember *int
		}
	}{}

	err = a.processStructInner(cfgWithBadNestedMember, attrs{})
	assert.ErrorIs(t, err, ErrConfigFieldUnsupported)
}

func TestProcessStructInnerAttributes(t *testing.T) {
	t.Parallel()

	a := newBase(t)

	// test-specific config
	type Config struct {
		One   string `asp.long:"long"`
		Two   string `asp:"combined"`
		Three string `asp:"both" asp.long:"preferred"`

		// nested (new for 0.2)
		Four struct {
			Five string `asp.long:"inner-long"`
		} `asp.long:"outer-long"`
	}

	err := a.processStructInner(Config{}, attrs{})
	assert.NoError(t, err)

	flags := a.cmd.PersistentFlags()

	// for each item, we have the original name (should not be a flag), the new name, the short name
	cases := map[string]string{
		"one":   "long",
		"two":   "combined",
		"three": "preferred",
		// "four": "outer-long",
		"five": "outer-long-inner-long",
	}

	for k, v := range cases {
		name := k
		long := v
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Nil(t, flags.Lookup(name))
			assert.NotNil(t, flags.Lookup(long))
		})
	}
}

// // I can't seem to test the "template does not parse" edge case.
// func TestProcessBadDescription(t *testing.T) {
// 	a := newBase(t)

// 	err := a.processStruct(struct {
// 		Bad int `asp.description:"{{if}}"`
// 	}{})
// 	assert.NoError(t, err)
// }
