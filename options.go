package asp

// Option represents an option to the [asp.Attach] method.
type Option func(*aspBase) error

// WithDefaultConfigName specifies the filename to use when searching for a
// config file.  This is typically the same as the app name.  If no default
// config name is given, *no* config file will be loaded by default, and the
// `--config` command-line flag *must* be given to use a config file.
func WithDefaultConfigName(cfgName string) Option {
	return func(a *aspBase) error {
		a.defaultCfgName = cfgName
		return nil
	}
}

// WithEnvPrefix specifies the prefix to use with environment variables.  If not
// passed to Attach(), the prefix "APP_" is assumed.
func WithEnvPrefix(prefix string) Option {
	return func(a *aspBase) error {
		a.envPrefix = prefix
		return nil
	}
}

// WithConfigFlag adds a `--config cfgFile` flag to the command being attached.
// Note that there is *not* an environment variable or config setting that
// mirrors this CLI-only flag.  This is set by default.
func WithConfigFlag(a *aspBase) error {
	a.withConfigFlag = true
	return nil
}

// WithoutConfigFlag prevents the addition a `--config cfgFile` flag to the
// command being attached.
func WithoutConfigFlag(a *aspBase) error {
	a.withConfigFlag = false
	return nil
}
