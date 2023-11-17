# Defaults

The [_Getting started_](01-getting-started.md) example lost the default values provided in the original implementation. To get them back, we provide them the value passed to `asp.Attach()`:

> The default `--config` and `--help` flags are omitted for brevity.

```go
asp.Attach(rootCmd, rootConfig{Author: "YOUR NAME", License: "apache"})
```

```
Flags:
      --author string         sets the author value (env: APP_AUTHOR) (default "YOUR NAME")
      --license string        sets the license value (env: APP_LICENSE) (default "apache")
      --project-base string   sets the project base value (env: APP_PROJECTBASE)
      --use-viper             sets the use viper value (env: APP_USEVIPER)
```

If you go back and look at the original implementation, there were actually _two_ defaults for the “author” value: the flags used `"YOUR NAME"`, but the viper bindings used `"NAME HERE <EMAIL ADDRESS>"`. Asp only has one place to specify a default.

To help keep code readable, and to put the defaults closer to the definition of the configuration type, define a variable:

```go
type rootConfig struct {
	ProjectBase string
	Author      string
	License     string
	UseViper    bool
}

var defaults = rootConfig{
	Author:  "YOUR NAME",
	License: "apache",
}
```

and then later:

```go
asp.Attach(rootCmd, defaults)
```
