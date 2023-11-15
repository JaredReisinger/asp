package asp

// We want to encourage the best behavior by making it as easy as possible to
// provide all of the setting variants:
//   * config file (the config struct itself)
//   * command line - long and short versions
//   * environment variable(s)
// ... and also description info.
//
// for example:
//   type Config struct {
//     Host string `asp:"host,h,APP_HOST,The host to use."`
//   }

import (
	// "builtin"

	"log" // REVIEW: maybe update to log/slog, go 1.21?
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// IncomingConfig is a placeholder generic type that exists only to allow us to
// define our inner and exposed value as strongly-typed to the originating
// configuration struct.
type IncomingConfig interface {
	interface{}
}

type contextKey struct{}

var ContextKey = contextKey{}

// Asp[T IncomingConfig] is an interface that represents the interface for
// settings/options.  After creating/initializing with a configuration structure
// (with default values), the methods on the interface allow for loading from
// command-line/config/environment, as well as lower-level access to the created
// viper instance and cobra command.  (In most cases these should not be needed,
// though!)
type Asp[T IncomingConfig] interface {
	// Config returns the aggregated configuration values, pulling from CLI
	// flags, environment variables, and implicit or explict config file.
	Config() *T

	// Command provides access to the [cobra.Command] that this instance of
	// [Asp] was attached to, in case additional Command customization is
	// needed.
	Command() *cobra.Command

	// Viper provides access to the [viper.Viper] that was created when this
	// instance of [Asp] was attached to the command, in case additional Viper
	// customization is needed.
	Viper() *viper.Viper

	// Execute(handler func(config T, args []string)) error

	// Debug outputs (via [log.Printf]) some diagnostic information to help
	// verify what settings/configs have been acquired.
	Debug()
}

// Attach adds to `cmd` the command-line arguments, and environment variable and
// configuration file bindings inferred from `config`.  If no `Option`s are
// provided, it defaults to `WithConfigFlag` and `WithEnvPrefix("APP_")`.
func Attach[T IncomingConfig](cmd *cobra.Command, config T, options ...Option) (Asp[T], error) {
	vip := viper.New()

	a := &asp[T]{
		aspBase: aspBase{
			// config: config,
			envPrefix:      "APP_",
			withConfigFlag: true,
			vip:            vip,
			cmd:            cmd,
		},
	}
	// log.Printf("initializing config for: %#v", config)

	var err error

	// handle any/all options...
	for _, opt := range options {
		err = opt(&a.aspBase)
		if err != nil {
			return nil, err
		}
	}

	if a.withConfigFlag {
		cmd.PersistentFlags().StringVar(&a.cfgFile, "config", "", "configuration file to load")
	}

	err = a.processStruct(config)
	if err != nil {
		return nil, err
	}

	if a.defaultCfgName != "" {
		vip.SetConfigName(a.defaultCfgName)

		vip.AddConfigPath(".")
		vip.AddConfigPath("$HOME/.config")
		vip.AddConfigPath("$HOME")
		vip.AddConfigPath("/etc")
	}

	return a, nil
}

// We use aspBase to represent everything *except* the generic-type-specific
// stuff... this allows us to omit type specifiers in things like Option.  In
// theory, we could also define AspBase as a non-generic interface, but we don't
// currently have any use cases where we really need a non-type-specific
// interface exposed.
type aspBase struct {
	defaultCfgName string
	envPrefix      string
	withConfigFlag bool

	vip     *viper.Viper
	cmd     *cobra.Command
	cfgFile string

	baseType reflect.Type
}

// I'm using the generic T to "seed" the type at the time that Attach() is
// called, but it "pollutes" all usage of the asp instance, when the *only* call
// that really benefits from it is asp.Config().  I *really* like not having to
// explicitly provide the type in the .Config() call,
type asp[T IncomingConfig] struct {
	aspBase
}

// func (a *asp[T]) Execute(handler func(config T, args []string)) error {
// 	// Set up run-handler for the cobra command...
// 	a.cmd.Run = func(cmd *cobra.Command, args []string) {
// 		log.Printf("BEFORE (INSIDE): %v", a.vip.AllSettings())
// 		// TODO: unmarshal the settings into the expected config type!
// 		cfgVal := reflect.New(a.baseType)
// 		handler(cfgVal.Interface().(T), args)
// 		log.Printf("AFTER (INSIDE): %v", a.vip.AllSettings())
// 	}

// 	log.Printf("BEFORE: %v", a.vip.AllSettings())

// 	// a.cmd.ParseFlags()
// 	err := a.cmd.Execute()
// 	log.Printf("error? %v", err)

// 	log.Printf("AFTER: %v", a.vip.AllSettings())
// 	return err
// }

func (a *aspBase) Command() *cobra.Command {
	return a.cmd
}

func (a *aspBase) Viper() *viper.Viper {
	return a.vip
}

func (a *aspBase) Debug() {
	log.Printf("asp.Debug: %#v", a.vip.AllSettings())
}

func (a *asp[T]) Config() *T {
	// Before reading the config, check to see if there was a `--config` option
	// that specifies a particular config file!
	expectCfgFile := false

	if a.withConfigFlag && a.cfgFile != "" {
		log.Printf("using config file %q", a.cfgFile)
		// a.vip.SetConfigName(a.cfgFile)
		// a.vip.AddConfigPath(".")
		a.vip.SetConfigFile(a.cfgFile)
		expectCfgFile = true
	}

	val := reflect.New(a.baseType)
	// log.Printf("created config: %+v", val.Interface())
	cfg := val.Interface().(*T)
	// log.Printf("viper settings: %#v", a.vip.AllSettings())
	err := a.vip.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		case *viper.ConfigFileNotFoundError:
			if expectCfgFile {
				log.Fatalf("specified config file %q not found", a.cfgFile)
			} else {
				log.Printf("no config file found... perhaps there are environment variables")
			}
		default:
			log.Fatalf("read config error: (%T) %s", err, err.Error())
		}
	}

	err = a.vip.Unmarshal(
		cfg,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				betterStringToTime(),
				stringToByteSlice(),
				stringToMapStringInt(),
				stringToMapStringString(),
				betterStringToSlice(","),
			)))

	if err != nil {
		log.Fatalf("unmarshal config error: %+v", err)
	}

	// log.Printf("returning merged config: %+v", cfg)
	return cfg
}
