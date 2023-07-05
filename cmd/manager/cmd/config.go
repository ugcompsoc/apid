package cmd

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ugcompsoc/apid/cmd/manager/utils"
)

// configCmd represents the config command
var configCmd *cobra.Command

func NewConfigCmd() *cobra.Command {
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Displays a config given a directory",
		Long: `Prints out the config in the default directory if no directory is specified.
If a directory is specified it will look for a compatible file in there and
will print it out instead.`,
		Run: VerifyConfig,
	}

	configCmd.AddCommand(NewCreateConfigCmd())
	configCmd.PersistentFlags().BoolP("print", "p", false, "Print config")
	configCmd.PersistentFlags().BoolP("secrets", "s", false, "Print secrets")

	return configCmd
}

func VerifyConfig(cmd *cobra.Command, args []string) {
	filename, _ := cmd.Flags().GetString("filename")
	err := utils.VerifyFilename(filename)
	if err != nil {
		cmd.Printf("An error occured while verifying the filename: %s\n", err)
		return
	}

	directory, _ := cmd.Flags().GetString("directory")
	absoluteFilePath := filepath.Join(directory, filename)
	c, err := utils.ExtractFile(absoluteFilePath)
	if err != nil {
		cmd.Printf("An error occured while extracting the file: %s\n", err)
		return
	}

	issues, err := c.Verify()
	if err != nil {
		cmd.Printf("An error occured while verifying the config: %s\n", err)
		return
	}
	if len(issues) != 0 {
		cmd.Printf("Error(s) were found while parsing %s, please address them:\n", absoluteFilePath)
		for _, err := range issues {
			cmd.Printf("  - %s\n", err)
		}
		return
	}

	cmd.Print("OK\n")

	print, _ := cmd.Flags().GetBool("print")
	printSecrets, _ := cmd.Flags().GetBool("secrets")
	if print {
		yamlStr, err := utils.PrintConfig(c, printSecrets)
		if err != nil {
			cmd.Printf("An error occured while attempting to print the config: %s\n", err)
			return
		}
		cmd.Printf("\n%s", yamlStr)
	}
}

func init() {
	configCmd = NewConfigCmd()
}
