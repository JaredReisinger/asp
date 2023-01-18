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

func TestProcessStruct(t *testing.T) {
	cmd := &cobra.Command{}
	vip := viper.New()

	a := &asp[TestConfig]{
		// config: config,
		envPrefix:      "APP_",
		withConfigFlag: true,
		vip:            vip,
		cmd:            cmd,
	}

	baseType, err := a.processStruct(TestConfig{}, "", a.envPrefix)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if baseType != reflect.TypeOf(TestConfig{}) {
		t.Error("unexpected type")
	}

	// TODO: check viper/cobra settings?
}

func TestProcessStructErrors(t *testing.T) {
	cmd := &cobra.Command{}
	vip := viper.New()

	a := &asp[TestConfig]{
		// config: config,
		envPrefix:      "APP_",
		withConfigFlag: true,
		vip:            vip,
		cmd:            cmd,
	}

	// check for expected errors...
	_, err := a.processStruct(nil, "", "")
	if !errors.Is(err, ErrConfigMustBeStruct) {
		t.Error("expected error")
	}

	cfgWithBadMember := struct {
		BadMember *int
	}{}

	_, err = a.processStruct(cfgWithBadMember, "", "")
	if !errors.Is(err, ErrConfigFieldUnsupported) {
		t.Error("expected error")
	}

	cfgWithBadNestedMember := struct {
		BadNest struct {
			BadMember *int
		}
	}{}

	_, err = a.processStruct(cfgWithBadNestedMember, "", "")
	if !errors.Is(err, ErrConfigFieldUnsupported) {
		t.Error("expected error")
	}
}
