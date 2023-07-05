package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ugcompsoc/apid/internal/config"
)

var ExampleConfig config.Config = config.Config{
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

var ExampleConfigYAML string = `log_level: debug
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
  username: test_username
  password: test_password
`

func TestPrintConfig(t *testing.T) {
	t.Run("will block out secrets", func(t *testing.T) {
		tempConfig := ExampleConfig
		yamlStr, err := PrintConfig(&tempConfig, false)
		assert.NoError(t, err, "expecting no error while printing config")
		assert.Contains(t, yamlStr, "username: '********'", "expected no secrets to be shown in config")
		assert.Contains(t, yamlStr, "password: '********'", "expected no secrets to be shown in config")
	})

	t.Run("will not be able to marshal config", func(t *testing.T) {
		_, err := PrintConfig(nil, true)
		assert.Error(t, err, "expecting error while printing config")
	})

	t.Run("will print config", func(t *testing.T) {
		tempConfig := ExampleConfig
		yamlStr, err := PrintConfig(&tempConfig, true)
		assert.NoError(t, err, "expecting no error while printing config")
		assert.Equal(t, ExampleConfigYAML, yamlStr, "expected config print was not received")
	})
}

func TestVerifyFilename(t *testing.T) {
	errStr := "The filename is not in the form [NAME].yml"

	t.Run("will fail with empty filename", func(t *testing.T) {
		err := VerifyFilename("")
		assert.Equal(t, errStr, err.Error(), "an unexpected error message was received")
	})

	t.Run("will fail with invalid filename #1", func(t *testing.T) {
		err := VerifyFilename("ttt")
		assert.Equal(t, errStr, err.Error(), "an unexpected error message was received")
	})

	t.Run("will fail with invalid filename #2", func(t *testing.T) {
		err := VerifyFilename("apid.ttt")
		assert.Equal(t, errStr, err.Error(), "an unexpected error message was received")
	})

	t.Run("will pass with valid filename", func(t *testing.T) {
		err := VerifyFilename("apid.yml")
		assert.NoError(t, err, "an unexpected error was received")
	})
}

func TestExtractFile(t *testing.T) {
	directory := "config_util_test"
	if _, err := os.Stat(directory); err != nil {
		err := os.Mkdir(directory, os.ModePerm)
		assert.NoError(t, err, "could not create test directory")
	}
	absoluteFilePath := filepath.Join(directory, "apid.yml")

	t.Run("will fail with empty filename", func(t *testing.T) {
		_, err := ExtractFile("")
		assert.Equal(t, "open : no such file or directory", err.Error(), "an unexpected error message was received")
	})

	t.Run("will fail with a nonexistent file", func(t *testing.T) {
		_, err := ExtractFile(absoluteFilePath)
		assert.Equal(t, "open "+absoluteFilePath+": no such file or directory", err.Error(), "an unexpected error message was received")
	})

	t.Run("will fail with an empty file", func(t *testing.T) {
		os.Create(absoluteFilePath)
		defer os.Remove(absoluteFilePath)
		_, err := ExtractFile(absoluteFilePath)
		assert.Equal(t, "The file apid.yml is completely empty. What do you want me to do with this?", err.Error(), "an unexpected error message was received")
	})

	t.Run("will fail with invalid file contents", func(t *testing.T) {
		err := ioutil.WriteFile(absoluteFilePath, []byte("b"), 0644)
		assert.NoError(t, err, "could not create test file")
		defer os.Remove(absoluteFilePath)
		_, err = ExtractFile(absoluteFilePath)
		assert.Equal(t, "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `b` into config.Config", err.Error(), "an unexpected error message was received")
	})

	t.Run("will fail with invalid file contents", func(t *testing.T) {
		err := ioutil.WriteFile(absoluteFilePath, []byte(ExampleConfigYAML), 0644)
		assert.NoError(t, err, "could not create test file")
		defer os.Remove(absoluteFilePath)
		c, err := ExtractFile(absoluteFilePath)
		assert.NoError(t, err, "expected no error in extracting config from file")
		assert.Equal(t, ExampleConfig, *c)
	})

	err := os.RemoveAll(directory)
	assert.NoError(t, err, "could not delete testing directory")
}
