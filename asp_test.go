package asp

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var defaultConfig = TestConfig{}

func TestAttach(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(cmd, defaultConfig)
	assert.NoError(t, err)

	aspActual := a.(*asp[TestConfig])

	assert.Equal(t, "APP_", aspActual.envPrefix)
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

	aspActual := a.(*asp[TestConfig])

	assert.Equal(t, false, aspActual.withConfigFlag)
}

func TestAttachWithDefaultConfigName(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := Attach(
		cmd, defaultConfig,
		WithDefaultConfigName("DEFAULT_CONFIG"),
	)
	assert.NoError(t, err)

	aspActual := a.(*asp[TestConfig])

	assert.Equal(t, "DEFAULT_CONFIG", aspActual.defaultCfgName)

	// see what happens without a --config value...
	_, err = aspActual.Config()
	assert.NoError(t, err)

	// also see what happens if the --config value is passed...
	aspActual.cfgFile = "asp_test_config.yaml"
	_, err = aspActual.Config()
	assert.NoError(t, err)
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
