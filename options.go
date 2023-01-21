package asp

// Option represents an option to the [asp.Attach] method.
type Option[T IncomingConfig] func(*asp[T]) error

// WithDefaultConfigName specifies the filename to use when searching for a
// config file.  This is typically the same as the app name.  If no default
// config name is given, *no* config file will be loaded by default, and the
// `--config` command-line flag *must* be given to use a config file.
func WithDefaultConfigName[T IncomingConfig](cfgName string) Option[T] {
	return func(a *asp[T]) error {
		a.defaultCfgName = cfgName
		return nil
	}
}

// WithEnvPrefix specifies the prefix to use with environment variables.  If not
// passed to Attach(), the prefix "APP_" is assumed.
func WithEnvPrefix[T IncomingConfig](prefix string) Option[T] {
	return func(a *asp[T]) error {
		a.envPrefix = prefix
		return nil
	}
}

// WithConfigFlag adds a `--config cfgFile` flag to the command being attached.
// Note that there is *not* an environment variable or config setting that
// mirrors this CLI-only flag.  This is set by default.
func WithConfigFlag[T IncomingConfig](a *asp[T]) error {
	a.withConfigFlag = true
	return nil
}

// WithoutConfigFlag prevents the addition a `--config cfgFile` flag to the
// command being attached.
func WithoutConfigFlag[T IncomingConfig](a *asp[T]) error {
	a.withConfigFlag = false
	return nil
}
