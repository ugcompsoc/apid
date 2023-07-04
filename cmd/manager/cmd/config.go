package cmd

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/ugcompsoc/apid/internal/config"
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
	configCmd.PersistentFlags().StringP("directory", "d", ".", "Directory for config")
	configCmd.PersistentFlags().StringP("filename", "f", "apid.yml", "filename for config")
	configCmd.PersistentFlags().BoolP("print", "p", false, "Print config")
	configCmd.PersistentFlags().BoolP("secrets", "s", false, "Print secrets")

	return configCmd
}

func verifyFileAndExtract(file []byte) (*config.Config, []string, error) {
	fileStr := string(file)
	if len(fileStr) == 0 {
		return nil, nil, errors.New("The file apid.yml is completely empty. What do you want me to do with this?")
	}
	c := config.Config{}
	err := yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, nil, err
	}

	issues, err := c.Verify()
	return &c, issues, err
}

func VerifyConfig(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	directory, _ := cmd.Flags().GetString("directory")
	filename, _ := cmd.Flags().GetString("filename")
	filenameRegex, err := regexp.Compile("^([a-z]+).yml$")
	// unless this regex is changes, this error will never be reached
	if err != nil {
		if debug {
			cmd.Printf("Error: %s\n", err)
		}
		cmd.Printf("An error occured while generating the filename regex\n")
		return
	}
	if !filenameRegex.MatchString(filename) {
		cmd.Print("The filename is not in the form [NAME].yml")
		return
	}

	absoluteFilePath := filepath.Join(directory, filename)
	file, err := ioutil.ReadFile(absoluteFilePath)
	if err != nil {
		if debug {
			cmd.Printf("Error: %s\n", err)
		}
		cmd.Printf("No file exists at path: %s\n", absoluteFilePath)
		return
	}

	c, errs, err := verifyFileAndExtract(file)
	if err != nil {
		if debug {
			cmd.Printf("Error: %s\n", err)
		}
		cmd.Printf("An error was encountered while verifing the file\n")
		return
	}
	if len(errs) != 0 {
		cmd.Printf("Error(s) were found while parsing %s, view them below and address them\n", filename)
		for _, err := range errs {
			cmd.Printf("  - %s\n", err)
		}
		return
	}

	cmd.Print("OK\n")

	print, _ := cmd.Flags().GetBool("print")
	if print {
		printSecrets, _ := cmd.Flags().GetBool("secrets")
		if !printSecrets {
			if len(c.Database.Username) != 0 {
				c.Database.Username = "********"
			}
			if len(c.Database.Password) != 0 {
				c.Database.Password = "********"
			}
		}

		// Should really never catch an error here unless the config values
		// are changed between here and when they are verified before or
		// or for whatever reason we somehow gave a value that isn't valid
		// in YAML
		cYaml, err := yaml.Marshal(&c)
		if err != nil {
			if debug {
				cmd.Printf("Error: %s\n", err)
			}
			cmd.Print("Could not marshall the config struct\n")
			return
		}
		cmd.Printf("\n%s", string(cYaml))
	}
}

func init() {
	configCmd = NewConfigCmd()
}
