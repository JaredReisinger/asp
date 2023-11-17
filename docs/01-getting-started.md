# Getting started

In the most basic case, getting started with `asp` is very easy if you already have a simple CLI app that’s using [`cobra`](https://cobra.dev/) to define it’s command structure, especially if you don’t have too many flags already implemented.

First, get the latest version of asp for use in your project:

```sh
go get github.com/jaredreisinger/asp@latest
```

Let’s assume you’re using Cobra and Go best practices, and have a `cmd/root.go` file with the root `cobra.Command` for your app. Further, let’s assume the final version of the example from [Cobra’s _Create rootCmd_ instructions](https://cobra.dev/#create-rootcmd):

```go
package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
love by spf13 and friends in Go.
Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
```

We'll start by defining a `struct` that holds all of the configuration values we see as flags (except for the config file; as it turn out, asp will provide this automatically):

```go
type rootConfig struct {
	ProjectBase string
	Author      string
	License     string
	UseViper    bool
}
```

Next, we’ll need to know how to tell asp to attach that configuration definition to `rootCmd`, and how to get the user-specified values when the command runs. Those are:

```go
asp.Attach(rootCmd, rootConfig{})
```

and

```go
asp.Get[rootConfig](cmd)
```

Aside from a bit of error handling and the import for asp, that’s really all you need!

```diff
package cmd

import (
	"fmt"
	"os"

+	"github.com/jaredreisinger/asp"
-	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
-	"github.com/spf13/viper"
)

+type rootConfig struct {
+	ProjectBase string
+	Author      string
+	License     string
+	UseViper    bool
+}

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
love by spf13 and friends in Go.
Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
+		cfg, err := asp.Get[rootConfig](cmd)
+		if err != nil {
+			fmt.Println(err)
+			os.Exit(1)
+		}

		// Do Stuff Here
+		// ... and use `cfg.ProjectBase`, `cfg.Author`, etc.
	},
}

func Execute() {
+	if err := asp.Attach(rootCmd, rootConfig{}); err != nil {
+		fmt.Println(err)
+		os.Exit(1)
+	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

-func init() {
-	cobra.OnInitialize(initConfig)
-	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
-	rootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
-	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
-	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
-	rootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")
-	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
-	viper.BindPFlag("projectbase", rootCmd.PersistentFlags().Lookup("projectbase"))
-	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
-	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
-	viper.SetDefault("license", "apache")
-}
-
-func initConfig() {
-	// Don't forget to read config either from cfgFile or from home directory!
-	if cfgFile != "" {
-		// Use config file from the flag.
-		viper.SetConfigFile(cfgFile)
-	} else {
-		// Find home directory.
-		home, err := homedir.Dir()
-		if err != nil {
-			fmt.Println(err)
-			os.Exit(1)
-		}
-
-		// Search config in home directory with name ".cobra" (without extension).
-		viper.AddConfigPath(home)
-		viper.SetConfigName(".cobra")
-	}
-
-	if err := viper.ReadInConfig(); err != nil {
-		fmt.Println("Can't read config:", err)
-		os.Exit(1)
-	}
-}
```

We’ve removed 39 lines, and added back only 17 (it would have been only 12 if we were using `cobra.CheckErr(err)`). If you look at the final code, it’s a lot less… everything, and easier to read:

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/jaredreisinger/asp"
	"github.com/spf13/cobra"
)

type rootConfig struct {
	ProjectBase string
	Author      string
	License     string
	UseViper    bool
}

var rootCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
love by spf13 and friends in Go.
Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := asp.Get[rootConfig](cmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_ = cfg

		// Do Stuff Here
		// ... and use `cfg.ProjectBase`, `cfg.Author`, etc.
	},
}

func Execute() {
	if err := asp.Attach(rootCmd, rootConfig{}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
```

Running the changed code (`go run . --help`) shows an almost-equivalent result:

```
A Fast and Flexible Static Site Generator built with
love by spf13 and friends in Go.
Complete documentation is available at http://hugo.spf13.com

Usage:
  hugo [flags]

Flags:
      --author string         sets the author value (env: APP_AUTHOR)
      --config string         configuration file to load
  -h, --help                  help for hugo
      --license string        sets the license value (env: APP_LICENSE)
      --project-base string   sets the project base value (env: APP_PROJECTBASE)
      --use-viper             sets the viper value (env: APP_USEVIPER)
```

You can see that all of the expected command-line flags are present, including `--config`. There are a few differences, though:

1. The original `--projectbase` flag is now `--project-base`, and `--viper` is `--use-viper`.

2. The flag descriptions aren’t quite what they were before, but they’re serviceable.

3. The help listing is now mentioning equivalent environment variables. You can use `APP_AUTHOR=me hugo ...` instead of `hugo --author me ...` if you want to.

4. The config file structure is driven directly from the type of the configuration structure. Instead of interpreting dotted names in the viper binding lines, you can just look at the declared type.

You can find out more about how to solve the first and second issues in XXX.
