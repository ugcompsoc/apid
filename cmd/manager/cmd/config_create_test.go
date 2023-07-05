package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfig(t *testing.T) {
	helpMessage := `Usage:
  manager config create [flags]

Flags:
      --database_host string               Database Host (default "mongodb://ugcompsoc_apid_local_db")
      --database_name string               Database Name (default "apid")
      --database_password string           Database Password
      --database_username string           Database Username
  -h, --help                               help for create
      --http_cors_allowed_orgins strings   HTTP CORS Allowed Origins; In the form '[ORIGIN,ORIGIN]' (default [*])
      --http_listen_address string         HTTP Listen Address; In the form of 'IP/DOMAIN:PORT' (default ":8080")
      --log_level string                   Log level; Available values: trace, disabled, panic, fatal, error, warn, info, or debug (default "debug")
      --timeouts_shutdown string           Shutdown Timeout (default "30s")
      --timeouts_startup string            Startup Timeout (default "30s")

Global Flags:
  -d, --directory string   Directory for config (default ".")
  -f, --filename string    filename for config (default "apid.yml")
  -p, --print              Print config
  -s, --secrets            Print secrets`

	directory := "config_create_test"
	if _, err := os.Stat(directory); err != nil {
		err := os.Mkdir(directory, os.ModePerm)
		assert.NoError(t, err, "could not create test directory")
	}
	absoluteFilePath := filepath.Join(directory, "apid.yml")

	argsWithDirectory := []string{"config", "create", "-d=" + directory}
	argsWithDatabaseSet := append(argsWithDirectory, []string{"--database_password=test_password", "--database_username=test_username"}...)

	errorsWereFoundGenerating := "Error(s) were found while generating the config, please address them:\n"
	errorsWereFoundVerifying := "Error(s) were found while parsing " + absoluteFilePath + ", please address them:\n"

	runs := []struct {
		name string
		args []string
		out  string
		err  error
	}{
		{
			name: "no default set for database username and password",
			args: []string{"config", "create"},
			out:  "Error: required flag(s) \"database_password\", \"database_username\" not set\n" + helpMessage,
			err:  errors.New("Error: required flag(s) \"database_password\", \"database_username\" not set"),
		},
		{
			name: "no default set for database username",
			args: []string{"config", "create", "--database_password='test'"},
			out:  "Error: required flag(s) \"database_username\" not set\n" + helpMessage,
			err:  errors.New("Error: required flag(s) \"database_username\" not set"),
		},
		{
			name: "no default set for database password",
			args: []string{"config", "create", "--database_username='testpassword'"},
			out:  "Error: required flag(s) \"database_password\" not set\n" + helpMessage,
			err:  errors.New("Error: required flag(s) \"database_password\" not set"),
		},
		{
			name: "no errors and no issues - happy path",
			args: argsWithDatabaseSet,
			out:  "OK",
		},
		{
			name: "will not parse startup timeout",
			args: append(argsWithDatabaseSet, "--timeouts_startup", "30"),
			out:  errorsWereFoundGenerating + "  - Could not parse startup timeout. Use the format '[NUMBER]s'",
		},
		{
			name: "will not parse shutdown timeout",
			args: append(argsWithDatabaseSet, "--timeouts_shutdown", "30"),
			out:  errorsWereFoundGenerating + "  - Could not parse shutdown timeout. Use the format '[NUMBER]s'",
		},
		{
			name: "database username has no default value",
			args: append(argsWithDatabaseSet, "--database_username", ""),
			out:  errorsWereFoundGenerating + "  - Database username has no default value and is required",
		},
		{
			name: "database password has no default value",
			args: append(argsWithDatabaseSet, "--database_password", ""),
			out:  errorsWereFoundGenerating + "  - Database password has no default value and is required",
		},
		{
			name: "log level is not set, from verify function",
			args: append(argsWithDatabaseSet, "--log_level", ""),
			out:  errorsWereFoundVerifying + "  - An invalid log level was specified",
		},
		{
			name: "will fail to write file because there is no directory",
			args: []string{"config", "create", "--directory=dir_not_in_existance", "--database_password=test_password", "--database_username=test_username"},
			out:  "Could not write file to dir_not_in_existance/apid.yml: open dir_not_in_existance/apid.yml: no such file or directory",
		},
		{
			name: "an invalid filename will cause an issue",
			args: append(argsWithDatabaseSet, "--filename=apid"),
			out:  "An error occured while verifying the filename: The filename is not in the form [NAME].yml",
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			out, err := execute(t, NewRootCmd(), run.args...)
			if run.err == nil {
				assert.NoError(t, err, "expected no error running manager")
			} else {
				assert.Error(t, err, "expected error running manager")
			}
			assert.Equal(t, run.out, out, "unexpected manager output")
			if _, err := os.Stat(directory + "/apid.yml"); err == nil {
				os.Remove(directory + "/apid.yml")
			}
		})
	}

	t.Run("prints multiple issues", func(t *testing.T) {
		out, err := execute(t, NewRootCmd(), append(argsWithDatabaseSet, []string{"--database_username=t", "--database_name=t"}...)...)
		assert.NoError(t, err, "expected no error running manager")
		assert.Equal(t, errorsWereFoundVerifying+`  - Mongo database name is not long enough
  - Mongo database username is not long enough`, out, "print to screen did not match expected error(s)")
	})

	t.Run("prints secrets from config", func(t *testing.T) {
		out, err := execute(t, NewRootCmd(), append(argsWithDatabaseSet, []string{"--print", "--secrets"}...)...)
		assert.NoError(t, err, "expected no error running manager")
		assert.Contains(t, out, "username: test_username", "expected secrets to be shown in config")
		assert.Contains(t, out, "password: test_password", "expected secrets to be shown in config")
	})

	t.Run("prints config to screen", func(t *testing.T) {
		out, err := execute(t, NewRootCmd(), append(argsWithDatabaseSet, []string{"--print"}...)...)
		assert.NoError(t, err, "expected no error running manager")
		assert.Equal(t, `OK

log_level: debug
timeouts:
  startup: 30s
  shutdown: 30s
http:
  listen_address: :8080
  cors:
    allowed_origins:
    - '*'
database:
  host: mongodb://ugcompsoc_apid_local_db
  name: apid
  username: '********'
  password: '********'`, out, "print to screen did not match expected config")
	})

	err := os.RemoveAll(directory)
	assert.NoError(t, err, "could not delete testing directory")
}
