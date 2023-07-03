package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

func NewRootCmd() *cobra.Command {
	newRootCmd := &cobra.Command{
		Use:   "manager",
		Short: "A CLI to manage APId deployments",
		Long: `This CLI will handle deployment addition, removals, and administration
	for the University Of Galway Computer Society's APId.`,
	}

	newRootCmd.PersistentFlags().Bool("debug", false, "Show debug messages")
	newRootCmd.AddCommand(NewConfigCmd())
	return newRootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd = NewRootCmd()
}
