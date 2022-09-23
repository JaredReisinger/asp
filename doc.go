// asp, the Automatic Settings Provider, an opinionated companion for viper and
// cobra.
//
// Why?  I really like viper and cobra for building 12-factor style
// applications, but there's a lot of overhead incurred (in lines-of-code) in
// creating the config, command-line, and environment variable settings for each
// option.  The goal of `asp` is the (a) reduce the redundant boilerplate by
// concisely defining all of the necessary information in the config struct
// itself, (b) to encourage good practices by ensuring that *every* option has
// config, command-line, *and* environment variable representation, and (c) to
// avoid possible typos that using string-based configuration lookups can
// cause--Go can't tell that `viper.Get("sommeSetting")` is misspelled at
// compile time... but it *can* tell that `config.sommeSetting` is invalid if
// the struct defines the member as `someSetting`.
//
// This is done by driving the command-line and environment variable settings
// from a tagged struct that's used to define the config file format.
//
//     type Config struct {
//       Host string `asp:""`
//     }
package asp

// The meta-viper (https://github.com/carlosvin/meta-viper) project does
// something similar, but I wanted more control over exactly how the
// command-line, env-var, and config file structure was defined.
