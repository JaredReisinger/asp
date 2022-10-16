package asp

import (
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
		WithoutConfigFlag[TestConfig],
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
		WithDefaultConfigName[TestConfig]("DEFAULT_CONFIG"),
	)
	if err != nil {
		t.Logf("got error: %v", err)
		t.FailNow()
	}

	aspActual := a.(*asp[TestConfig])

	expect(t, 0, "defaultCfgName", "DEFAULT_CONFIG", aspActual.defaultCfgName)
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
