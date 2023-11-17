# Roadmap / backlog

As I use `asp` in other tools, I sometimes stumble across missing features, or things I think could be improved. These are those ideas, in no particular order.

- **Allow sub-structs to define a different CLI-flag or environment-variable “part” name.**

  > I’ve run across cases where the name in the configuration isn’t quite the name I want to use in flags or the environment. For example, using `LaTeX` in code reflects the proper name of the tool, but results in flags with `la-te-x` in the name because we infer word-parts from the capitalization. We should be able to use `asp.long` on the struct field to define the “effective” name for the field.

  _Sub-struct names added in [v0.2.0](https://github.com/JaredReisinger/asp/releases/tag/v0.2.0)!_

- **Support ISO durations**

  > Go’s native `time.\Duration` only parses hours or smaller units, because days suffer from daylight saving time issues (making some days 23 hours, and some 25), months have different lengths, and years have leap-days. But in practice, some applications want an easy way to configure “2 months”. Using a `string` and parsing as-needed works, but durations are common enough, and the ISO duration format is standard enough that it would be nice to support it directly. (If there’s an existing ISO duration library, maybe I could use that?)

- **Support extensible type deserializtion?**

  > Maybe I do? Need to re-read the code and document it!

  _Added in [v0.2.0](https://github.com/JaredReisinger/asp/commit/758d7077bbc998905cf4361bf44d71d2fd799a35)!_

# Thoughts on `asp` _(a brief history)_

I started with `cobra` to provide CLI support for some simple tools, and then naturally looked at `viper` when I started to think about 12-factor-izing things. To me, the amount of `viper`-specific configuration code felt like it was going to dwarf the tool I was writing.

What I wanted was an _easy_ way to make a flag/config/environment-variable _something_ for each of the configuration settings I had. At the time I starting thinking about `asp`, I didn’t see anything that matched what I wanted.

## What problem is `asp` attempting to solve?

As I see it, writing a well-behaved 12-factor app _ought_ to be as easy as listing out all of your configuration values and types—which is basically what a Go `struct` is! This is `asp`’s motto in a nutshell: _“You define your configuration struct, and `asp` will take care of the rest.”_

## Why not `viper`? Why not `envconfig`?

For the vast majority of cases (e.g. the 80:20 rule), a CLI application needs to be able to read/load settings from CLI flags, from a configuration file, and from environment variables. There are libraries to do each of these individually, and libraries to do some subset of these things, but I wasn’t seeing a solution that tried to solve _all_ of them holistically.

While `viper` claims to be a “complete configuration solution”, its CLI flag support works “with” the flags you write for `cobra` (or some other flag system), rather than providing its own. It also requires you to write a lot of boilerplate code to define the various configuration settings, the mapping between those settings and CLI flags, the mapping between those settings and environment variables, and so on.

That said, `asp` _absolutely relies_ on the functionality that `viper` (and `cobra`) provides; it just tries to handle all of the boilerplate code on your behalf, so that the only thing you need to worry about is defining your configuration structure.

### Feature comparison

| Feature                    | `asp`  | `cobra` | `viper` | `envconfig` |
| -------------------------- | :----: | :-----: | :-----: | :---------: |
| CLI flags                  | “free” |   yes   |  yes\*  |     no      |
| environment variables      | “free” |   no    |   yes   |     yes     |
| configuration structure    |  yes   |   no    |   yes   |     no      |
| configuration file (read)  | “free” |   no    |   yes   |     no      |
| configuration file (write) |   no   |   no    |   yes   |     no      |
| remote configuration file  |   no   |   no    |   yes   |     no      |

(\* As mentioned above, `viper` “works with” another package’s flags.)

## Why is `asp` “config structure first”?

Ultimately, your application logic just wants to be able to read its needed settings, and a configuration structure (a) provides build-time static typing and name safety, and (b) allows for easy aggregation of parts (sub-structs and anonymous structs). It is a much “higher-level” description of the needed data than either CLI flags or environment variables, and has a nearly 1:1 correlation with a config file.

## Significant changes

### 0.1 to 0.2

#### Nested struct names

There’s an interesting challenge with deeply-nested configs; the outer struct _**should be in control**_ of the final naming of the flags and environment variables. (It’s not necessary for the config, because it’s already providing a structural namespace.) This allows embedded structs to provide leaf names _without_ getting collisions. Version 0.1’s `asp.long` attribute assumed it provided an absolute, unscoped flag name:

```go
type Config struct {
  Inner struct {
    Value string `asp.long:"something"`
  }
}
```

would create the flag `--something`, not `--inner-something`. This design causes:

```go
type Config struct {
  Inner1 Inner
  Inner2 Inner
}

type Inner struct {
    Value string `asp.long:"something"`
}
```

to break, attempting to use `--something` twice. There’s an analogous problem with environment variables for the same reason.

Using the guideline of “outer struct is in control”, a first attempt to solve this is to use the outer field name as a part of the flag/envvar _always_. This means the inner `asp.long` is no longer absolute. This is a **BREAKING CHANGE**. The mitigation to get the old behavior is for the outer field to explicitly remove itself from the name, using `asp.long:""`. But, this also means that you can _**no longer**_ create a single field with a flag or envvar name of different “depth” than its siblings:

```go
type Config struct {
  Inner struct {
    Value string `asp.long:"something"`
    Sibling string
  }

  Unnamed struct {
    Value string `asp.long:"another"`
    Sibling string
  } `asp.long:""`
}
```

| field           | 0.1                 | 0.2                 |
| --------------- | ------------------- | ------------------- |
| Inner.Value     | `--something`       | `--inner-something` |
| Inner.Sibling   | `--inner-sibling`   | `--inner-sibling`   |
| Unnamed.Value   | `--another`         | `--another`         |
| Unnamed.Sibling | `--unnamed-sibling` | `--sibling`         |

_(Note that sub-struct fields did **not** have their tags evaluated in 0.1.)_
