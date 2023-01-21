package asp

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func expectBool(t *testing.T, i int, label string, expected bool, actual bool) {
	if actual != expected {
		t.Logf("case %d (%s): expected %t, got %t", i, label, expected, actual)
		t.Fail()
	}
}

var defaultConfig = TestConfig{}

func TestAttach(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(cmd, defaultConfig)
	if err != nil {
		t.Logf("got error: %v", err)
		t.FailNow()
	}

	aspActual := a.(*asp[TestConfig])

	expect(t, 0, "envPrefix", "APP_", aspActual.envPrefix)
	expectBool(t, 0, "withConfigFlag", aspActual.withConfigFlag, true)
	expect(t, 0, "defaultCfgName", "", aspActual.defaultCfgName)

	if a.Command() != cmd {
		t.Logf("expected same Commmand")
		t.Fail()
	}

	if a.Viper() == nil {
		t.Logf("expected Viper")
		t.Fail()
	}

	// check debug output...
	a.Debug()
}

func TestAttachWithoutConfigFlag(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(
		cmd, defaultConfig,
		WithoutConfigFlag,
	)
	if err != nil {
		t.Logf("got error: %v", err)
		t.FailNow()
	}

	aspActual := a.(*asp[TestConfig])

	expectBool(t, 0, "withConfigFlag", aspActual.withConfigFlag, false)
}

func TestAttachWithDefaultConfigName(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(
		cmd, defaultConfig,
		WithDefaultConfigName("DEFAULT_CONFIG"),
	)
	if err != nil {
		t.Logf("got error: %v", err)
		t.FailNow()
	}

	aspActual := a.(*asp[TestConfig])

	expect(t, 0, "defaultCfgName", "DEFAULT_CONFIG", aspActual.defaultCfgName)

	// see what happens without a --config value...
	aspActual.Config()

	// also see what happens if the --config value is passed...
	aspActual.cfgFile = "asp_test_config.yaml"
	aspActual.Config()
}

func TestAttachWithBogusOption(t *testing.T) {
	cmd := &cobra.Command{}

	bogusError := errors.New("bogus")
	bogusOption := func(a *aspBase) error {
		return bogusError
	}

	a, err := Attach(
		cmd, defaultConfig,
		bogusOption,
	)

	if !errors.Is(err, bogusError) {
		t.Fail()
	}

	if a != nil {
		t.Fail()
	}
}

func TestAttachWithUnsupportedConfig(t *testing.T) {
	cmd := &cobra.Command{}

	badConfig := struct {
		BadMember *int // we don't support pointer members!
	}{}

	_, err := Attach(cmd, badConfig)
	if !errors.Is(err, ErrConfigFieldUnsupported) {
		t.Fail()
	}

}

func TestConfigResults(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(cmd, defaultConfig)
	if err != nil {
		t.Logf("got error: %v", err)
		t.FailNow()
	}

	cfg := a.Config()

	// TODO: need more/better testing here!
	expect(t, 0, "misc", "", cfg.String)
}
