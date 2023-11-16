package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

type Config struct {
	SomeValue       string
	SomeFlag        bool
	ManyNumbers     []int
	MapStringInt    map[string]int
	MapStringString map[string]string

	SubSection struct {
		NamesLikeThis string
	}

	Verbose bool `asp.short:"v" asp.desc:"get noisy"`
}

func main() {
	defaults := Config{
		SomeValue: "DEFAULT STRING!",
	}

	cmd := &cobra.Command{
		Run: commandHandler,
	}

	err := asp.Attach(
		cmd, defaults,
		asp.WithDefaultConfigName("asp-example"),
	)
	cobra.CheckErr(err)

	err = cmd.Execute()
	cobra.CheckErr(err)
}

func commandHandler(cmd *cobra.Command, args []string) {
	// get the config using the asp.Asp instance attached to cmd
	config, err := asp.Get[Config](cmd)
	cobra.CheckErr(err)

	log.Printf("got config: %#v", config)
}
