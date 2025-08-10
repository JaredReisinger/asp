# Options

The `asp.Attach()` method also takes a series of `asp.Option` values which can customize behavior:

| option                                         | behavior                                                                                                                                                   |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `asp.WithConfigFlag` / `asp.WithoutConfigFlag` | turns on/off the `--config` flag (on by default)                                                                                                           |
| `asp.WithDecodeHook(`_hook_`)`                 | overrides the default unmarhsaling hook to add support for custom types                                                                                    |
| `asp.WithDefaultConfigName(`_name_`)`          | tells asp (viper) to look for config files named _name_ ([in many common formats](https://github.com/spf13/viper?tab=readme-ov-file#reading-config-files)) |
| `asp.WithEnvPrefix(`_prefix_`)`                | overrides the default `APP` prefix for generated environment variable names                                                                                |

The env-prefix and default config name options are the ones most likely to be used. To change asp to prefix environment variables with `MYAPP`, and look for a “myapp” config file, use an `asp.Attach()` call like:

```go
asp.Attach(cmd, config{}, asp.WithEnvPrefix("MYAPP"), asp.WithDefaultConfigName("myapp"))
```

### WithConfigFlag / WithoutConfigFlag

Turns on  (or off) the `--config` flag.

### WithDecodeHook

Under the covers, viper uses [mapstructure](https://pkg.go.dev/github.com/go-viper/mapstructure/v2) to decode values into typed structures. (See also [the viper doc](https://github.com/spf13/viper?tab=readme-ov-file#decoding-custom-formats).) You can use `asp.WithDecodeHook()` to provide your own [mapstructure.DecodeHookFunc](https://pkg.go.dev/github.com/go-viper/mapstructure/v2#DecodeHookFunc) to decode additional values. The default asp decoders are exported as well, so that you can leverage/combine them if you want.

### WithDefaultConfigName

Provides a default name for a config file to load.

### WithEnvPrefix

Allows you to provide a value to override the default `APP` environment variable name prefix.
