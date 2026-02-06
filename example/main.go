package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/jaredreisinger/asp"
)

type rootConfig struct {
	SomeValue       string
	SomeFlag        bool
	ManyNumbers     []int
	MapStringInt    map[string]int
	MapStringString map[string]string

	SubSection struct {
		NamesLikeThis string
	}

	SecretValue string `asp.sensitive:"true"`

	Verbose bool `asp.short:"v" asp.desc:"get noisy"`
}

type childConfig struct {
	AnotherValue string
	AnotherFlag  bool
}

var (
	rootCmd = &cobra.Command{
		Use: "root",
		Run: runRoot,
	}

	childCmd = &cobra.Command{
		Use: "child",
		Run: runChild,
	}
)

func init() {
	err := asp.Attach(rootCmd, rootConfig{
		SomeValue: "DEFAULT STRING!",
	}, asp.WithDefaultConfigName("asp-example"))
	cobra.CheckErr(err)

	err = asp.Attach(childCmd, childConfig{})
	cobra.CheckErr(err)

	rootCmd.AddCommand(childCmd)
	cobra.EnableTraverseRunHooks = true
}

func main() {
	err := rootCmd.Execute()
	cobra.CheckErr(err)
}

func runRoot(cmd *cobra.Command, args []string) {
	// get the config using the asp.Asp instance attached to cmd
	config, err := asp.Get[rootConfig](cmd)
	cobra.CheckErr(err)

	log.Printf("got config: %#v", config)

	s, err := asp.SerializeFlags(config, true)
	cobra.CheckErr(err)
	log.Printf("serialized flags: %s", s)
}

func runChild(cmd *cobra.Command, args []string) {
	// get the rootCfg using the asp.Asp instance attached to cmd
	rootCfg, err := asp.Get[rootConfig](cmd.Root())
	cobra.CheckErr(err)
	log.Printf("got root config: %#v", rootCfg)

	childCfg, err := asp.Get[childConfig](cmd)
	cobra.CheckErr(err)
	log.Printf("got child config: %#v", childCfg)

	s, err := asp.SerializeFlags(rootCfg, true)
	cobra.CheckErr(err)
	log.Printf("serialized root flags: %s", s)

	s, err = asp.SerializeFlags(childCfg, true)
	cobra.CheckErr(err)
	log.Printf("serialized child flags: %s", s)
}
