/*
[asp], the Automatic Settings Provider, an opinionated companion for [viper] and [cobra].

  - [Getting started]
  - [Config processing]
  - [Defaults]
  - [Options]
  - [Config tags]

## Why does this exist?

The [cobra] package provides excellent command-line flag functionality, and [viper] provides a rich configuration store and environment variable binding… _but_… there’s a lot of boilerplate and redundant code if you want to achieve the nirvana of CLI flags and environment variables _and_ configuration file support for _**all**_ flags/settings in your application. The [asp] package attempts to reduce this boilerplate by capturing it all from your “canonical” configure structure definition.

### The goals of [asp] are to

 1. reduce the redundant boilerplate by concisely defining all of the necessary information in the config struct itself;

 2. encourage good practices by ensuring that _every_ option has config, command-line _and_ environment variable representation;

 3. avoid possible typos that using string-based configuration lookups can cause—Go can’t tell that `viper.Get("SommeSetting")` is misspelled at compile time—but it _can_ tell that `config.SommeSetting` is invalid if the struct defines the member as `config.SomeSetting`.

[Getting started]: https://github.com/JaredReisinger/asp/blob/main/docs/01-getting-started.md
[Config processing]: https://github.com/JaredReisinger/asp/blob/main/docs/02-config-processing.md
[Defaults]: https://github.com/JaredReisinger/asp/blob/main/docs/03-defaults.md
[Options]: https://github.com/JaredReisinger/asp/blob/main/docs/04-options.md
[Config tags]: https://github.com/JaredReisinger/asp/blob/main/docs/05-config-tags.md
*/
package asp

// The meta-viper (https://github.com/carlosvin/meta-viper) project does
// something similar, but I wanted more control over exactly how the
// command-line, env-var, and config file structure was defined.
