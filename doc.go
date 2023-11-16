// [asp], the Automatic Settings Provider, an opinionated companion for [viper] and [cobra].
//
// ## Why does this exist?
//
// The [cobra] package provides excellent command-line flag functionality, and [viper] provides a rich configuration store and environment variable binding… _but_… there’s a lot of boilerplate and redundant code if you want to achieve the nirvana of CLI flags and environment variables _and_ configuration file support for _**all**_ flags/settings in your application. The [asp] package attempts to reduce this boilerplate by capturing it all from your “canonical” configure structure definition.
//
// # The goals of [asp] are to
//
//  1. reduce the redundant boilerplate by concisely defining all of the necessary information in the config struct itself;
//
//  2. encourage good practices by ensuring that _every_ option has config, command-line _and_ environment variable representation;
//
//  3. avoid possible typos that using string-based configuration lookups can cause—Go can’t tell that `viper.Get("sommeSetting")` is misspelled at compile time—but it _can_ tell that `config.sommeSetting` is invalid if the struct defines the member as `config.someSetting`.
//
// ## Getting started
//
// Assuming that you have a [cobra.Command]-based tool stubbed out, all you need to do is:
//
//  1. Create a type that defines your configuration syntax.
//
//  2. Call [asp.Attach] with your [cobra.Command] instance and an instance of your configuration struct that contains any default values.
//
//  3. Additionally, because of the way that the [cobra.Command.Run] function is called, you will probably want to inject the returned [asp.Asp] interface into your command context so that you can retrieve it inside of your `Run` implementation.
//
// Here’s a contrived example from [example/main.go] that highlights some of the automatic name processing, supported types, and available overrides (on the `Verbose` member):
//
//	package main
//
//	import (
//		"context"
//		"log"
//
//		"github.com/spf13/cobra"
//		"github.com/jaredreisinger/asp"
//	)
//
//	type Config struct {
//		SomeValue       string
//		SomeFlag        bool
//		ManyNumbers     []int
//		MapStringInt    map[string]int
//		MapStringString map[string]string
//
//		SubSection struct {
//			NamesLikeThis string
//		}
//
//		Verbose bool `asp.short:"v" asp.desc:"get noisy"`
//	}
//
//	func main() {
//		defaults := Config{
//			SomeValue: "DEFAULT STRING!",
//		}
//
//		cmd := &cobra.Command{
//			Run: commandHandler,
//		}
//
//		err := asp.Attach(
//			cmd, defaults,
//			asp.WithDefaultConfigName("asp-example"),
//		)
//		cobra.CheckErr(err)
//
//		err = cmd.Execute()
//		cobra.CheckErr(err)
//	}
//
//	func commandHandler(cmd *cobra.Command, args []string) {
//		// get the config using the asp.Asp instance attached to cmd
//		config, err := asp.Get[Config](cmd)
//		cobra.CheckErr(err)
//
//		log.Printf("got config: %#v", config)
//	}
//
// If you try running this, and use the `–help` flag (thanks, cobra!), you’ll get:
//
//	Usage:
//		[flags]
//
//	Flags:
//		    --config string                        configuration file to load
//		-h, --help                                 help for this command
//		    --many-numbers ints                    sets the ManyNumbers value (or use APP_MANYNUMBERS)
//		    --map-string-int stringToInt           sets the MapStringInt value (or use APP_MAPSTRINGINT) (default [])
//		    --map-string-string stringToString     sets the MapStringString value (or use APP_MAPSTRINGSTRING) (default [])
//		    --some-flag                            sets the SomeFlag value (or use APP_SOMEFLAG)
//		    --some-value string                    sets the SomeValue value (or use APP_SOMEVALUE) (default "DEFAULT STRING!")
//		    --sub-section-names-like-this string   sets the SubSection.NamesLikeThis value (or use APP_SUBSECTION_NAMESLIKETHIS)
//		-v, --verbose                              get noisy (or use APP_VERBOSE)
//
// Simply by calling [asp.Attach()], you’ve gotten the CLI flags, along with environment variable support _and_ config file parsing. The tagging on the `Config.Verbose` member alters the default help string and adds the `-v` shorthand. Try running this tool using variations of flags and environment values (and/or config!), and you will see the resulting config (line-breaks and whitespace added for legibility):
//
//	$ APP_VERBOSE=true go run ./example/main.go --some-value "from the CLI"
//
//	2023/01/17 23:57:44 got config: &main.Config{
//	  SomeValue:        "from the CLI",
//	  SomeFlag:         false,
//	  ManyNumbers:      []int(nil),
//	  MapStringInt:     map[string]int{},
//	  MapStringString:  map[string]string{},
//	  SubSection:       struct {NamesLikeThis string}{
//	    NamesLikeThis:     ""
//	  },
//	  Verbose:          true
//	}
//
// ## Getting fancy…
//
// When processing your configuration struct, [asp] turns `NamesLikeThis` into CLI flags with `--names-like-this`, and environment variables with `APP_NAMESLIKETHIS`. Named (non-anonymous) struct members are handled simiarly, with `SubSection. NamesLikeThis` becoming `--sub-section-names-like-this` and `APP_SUBSECTION_NAMESLIKETHIS`.
//
// You can provide your own values for these with the following tags:
//
//	tag         | meaning
//	----------- | ----------------------------------------------------------------------------------------------------------------
//	`asp.long`  | the long `--some-name` style CLI flag
//	`asp.short` | the short `-n` style CLI flag
//	`asp.env`   | the environment variable (prepended with envPrefix; `APP_` by default)
//	`asp.desc`  | the help text to show for the flag; any given value is automatically suffixed with the environment variable name
//
// You can see an example of the `asp.short` and `asp.desc` tags on the `Verbose` member in the above Getting started section.
//
// > _Technically, the default long-CLI name behavior means that [asp] can create conflicting CLI names for `Foo.BarBaz` and `FooBar.Baz` — they’d both become `--foo-bar-baz` — but in practice this isn’t very likely, and you can always provide your own `asp.long` value to mitigate the problem._
//
// ## Supported types
//
// Many of the flag types supported by [cobra] (by [pflags], really) are supported:
//
//	type                | format (CLI argument and/or environment)
//	------------------- | ------------------------------------------------------------------------------------------------------
//	`bool`              | no argument needed for `true`, environment variables can use `true`, `false`, `1`, `0`, `yes`, or `no`
//	`int`               | the number as a string
//	`uint`              | the number as a string
//	`string`            | the string value
//	`[]int`             | comma-separated numbers
//	`[]uint`            | comma-separated numbers
//	`[]string`          | comma-separated strings; note that an individual string value thus cannot itself contain a comma!
//	`map[string]int`    | comma-separated, equal-delimited string/number pairs (like `"a=1,b=4"`)
//	`map[string]string` | comma-separated, equal-delimited string/string pairs (like `"a=foo,b=bar"`)
//	`time.Time`         | RFC3399Nano format, or the literal `now`
//	`time.Duration`     | ISO8601 duration format
//
// ## Embedded anonymous structs
//
// It’s reasonable to want to compose app configuration out of sub-parts, and embed anonymous structs to make those values transparently available at runtime. This is [asp]’s default behavior with anonymous structs, but there are a few caveats about which you need to be aware:
//
//   - The config, flag, and environment variable names for an anonymous embedded struct _**do not**_ include the name of the anonymous embedded struct itself. If you want to include the struct name, simply don’t make it an anonymous embed, and ignore the rest of this section entirely.
//
//   - When writing the anonymous embedded struct reference, you need to include a [mapstructure] tag to “squash” the members to the parent map for deserialization. It would be ideal if [asp] could somehow default this for you, but it cannot. (I wish it were the default for [mapstructure], but alas, it is not.) For example:
//
//     type CommonFields struct {
//     FirstName string
//     LastName  string
//     }
//
//     type Config struct {
//     CommonFields `mapstructure:",squash"` // <==== this is the needed tag!
//     More         string
//     }
//
//     Without the [mapstructure] “squash” option, the [viper] configuration file values won’t map to the final config object correctly.
//
//   - When you write a config file (in YAML, TOML, or what-have-you), you must write as though the embedded fields exist directly in the parent:
//
//     # config.yaml
//     firstName: John
//     lastName: Doe
//     more: used for an unknown person
//     # --*NOT*--
//     # commonFields:
//     #   firstName: John
//     #   lastName: Doe
//     # more: ...
//
//   - As per standard Go behavior, however, while you will be able to “read” values from your loaded configuration using the embedded struct field shorthand (`config.FirstName`), you _cannot_ programmatically construct your config that way. In this case, for example to create defaults, you will need to provide the embedded struct explicitly:
//
//     var Default = Config{
//     CommonFields: CommonFields{
//     FirstName: "Mia",
//     },
//     }
//
// The examples above will result in:
//
//	--first-name string   sets the FirstName value (or use APP_FIRSTNAME) (default "Mia")
//	--last-name string    sets the LastName value (or use APP_LASTNAME)
//	--more string         sets the More value (or use APP_MORE)
package asp

// The meta-viper (https://github.com/carlosvin/meta-viper) project does
// something similar, but I wanted more control over exactly how the
// command-line, env-var, and config file structure was defined.
