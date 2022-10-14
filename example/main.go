package main

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

type Config struct {
	SomeValue   string
	SomeFlag    bool
	ManyNumbers []int

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

	a, err := asp.Attach(
		cmd, defaults,
		asp.WithDefaultConfigName[Config]("asp-example"),
	)
	cobra.CheckErr(err)

	// Ensure the `asp.Asp` value is available to the command handler when it
	// runs.  You can also store the returned value in a global, but using
	// context helps when you have more than one command with differing config
	// structures.
	cmd.ExecuteContext(
		context.WithValue(context.Background(), asp.ContextKey, a))
}

func commandHandler(cmd *cobra.Command, args []string) {
	// Extract the `asp.Asp` from the context and get the parsed config.
	a := cmd.Context().Value(asp.ContextKey).(asp.Asp[Config])
	config := a.Config()

	log.Printf("got config: %#v", config)
}
