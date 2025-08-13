# Config tags

If we look again at the configuration type declaration from [_Getting started_](01-getting-started.md):

> The default `--config` and `--help` flags are omitted for brevity.

```go
type rootConfig struct {
    ProjectBase string
    Author      string
    License     string
    UseViper    bool
}
```

```
Flags:
      --author string         sets the author value (env: APP_AUTHOR)
      --license string        sets the license value (env: APP_LICENSE)
      --project-base string   sets the project base value (env: APP_PROJECTBASE)
      --use-viper             sets the viper value (env: APP_USEVIPER)
```

…we can see that the flag descriptions aren’t quite what the original implementation had, nor are there any short flag names. We can fix these with asp tags:

```go
type rootConfig struct {
    ProjectBase string `asp.short:"b" asp.long:"projectbase" asp.desc:"base project directory eg. github.com/spf13/"`
    Author      string `asp.short:"a" asp.desc:"Author name for copyright attribution"`
    License     string `asp.short:"l" asp.desc:"Name of license for the project"`
    UseViper    bool   `asp.desc:"Use Viper for configuration"`
}
```

```
Flags:
  -a, --author string        Author name for copyright attribution (env: APP_AUTHOR)
  -l, --license string       Name of license for the project (env: APP_LICENSE)
  -b, --projectbase string   base project directory eg. github.com/spf13/ (env: APP_PROJECTBASE)
      --use-viper            Use Viper for configuration (env: APP_USEVIPER)
```

And now we have very good fidelity with the original implementation.

## Tags

| tag                              | meaning                                                                                     |
| -------------------------------- | ------------------------------------------------------------------------------------------- |
| [`asp`](#asp)                    | combination of the other four values, comma-separated in this order: `long,short,env,desc`. |
| [`asp.desc`](#aspdesc)           | help text to show for the flag; (processed as a template)                                   |
| [`asp.env`](#aspenv)             | environment variable (prepended with envPrefix; `APP` by default)                           |
| [`asp.long`](#asplong)           | long `--some-name` style CLI flag                                                           |
| [`asp.short`](#aspshort)         | short `-n` style CLI flag                                                                   |
| [`asp.sensitive`](#aspsensitive) | indicates that the value is "sensitive" and should be redacted from SerializeFlags output.  |

If you are consistently providing most or all of the values, the `asp` tag is a bit more concise.

### `asp`

The “all the tags” tag, `asp:""` allows you to specify the long, short, env, desc, and sensitive values, separated by commas. The "explicit" tags always take precedence, but any non-empty portions of `asp` take precedence over the default fallback values. To _omit_ a value, the explicit attribute tag must be used.

### `asp.desc`

Sets the usage description for the flag. This is a [Go-style template string](https://pkg.go.dev/text/template) with several values and functions available.

> Examples assume the config with default options:
>
> ```go
> type config struct {
>     Outer struct {
>         Inner string `asp.short:"x"` // <------- this field
>     }
> }
> ```
>
> or called with `.Name` for the functions.

| name               | kind     | description                                                                                                      | example           |
| ------------------ | -------- | ---------------------------------------------------------------------------------------------------------------- | ----------------- |
| .Name              | value    | complete field name, `.`-delimited                                                                               | `Outer.Inner`     |
| .Long              | value    | long-form flag name, `-`-delimited                                                                               | `outer-inner`     |
| .Short             | value    | short-form flag name                                                                                             | `x`               |
| .Env               | value    | environment variable name, `_`-delimited                                                                         | `APP_OUTER_INNER` |
| .NoEnv             | value    | empty string sentinel to prevent automatic appending of environment variable                                     |                   |
| .ParentName        | value    | complete _parent_ field name, `.`-delimited                                                                      | `Outer`           |
| camel              | function | [`strcase.ToCamel`](https://pkg.go.dev/github.com/iancoleman/strcase#ToCamel) function                           | `OuterInner`      |
| delimited          | function | [`strcase.ToDelimited`](https://pkg.go.dev/github.com/iancoleman/strcase#ToDelimited) function                   | `outer.inner`     |
| kebab              | function | [`strcase.ToKebab`](https://pkg.go.dev/github.com/iancoleman/strcase#ToKebab) function                           | `outer-inner`     |
| lowerCamel         | function | [`strcase.ToLowerCamel`](https://pkg.go.dev/github.com/iancoleman/strcase#ToLowerCamel) function                 | `outerInner`      |
| screamingDelimited | function | [`strcase.ToScreamingDelimited`](https://pkg.go.dev/github.com/iancoleman/strcase#ToScreamingDelimited) function | `OUTER.INNER`     |
| screamingKebab     | function | [`strcase.ToScreamingKebab`](https://pkg.go.dev/github.com/iancoleman/strcase#ToScreamingKebab) function         | `OUTER-INNER`     |
| screamingSnake     | function | [`strcase.ToScreamingSnake`](https://pkg.go.dev/github.com/iancoleman/strcase#ToScreamingSnake) function         | `OUTER_INNER`     |
| snake              | function | [`strcase.ToSnake`](https://pkg.go.dev/github.com/iancoleman/strcase#ToSnake) function                           | `outer_inner`     |
| snakeWithIgnore    | function | [`strcase.ToSnakeWithIgnore`](https://pkg.go.dev/github.com/iancoleman/strcase#ToSnakeWithIgnore) function       | `outer_inner`     |

> [!NOTE]
>
> As of v0.2.3, the [sprig](https://masterminds.github.io/sprig/) functions are also available.

If no `asp.desc:"DESCRIPTION"` or `asp:",,,DESCRIPTION"` tag is present, asp defaults to:

```go
"sets the {{delimited .Name ' '}} value (env: {{.Env}})"
```

which results in:

```
"sets the outer inner value (env: APP_OUTER_INNER)"
```

If an description override is present, asp looks to see if either `{{.Env}}` or `{{.NoEnv}}` is included in the string. If not, it automatically appends `" (env: {{.Env}})"` You can use `{{.Env}}` to include the environment variable name in a specific location in the description, or use `{{.NoEnv}}` to indicate that the automatic appending should be avoided.

### `asp.env`

Provides an override value for “this field’s” portion of the an environment variable name. In the case of a value field, the terminal term in the name; for a nested struct, a middle part of the name. Explicitly setting an empty string (`asp.env:""`) will omit that segment in the name.

### `asp.long`

Provides an override value for “this field’s” portion of the a long flag name. In the case of a value field, the terminal term in the name; for a nested struct, a middle part of the name. Explicitly setting an empty string (`asp.long:""`) will omit that segment in the name.

### `asp.short`

Provides the short flag (single character) for the field.

### `asp.sensitive`

If set to `true` (`asp.sensitive:"true"`), the SerializeFlags function will use `[REDACTED]` in place of the actual value.
