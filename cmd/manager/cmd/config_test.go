package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/ugcompsoc/apid/internal/config"
	"gopkg.in/yaml.v2"
)

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()

	return strings.TrimSpace(buf.String()), err
}

func TestVerifyConfig(t *testing.T) {
	directory := "config_test"
	if _, err := os.Stat(directory); err != nil {
		err := os.Mkdir(directory, os.ModePerm)
		assert.NoError(t, err, "could not create test directory")
	}
	absoluteFilePath := filepath.Join(directory, "apid.yml")

	argsWithDirectory := []string{"config", "--directory=" + directory}

	runs := []struct {
		name       string
		args       []string
		out        string
		BeforeWork func()
	}{
		{
			name:       "prints file not found error",
			args:       argsWithDirectory,
			out:        "No file exists at path: config_test/apid.yml",
			BeforeWork: func() {},
		},
		{
			name:       "prints file not found error and debug error",
			args:       append(argsWithDirectory, []string{"--debug"}...),
			out:        "Error: open config_test/apid.yml: no such file or directory\nNo file exists at path: config_test/apid.yml",
			BeforeWork: func() {},
		},
		{
			name: "prints file is empty",
			args: argsWithDirectory,
			out:  "An error was encountered while verifing the file",
			BeforeWork: func() {
				_, err := os.Create(absoluteFilePath)
				assert.NoError(t, err, "expected no error in creating testing file")
			},
		},
		{
			name: "prints file is empty with debug flag",
			args: append(argsWithDirectory, []string{"--debug"}...),
			out:  "Error: The file apid.yml is completely empty. What do you want me to do with this?\nAn error was encountered while verifing the file",
			BeforeWork: func() {
				_, err := os.Create(absoluteFilePath)
				assert.NoError(t, err, "expected no error in creating testing file")
			},
		},
		{
			name: "prints could not unmarshal if file contains a malformed variable",
			args: argsWithDirectory,
			out:  "An error was encountered while verifing the file",
			BeforeWork: func() {
				file, err := os.Create(absoluteFilePath)
				file.Write([]byte("b"))
				assert.NoError(t, err, "expected no error in creating testing file")
			},
		},
		{
			name: "prints could not unmarshal if file contains a malformed variable and debug error",
			args: append(argsWithDirectory, []string{"--debug"}...),
			out:  "Error: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `b` into config.Config\nAn error was encountered while verifing the file",
			BeforeWork: func() {
				file, err := os.Create(absoluteFilePath)
				file.Write([]byte("b"))
				assert.NoError(t, err, "expected no error in creating testing file")
			},
		},
		{
			name:       "an invalid filename will cause an issue",
			args:       append(argsWithDirectory, "--filename=apid"),
			out:        "The filename is not in the form [NAME].yml",
			BeforeWork: func() {},
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			run.BeforeWork()
			out, err := execute(t, NewRootCmd(), run.args...)
			assert.NoError(t, err, "expected no error running manager")
			assert.Equal(t, run.out, out, "unexpected manager output")
			if _, err := os.Stat(directory + "/apid.yml"); err == nil {
				os.Remove(directory + "/apid.yml")
			}
		})
	}

	// setup
	c := &config.Config{
		LogLevel: "debug",
		Timeouts: config.Timeouts{
			Shutdown: 30 * time.Second,
			Startup:  30 * time.Second,
		},
		HTTP: config.HTTP{
			ListenAddress: ":8080",
			CORS: config.CORS{
				AllowedOrigins: []string{"*"},
			},
		},
		Database: config.Database{
			Host:     "mongodb://ugcompsoc_apid_local_db",
			Name:     "apid",
			Username: "test_username",
			Password: "test_password",
		},
	}

	cYaml, err := yaml.Marshal(c)
	assert.NoError(t, err, "expected no error marshalling config to yaml")
	err = ioutil.WriteFile(absoluteFilePath, cYaml, 0644)
	assert.NoError(t, err, "expected no error writing file")

	t.Run("prints ok to screen", func(t *testing.T) {
		out, err := execute(t, NewRootCmd(), argsWithDirectory...)
		assert.NoError(t, err, "expected no error running manager")
		assert.Equal(t, `OK`, out, "print to screen did not match expected config")
	})

	t.Run("prints secrets from config", func(t *testing.T) {
		out, err := execute(t, NewRootCmd(), append(argsWithDirectory, []string{"--print", "--secrets"}...)...)
		assert.NoError(t, err, "expected no error running manager")
		assert.Contains(t, out, "username: test_username", "expected secrets to be shown in config")
		assert.Contains(t, out, "password: test_password", "expected secrets to be shown in config")
	})

	t.Run("prints config to screen", func(t *testing.T) {
		out, err := execute(t, NewRootCmd(), append(argsWithDirectory, "--print")...)
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

	t.Run("prints error(s) to screen", func(t *testing.T) {
		c.Database.Name = ""
		cYaml, err := yaml.Marshal(c)
		assert.NoError(t, err, "expected no error marshalling config to yaml")
		err = ioutil.WriteFile(absoluteFilePath, cYaml, 0644)
		out, err := execute(t, NewRootCmd(), append(argsWithDirectory, "--secrets")...)
		assert.Equal(t, `Error(s) were found while parsing apid.yml, view them below and address them
  - Mongo database name is not long enough`, out, "print to screen did not match expected error(s)")

		c.Database.Username = ""
		cYaml, err = yaml.Marshal(c)
		assert.NoError(t, err, "expected no error marshalling config to yaml")
		err = ioutil.WriteFile(absoluteFilePath, cYaml, 0644)
		out, err = execute(t, NewRootCmd(), append(argsWithDirectory, "--secrets")...)
		assert.Equal(t, `Error(s) were found while parsing apid.yml, view them below and address them
  - Mongo database name is not long enough
  - Mongo database username is not long enough`, out, "print to screen did not match expected error(s)")
	})

	err = os.RemoveAll(directory)
	assert.NoError(t, err, "could not delete testing directory")
}
