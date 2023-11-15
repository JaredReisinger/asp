package asp

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AnonymousEmbedded struct {
	EmbeddedInt int
}

type TestConfig struct {
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

	Nested struct {
		Dummy int
	}

	AnonymousEmbedded
}

func newBase(t *testing.T) *aspBase {
	t.Helper()
	return &aspBase{
		// config: config,
		// envPrefix:      "APP_",
		// withConfigFlag: true,
		cmd: &cobra.Command{},
		vip: viper.New(),
	}
}

func TestProcessStruct(t *testing.T) {
	a := newBase(t)

	err := a.processStruct(TestConfig{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if a.baseType != reflect.TypeOf(TestConfig{}) {
		t.Error("unexpected type")
	}

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
		if f == nil {
			t.Errorf("expected flag for %q, got nil", n)
		}
	}
}

func TestProcessStructInnerErrors(t *testing.T) {
	a := newBase(t)

	// check for expected errors...
	err := a.processStructInner(nil, attrs{})
	if !errors.Is(err, ErrConfigMustBeStruct) {
		t.Error("expected error")
	}

	cfgWithBadMember := struct {
		BadMember *int
	}{}

	err = a.processStructInner(cfgWithBadMember, attrs{})
	if !errors.Is(err, ErrConfigFieldUnsupported) {
		t.Error("expected error")
	}

	cfgWithBadNestedMember := struct {
		BadNest struct {
			BadMember *int
		}
	}{}

	err = a.processStructInner(cfgWithBadNestedMember, attrs{})
	if !errors.Is(err, ErrConfigFieldUnsupported) {
		t.Error("expected error")
	}
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
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

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

			if flags.Lookup(name) != nil {
				t.Errorf("did not expect flag for %q, got one", name)
			}

			if flags.Lookup(long) == nil {
				t.Errorf("expect flag for tag name %q, got nil", long)
			}

		})
	}
}
