package config

import (
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetZeroLogLevel(t *testing.T) {
	t.Run("check all levels can be converted", func(t *testing.T) {
		c := &Config{}

		c.LogLevel = "trace"
		assert.Equal(t, zerolog.TraceLevel, c.GetZeroLogLevel(), "expected to return trace zerolog ENUM")

		c.LogLevel = "disabled"
		assert.Equal(t, zerolog.Disabled, c.GetZeroLogLevel(), "expected to return disabled zerolog ENUM")

		c.LogLevel = "panic"
		assert.Equal(t, zerolog.PanicLevel, c.GetZeroLogLevel(), "expected to return panic zerolog ENUM")

		c.LogLevel = "fatal"
		assert.Equal(t, zerolog.FatalLevel, c.GetZeroLogLevel(), "expected to return fatal zerolog ENUM")

		c.LogLevel = "error"
		assert.Equal(t, zerolog.ErrorLevel, c.GetZeroLogLevel(), "expected to return error zerolog ENUM")

		c.LogLevel = "warn"
		assert.Equal(t, zerolog.WarnLevel, c.GetZeroLogLevel(), "expected to return warn zerolog ENUM")

		c.LogLevel = "info"
		assert.Equal(t, zerolog.InfoLevel, c.GetZeroLogLevel(), "expected to return info zerolog ENUM")

		c.LogLevel = "debug"
		assert.Equal(t, zerolog.DebugLevel, c.GetZeroLogLevel(), "expected to return debug zerolog ENUM")

		c.LogLevel = "dummy"
		assert.Equal(t, zerolog.NoLevel, c.GetZeroLogLevel(), "expected to return nolevel zerolog ENUM")
	})
}

var validConfig Config = Config{
	LogLevel: "debug",
	Timeouts: Timeouts{
		Shutdown: 30 * time.Second,
		Startup:  30 * time.Second,
	},
	HTTP: HTTP{
		ListenAddress: ":8080",
		CORS: CORS{
			AllowedOrigins: []string{"*"},
		},
	},
	Database: Database{
		Host:     "mongodb://ugcompsoc_apid_local_db",
		Name:     "apid",
		Username: "test_username",
		Password: "test_password",
	},
}

type Run struct {
	name        string
	beforeWork  func()
	verifyFunc  func() ([]string, error)
	issue       string
	expectIssue bool
}

func (r *Run) verifyIssuesAndError(t *testing.T) {
	r.beforeWork()
	issues, err := r.verifyFunc()
	// Ideally we should check here if we should expect an error although we should
	// never get an error as the regex is static. If the regex is changed and is invalid
	// this will catch it.
	assert.NoError(t, err, "expected no error from verify function")
	if r.expectIssue {
		assert.Contains(t, issues, r.issue, "expected an issue but was not found")
	} else {
		assert.NotContains(t, issues, r.issue, "expected no issue but was found")
	}
}

func TestTimeoutsVerify(t *testing.T) {
	var testConfig Config

	runs := []Run{
		// Timeouts: Startup
		{
			name: "expect no startup timeout error",
			beforeWork: func() {
				testConfig.Timeouts.Startup = 30 * time.Second
			},
			issue:       "Startup timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: false,
		},
		{
			name: "expect a startup timeout error when an invalid time duration is given",
			beforeWork: func() {
				testConfig.Timeouts.Startup = 30 * time.Hour
			},
			issue:       "Startup timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: true,
		},
		{
			name: "expect a startup timeout error when no time duration is given",
			beforeWork: func() {
				testConfig.Timeouts.Startup = 0
			},
			issue:       "Startup timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: true,
		},
		// Timeouts: Shutdown
		{
			name: "expect no shutdown timeout error",
			beforeWork: func() {
				testConfig.Timeouts.Shutdown = 30 * time.Second
			},
			issue:       "Shutdown timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: false,
		},
		{
			name: "expect a shutdown timeout error when an invalid time duration is given",
			beforeWork: func() {
				testConfig.Timeouts.Shutdown = 30 * time.Hour
			},
			issue:       "Shutdown timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: true,
		},
		{
			name: "expect a shutdown timeout error when no time duration is given",
			beforeWork: func() {
				testConfig.Timeouts.Shutdown = 0
			},
			issue:       "Shutdown timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: true,
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			testConfig = validConfig
			run.verifyFunc = testConfig.Timeouts.Verify
			run.verifyIssuesAndError(t)
		})
	}
}

func TestHTTPVerify(t *testing.T) {
	var testConfig Config

	runs := []Run{
		// HTTP: Listen Address
		{
			name:        "expect no HTTP listen address issue",
			beforeWork:  func() {},
			issue:       "HTTP listen address is not valid",
			expectIssue: false,
		},
		{
			name: "expect HTTP listen address issue when given nothing",
			beforeWork: func() {
				testConfig.HTTP.ListenAddress = ""
			},
			issue:       "HTTP listen address is not valid",
			expectIssue: true,
		},
		{
			name: "expect HTTP listen address issue when given nothing",
			beforeWork: func() {
				testConfig.HTTP.ListenAddress = "...abc"
			},
			issue:       "HTTP listen address is not valid",
			expectIssue: true,
		},
		// HTTP: CORS Allowed Origins
		{
			name:        "expect no HTTP CORS issue with * allowed",
			beforeWork:  func() {},
			issue:       "The allowed origin * is invalid",
			expectIssue: false,
		},
		{
			name: "expect no HTTP CORS issue with domain allowed",
			beforeWork: func() {
				testConfig.HTTP.CORS.AllowedOrigins = []string{"api.compsoc.ie"}
			},
			issue:       "The allowed origin api.compsoc.ie is invalid",
			expectIssue: false,
		},
		{
			name: "expect no HTTP CORS issue with domain and * allowed",
			beforeWork: func() {
				testConfig.HTTP.CORS.AllowedOrigins = []string{"api.compsoc.ie", "*"}
			},
			issue:       "The allowed origin",
			expectIssue: false,
		},
		{
			name: "expect HTTP CORS issue when none is given",
			beforeWork: func() {
				testConfig.HTTP.CORS.AllowedOrigins = []string{}
			},
			issue:       "No allowed origins specified",
			expectIssue: true,
		},
		{
			name: "expect HTTP CORS issue when invalid domain is given",
			beforeWork: func() {
				testConfig.HTTP.CORS.AllowedOrigins = []string{"apid.compsoc.i"}
			},
			issue:       "The allowed origin apid.compsoc.i is invalid",
			expectIssue: true,
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			testConfig = validConfig
			run.verifyFunc = testConfig.HTTP.Verify
			run.verifyIssuesAndError(t)
		})
	}
}

func TestDatabaseVerify(t *testing.T) {
	var testConfig Config

	runs := []Run{
		// Mongo
		{
			name:        "expect no Mongo host issue",
			beforeWork:  func() {},
			issue:       "Mongo host is not valid",
			expectIssue: false,
		},
		{
			name: "expect Mongo issues when host is empty",
			beforeWork: func() {
				testConfig.Database.Host = ""
			},
			issue:       "Mongo host is not valid",
			expectIssue: true,
		},
		{
			name: "expect Mongo issues when malformed host",
			beforeWork: func() {
				testConfig.Database.Host = "mongo://test"
			},
			issue:       "Mongo host is not valid",
			expectIssue: true,
		},
		{
			name: "expect Mongo database name issue when not set",
			beforeWork: func() {
				testConfig.Database.Name = ""
			},
			issue:       "Mongo database name is not long enough",
			expectIssue: true,
		},
		{
			name: "expect Mongo database name issue is less than 3 characters",
			beforeWork: func() {
				testConfig.Database.Name = "ab"
			},
			issue:       "Mongo database name is not long enough",
			expectIssue: true,
		},
		{
			name: "expect no Mongo database name issue",
			beforeWork: func() {
				testConfig.Database.Name = "abc"
			},
			issue:       "Mongo database name is not long enough",
			expectIssue: false,
		},
		{
			name: "expect Mongo database username issue when not set",
			beforeWork: func() {
				testConfig.Database.Username = ""
			},
			issue:       "Mongo database username is not long enough",
			expectIssue: true,
		},
		{
			name: "expect Mongo database username issue is less than 3 characters",
			beforeWork: func() {
				testConfig.Database.Username = "ab"
			},
			issue:       "Mongo database username is not long enough",
			expectIssue: true,
		},
		{
			name: "expect no Mongo database username issue",
			beforeWork: func() {
				testConfig.Database.Username = "abc"
			},
			issue:       "Mongo database username is not long enough",
			expectIssue: false,
		},
		{
			name: "expect Mongo database password issue when not set",
			beforeWork: func() {
				testConfig.Database.Password = ""
			},
			issue:       "Mongo database password is not long enough",
			expectIssue: true,
		},
		{
			name: "expect Mongo database password issue is less than 3 characters",
			beforeWork: func() {
				testConfig.Database.Password = "ab"
			},
			issue:       "Mongo database password is not long enough",
			expectIssue: true,
		},
		{
			name: "expect no Mongo database password issue",
			beforeWork: func() {
				testConfig.Database.Password = "abc"
			},
			issue:       "Mongo database password is not long enough",
			expectIssue: false,
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			testConfig = validConfig
			run.verifyFunc = testConfig.Database.Verify
			run.verifyIssuesAndError(t)
		})
	}
}

func TestConfig(t *testing.T) {
	var testConfig Config

	t.Run("expect no issues or error", func(t *testing.T) {
		testConfig := validConfig
		issues, err := testConfig.Verify()
		assert.NoError(t, err, "expected no error from verify function")
		assert.Empty(t, issues, "expected no issues from verify function")
	})

	t.Run("expect multiple issues", func(t *testing.T) {
		testConfig := validConfig
		testConfig.LogLevel = ""
		testConfig.Timeouts.Startup = 0
		issues, err := testConfig.Verify()
		assert.NoError(t, err)
		assert.Contains(t, issues, "An invalid log level was specified", "expected an issue but was not found")
		assert.Contains(t, issues, "Startup timeout should be represented in the form '{int}s', e.g. '30s'", "expected an issue but was not found")
	})

	runs := []Run{
		// Log level
		{
			name:        "expect no log level error",
			beforeWork:  func() {},
			issue:       "An invalid log level was specified",
			expectIssue: false,
		},
		{
			name: "expect log level error when invalid log level is given",
			beforeWork: func() {
				testConfig.LogLevel = "a"
			},
			issue:       "An invalid log level was specified",
			expectIssue: true,
		},
		{
			name: "expect log level error when no LogLevel is not given",
			beforeWork: func() {
				testConfig.LogLevel = ""
			},
			issue:       "An invalid log level was specified",
			expectIssue: true,
		},
		// Timeouts issues retrieved sanity check
		{
			name: "expect timeouts issue to exist",
			beforeWork: func() {
				testConfig.Timeouts.Startup = 0
			},
			issue:       "Startup timeout should be represented in the form '{int}s', e.g. '30s'",
			expectIssue: true,
		},
		// HTTP issues retrieved sanity check
		{
			name: "expect http issue to exist",
			beforeWork: func() {
				testConfig.HTTP.ListenAddress = ""
			},
			issue:       "HTTP listen address is not valid",
			expectIssue: true,
		},
		// Database issues retrieved sanity check
		{
			name: "expect database issue to exist",
			beforeWork: func() {
				testConfig.Database.Name = ""
			},
			issue:       "Mongo database name is not long enough",
			expectIssue: true,
		},
	}

	for _, run := range runs {
		t.Run(run.name, func(t *testing.T) {
			testConfig = validConfig
			run.verifyFunc = testConfig.Verify
			run.verifyIssuesAndError(t)
		})
	}
}
