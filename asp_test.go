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

func TestAttachInstance(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := AttachInstance(cmd, defaultConfig)
	assert.NoError(t, err)

	aspActual := a.(*asp[aspTestConfig])

	assert.Equal(t, "APP", aspActual.envPrefix)
	assert.Equal(t, true, aspActual.withConfigFlag)
	assert.Equal(t, "", aspActual.defaultCfgName)

	assert.Equal(t, cmd, a.Command())
	assert.NotNil(t, a.Viper())
}

func TestAttachInstanceWithoutConfigFlag(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := AttachInstance(
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

func TestAttachInstanceWithConfigFlag(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := AttachInstance(
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

func TestAttachInstanceWithDefaultConfigName(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := AttachInstance(
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

func TestAttachInstanceWithBogusOption(t *testing.T) {
	cmd := &cobra.Command{}

	bogusError := errors.New("bogus")
	bogusOption := func(a *aspBase) error {
		return bogusError
	}

	a, err := AttachInstance(
		cmd, defaultConfig,
		bogusOption,
	)
	assert.ErrorIs(t, err, bogusError)
	assert.Nil(t, a)
}

func TestAttachInstanceWithUnsupportedConfig(t *testing.T) {
	cmd := &cobra.Command{}

	badConfig := struct {
		BadMember *int // we don't support pointer members!
	}{}

	_, err := AttachInstance(cmd, badConfig)
	assert.ErrorIs(t, err, ErrConfigFieldUnsupported)
}

func TestConfigResults(t *testing.T) {
	cmd := &cobra.Command{}

	a, err := AttachInstance(cmd, defaultConfig)
	assert.NoError(t, err)

	cfg, err := a.Config()
	assert.NoError(t, err)

	// TODO: need more/better testing here!
	assert.Equal(t, "", cfg.String)
}

func TestAttachedCommand(t *testing.T) {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := Get[aspTestConfig](cmd)
			assert.NoError(t, err)
			assert.Equal(t, "test", cfg.String)
		},
	}

	err := Attach(cmd, defaultConfig)
	assert.NoError(t, err)

	cmd.SetArgs([]string{"--string", "test"})
	cmd.Execute()
}

func TestAttachedCommandWrongType(t *testing.T) {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := Get[struct{}](cmd)
			assert.ErrorIs(t, err, ErrConfigTypeMismatch)
			assert.Nil(t, cfg)
		},
	}

	err := Attach(cmd, defaultConfig)
	assert.NoError(t, err)

	cmd.Execute()
}
