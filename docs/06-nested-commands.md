# Nested commands

Asp was originally designed with only simple, single-command apps in mind, or apps where _all_ of the flags were defined on the child command.  As of 0.4.1, Asp better supports complex nested-command designs.

To use `asp.Attach` at multiple levels of a `cobra.Command` tree, make sure to set [`cobra.EnableTraverseRunHooks`](https://pkg.go.dev/github.com/spf13/cobra#EnableTraverseRunHooks) to `true`. Without this, only the `asp.Attach` for the invoked command will result in a `Get`-able config; the parent command configurations will not be stashed in those commands' context.

```go
type rootConfig struct {
    Flag bool
}

type childConfig struct {
    AnotherFlag bool
}

var (
    rootCmd = &cobra.Command{ ... }
    childCmd = &cobra.Command{
        Run: runChild,
    }
)

func init() {
  cobra.EnableTraverseRunHooks = true
  asp.Attach(rootCmd, rootConfig{})
  asp.Attach(childCmd, childConfig{})
  rootCmd.AddCommand(childCmd)
}

func runChild(cmd *cobra.Command, args []string) {
    // Note that the root config must be fetched *from the root command*
    rootCfg, err := asp.Get[rootConfig](cmd.Root())
    if err != nil {
      return err
    }

    childCfg, err := asp.Get[childConfig](cmd)
    if err != nil {
      return err
    }
}

func main() {
  rootCmd.Execute()
}
```

## Future thoughts

Requiring the child to make multiple `asp.Get()` calls is a bit awkward. Nicer would be a mechanism to include the parents' configs as a member in the child config, in a way that can automatically climb the command tree and extract the values. That constitutes a new feature, however, and is a bigger change than the "fix" for 0.4.1 warrants.
