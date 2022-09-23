package asp

type Option[T IncomingConfig] func(*asp[T]) error

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
