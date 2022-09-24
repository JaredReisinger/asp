package asp

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		t.Logf("unexpected error: %v", err)
		t.Fail()
	}

	if baseType != reflect.TypeOf(TestConfig{}) {
		t.Fail()
	}

	// TODO: check viper/cobra settings?
}
