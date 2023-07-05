package cmd

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/ugcompsoc/apid/cmd/manager/utils"
	"github.com/ugcompsoc/apid/internal/config"
	"gopkg.in/yaml.v2"
)

// createConfigCmd represents the config command
var createConfigCmd *cobra.Command

func NewCreateConfigCmd() *cobra.Command {
	createConfigCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a config for the APId",
		Long: `Given a set of APId configuration flags and variables, a config file will
be generated in the folder specified.`,
		Run: CreateConfig,
	}

	createConfigCmd.Flags().String("log_level", "debug", "Log level; Available values: trace, disabled, panic, fatal, error, warn, info, or debug")
	createConfigCmd.Flags().String("timeouts_startup", "30s", "Startup Timeout")
	createConfigCmd.Flags().String("timeouts_shutdown", "30s", "Shutdown Timeout")
	createConfigCmd.Flags().String("http_listen_address", ":8080", "HTTP Listen Address; In the form of 'IP/DOMAIN:PORT'")
	createConfigCmd.Flags().StringSlice("http_cors_allowed_orgins", []string{"*"}, "HTTP CORS Allowed Origins; In the form '[ORIGIN,ORIGIN]'")
	createConfigCmd.Flags().String("database_host", "mongodb://ugcompsoc_apid_local_db", "Database Host")
	createConfigCmd.Flags().String("database_name", "apid", "Database Name")
	createConfigCmd.Flags().String("database_username", "", "Database Username")
	createConfigCmd.Flags().String("database_password", "", "Database Password")
	createConfigCmd.MarkFlagRequired("database_username")
	createConfigCmd.MarkFlagRequired("database_password")

	return createConfigCmd
}

func CreateConfig(cmd *cobra.Command, args []string) {
	c := &config.Config{}
	issues := []string{}
	var err error

	c.LogLevel, _ = createConfigCmd.Flags().GetString("log_level")
	startupTimeout, _ := createConfigCmd.Flags().GetString("timeouts_startup")
	c.Timeouts.Startup, err = time.ParseDuration(startupTimeout)
	if err != nil {
		issues = append(issues, "Could not parse startup timeout. Use the format '[NUMBER]s'")
	}
	shutdownTimeout, _ := createConfigCmd.Flags().GetString("timeouts_shutdown")
	c.Timeouts.Shutdown, err = time.ParseDuration(shutdownTimeout)
	if err != nil {
		issues = append(issues, "Could not parse shutdown timeout. Use the format '[NUMBER]s'")
	}
	c.HTTP.ListenAddress, _ = createConfigCmd.Flags().GetString("http_listen_address")
	c.HTTP.CORS.AllowedOrigins, _ = createConfigCmd.Flags().GetStringSlice("http_cors_allowed_orgins")
	c.Database.Host, _ = createConfigCmd.Flags().GetString("database_host")
	c.Database.Name, _ = createConfigCmd.Flags().GetString("database_name")
	c.Database.Username, _ = createConfigCmd.Flags().GetString("database_username")
	c.Database.Password, _ = createConfigCmd.Flags().GetString("database_password")
	if c.Database.Username == "" {
		issues = append(issues, "Database username has no default value and is required")
	}
	if c.Database.Password == "" {
		issues = append(issues, "Database password has no default value and is required")
	}

	if len(issues) != 0 {
		cmd.Print("Error(s) were found while generating the config, please address them:\n")
		for _, issue := range issues {
			cmd.Printf("  - %s\n", issue)
		}
		return
	}

	filename, _ := cmd.Flags().GetString("filename")
	err = utils.VerifyFilename(filename)
	if err != nil {
		cmd.Printf("An error occured while verifying the filename: %s\n", err)
		return
	}
	directory, _ := cmd.Flags().GetString("directory")
	absoluteFilePath := filepath.Join(directory, filename)

	issues, err = c.Verify()
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

	cYaml, err := yaml.Marshal(c)
	if err != nil {
		cmd.Printf("Could not marshall the config struct: %s\n", err)
		return
	}
	err = os.WriteFile(absoluteFilePath, cYaml, 0644)
	if err != nil {
		cmd.Printf("Could not write file to %s: %s\n", absoluteFilePath, err)
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
	createConfigCmd = NewCreateConfigCmd()
}
