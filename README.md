# asp — Automatic Settings Processor

[![Go reference](https://img.shields.io/badge/pkg.go.dev-reference-007D9C?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/jaredreisinger/asp)
![Go version](https://img.shields.io/github/go-mod/go-version/jaredreisinger/asp?logo=go&logoColor=white&style=flat-square)
[![GitHub build](https://img.shields.io/github/actions/workflow/status/jaredreisinger/asp/pipeline.yaml?branch=main&logo=github&style=flat-square)](https://github.com/jaredreisinger/asp/actions/workflows/pipeline.yaml)
[![Codecov](https://img.shields.io/codecov/c/github/jaredreisinger/asp?logo=codecov&label=codedov&style=flat-square)](https://codecov.io/gh/JaredReisinger/asp)
[![License](https://img.shields.io/github/license/jaredreisinger/asp?style=flat-square)](https://github.com/JaredReisinger/asp/blob/main/LICENSE)

`asp`, the Automatic Settings Provider, an opinionated companion for `viper` and `cobra`.

- [Getting started](https://github.com/JaredReisinger/asp/blob/main/docs/01-getting-started.md)
- [Config processing](https://github.com/JaredReisinger/asp/blob/main/docs/02-config-processing.md)
- [Defaults](https://github.com/JaredReisinger/asp/blob/main/docs/03-defaults.md)
- [Options](https://github.com/JaredReisinger/asp/blob/main/docs/04-options.md)
- [Config tags](https://github.com/JaredReisinger/asp/blob/main/docs/05-config-tags.md)

## Why does this exist?

The `cobra` package provides excellent command-line flag functionality, and `viper` provides a rich configuration store and environment variable binding… _but_… there’s a lot of boilerplate and redundant code if you want to achieve the nirvana of CLI flags and environment variables _and_ configuration file support for _**all**_ flags/settings in your application. The `asp` package attempts to reduce this boilerplate by capturing it all from your “canonical” configure structure definition.

The goals of `asp` are to

1. reduce the redundant boilerplate by concisely defining all of the necessary information in the config struct itself;

2. encourage good practices by ensuring that _every_ option has config, command-line _and_ environment variable representation;

3. avoid possible typos that using string-based configuration lookups can cause—Go can’t tell that `viper.Get("sommeSetting")` is misspelled at compile time—but it _can_ tell that `config.sommeSetting` is invalid if the struct defines the member as `config.someSetting`.
