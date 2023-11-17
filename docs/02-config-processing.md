# Config processing

Let’s look again at the configuration type declaration from [_Getting started_](01-getting-started.md), and what we get as a result:

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

You can see that the flag names and environment variable names come from the configuration field names. Both `Author` and `License` are simple names, lower-cased for the flag, and upper-cased for the environment variable. `ProjectBase` and `UseViper` on the other hand, show that asp is recognizing that there are two words in the name, and using standard CLI flag behavior by hyphenating them (`project-base` and `use-viper`). Meanwhile, the environment variables are left as single uninterrupted uppercase terms (`PROJECTBASE` and `USEVIPER`).

## Types

Asp supports many of the flag types supported by `cobra` (by `pflags`, really):

| type                | format (CLI argument and/or environment)                                                               |
| ------------------- | ------------------------------------------------------------------------------------------------------ |
| `bool`              | no argument needed for `true`, environment variables can use `true`, `false`, `1`, `0`, `yes`, or `no` |
| `int`               | number as a string                                                                                     |
| `uint`              | number as a string                                                                                     |
| `string`            | the string value                                                                                       |
| `[]int`             | comma-separated numbers                                                                                |
| `[]uint`            | comma-separated numbers                                                                                |
| `[]string`          | comma-separated strings; note that an individual string value cannot itself contain a comma!           |
| `map[string]int`    | comma-separated, equal-delimited string/number pairs (like `"a=1,b=4"`)                                |
| `map[string]string` | comma-separated, equal-delimited string/string pairs (like `"a=foo,b=bar"`)                            |
| `time.Time`         | RFC3399Nano format, or the literals `now`, `local`, or `utc`                                           |
| `time.Duration`     | ~~ISO8601 duration format~~ [Go duration format](https://pkg.go.dev/time#ParseDuration)                |

To extend this list, or change the parsing behavior, see [WithDecodeHook](04-options.md#withdecodehook), but be aware that asp cannot currently map non-default types to flags.

## Nested fields

There are times you want to provide more semantic structure to your configuration, so let's see what happens if we make the `Author` a structure with both `Name` and `Email` properties:

```go
type rootConfig struct {
    Author      struct {
        Name  string
        Email string
    }
}
```

This results in:

```
Flags:
      --author-email string   sets the author email value (env: APP_AUTHOR_EMAIL)
      --author-name string    sets the author name value (env: APP_AUTHOR_NAME)
```

Nested `struct`s result in flags and environment variables with the parent field name as a prefix, and with an underscore separator. A nested type can be used any number of times, since the field names in the parent will be unique:

```go
type person struct {
    Name  string
    Email string
}

type rootConfig struct {
    Author     person
    Editor     person
    Supervisor person
}
```

```
Flags:
      --author-email string       sets the author email value (env: APP_AUTHOR_EMAIL)
      --author-name string        sets the author name value (env: APP_AUTHOR_NAME)
      --editor-email string       sets the editor email value (env: APP_EDITOR_EMAIL)
      --editor-name string        sets the editor name value (env: APP_EDITOR_NAME)
      --supervisor-email string   sets the supervisor email value (env: APP_SUPERVISOR_EMAIL)
      --supervisor-name string    sets the supervisor name value (env: APP_SUPERVISOR_NAME)
```

There is one caveat, however. Since asp also hyphenates multi-word field names, you can create a valid Go structure that results in flag name collision:

```go
type rootConfig struct {
    AuthorEmail string
    Author      struct {
        Email string
    }
}
```

It compiles successfully, but running the code will panic when the second `--author-email` flag is added. This is a contrived example, and should very rarely ever happen in actual code; you might not have collision in field references, but it would be confusing to have both `cfg.AuthorEmail` and `cfg.Author.Email` in the code.

## Anonymous structs

In some cases, you’ll have a struct whose fields you want to use _directly_ in the enclosing configuration, as an anonymous struct. Asp supports this as well:

```go
type Person struct {
    Name  string
    Email string
}

type rootConfig struct {
    Person `mapstructure:",squash"`
}
```

```
Flags:
      --email string    sets the email value (env: APP_EMAIL)
      --name string     sets the name value (env: APP_NAME)
```

Do note, however, that as written above, asp (viper) will expect a config file like:

```yaml
person:
  name: some name
  email: someone@example.org
```

If you want the config file flattened in the same way that the flags and environment variables are—which will be more consistent and “obvious” to your users—you need to add an annotation to the field:

```go
type rootConfig struct {
    Person `mapstructure:",squash"`
}
```

This tells the unmarshaling code to [expect the anonymous structs field inlined in the parent](https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Embedded_Structs_and_Squashing):

```yaml
# not under a "person" key!
name: some name
email: someone@example.org
```
