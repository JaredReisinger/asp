package asp

import (
	"errors"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type aspTestConfig struct {
	String   string
	Time     time.Time
	Duration time.Duration
	Bool     bool
	Int      int
}

var defaultConfig = aspTestConfig{}

func TestAttach(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(cmd, defaultConfig)
	assert.NoError(t, err)

	aspActual := a.(*asp[aspTestConfig])

	assert.Equal(t, "APP", aspActual.envPrefix)
	assert.Equal(t, true, aspActual.withConfigFlag)
	assert.Equal(t, "", aspActual.defaultCfgName)

	assert.Equal(t, cmd, a.Command())
	assert.NotNil(t, a.Viper())

	// check debug output...
	// a.Debug()
}

func TestAttachWithoutConfigFlag(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(
		cmd, defaultConfig,
		WithoutConfigFlag,
	)
	assert.NoError(t, err)

	aspActual := a.(*asp[aspTestConfig])

	assert.Equal(t, false, aspActual.withConfigFlag)
}

func testConfigFiles(t *testing.T, aspActual *asp[aspTestConfig]) {
	t.Helper()

	var err error

	aspActual.cfgFile = "asp_test_config.yaml"
	_, err = aspActual.Config()
	assert.NoError(t, err)

	aspActual.cfgFile = "asp_test_config_bad.yaml"
	_, err = aspActual.Config()
	assert.Error(t, err)

	aspActual.cfgFile = "asp_test_config_missing.yaml"
	_, err = aspActual.Config()
	assert.Error(t, err)
}

func TestAttachWithConfigFlag(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(
		cmd, defaultConfig,
		WithConfigFlag,
	)
	assert.NoError(t, err)

	aspActual := a.(*asp[aspTestConfig])

	assert.Equal(t, true, aspActual.withConfigFlag)

	// see what happens without a --config value...
	_, err = aspActual.Config()
	assert.NoError(t, err)

	testConfigFiles(t, aspActual)
}

func TestAttachWithDefaultConfigName(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(
		cmd, defaultConfig,
		WithDefaultConfigName("DEFAULT_CONFIG"),
	)
	assert.NoError(t, err)

	aspActual := a.(*asp[aspTestConfig])

	assert.Equal(t, "DEFAULT_CONFIG", aspActual.defaultCfgName)

	// see what happens without a --config value...
	_, err = aspActual.Config()
	assert.NoError(t, err)

	testConfigFiles(t, aspActual)
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
	assert.ErrorIs(t, err, bogusError)
	assert.Nil(t, a)
}

func TestAttachWithUnsupportedConfig(t *testing.T) {
	cmd := &cobra.Command{}

	badConfig := struct {
		BadMember *int // we don't support pointer members!
	}{}

	_, err := Attach(cmd, badConfig)
	assert.ErrorIs(t, err, ErrConfigFieldUnsupported)
}

func TestConfigResults(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(cmd, defaultConfig)
	assert.NoError(t, err)

	cfg, err := a.Config()
	assert.NoError(t, err)

	// TODO: need more/better testing here!
	assert.Equal(t, "", cfg.String)
}
