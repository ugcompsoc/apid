package config

import (
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/rs/zerolog"
)

var domainRegexStr string = "([a-z0-9]+(-[a-z0-9]+)*.)+[a-z]{2,}"

type Timeouts struct {
	Startup  time.Duration `mapstructure:"startup" yaml:"startup"`
	Shutdown time.Duration `mapstructure:"shutdown" yaml:"shutdown"`
}

type CORS struct {
	AllowedOrigins []string `mapstructure:"allowed_origins" yaml:"allowed_origins"`
}

type HTTP struct {
	ListenAddress string `mapstructure:"listen_address" yaml:"listen_address"`
	CORS          CORS
}

type Database struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Name     string `mapstructure:"name" yaml:"name"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

// Config describes the configuration for Server
type Config struct {
	LogLevel string `mapstructure:"log_level" yaml:"log_level"`
	Timeouts Timeouts
	HTTP     HTTP
	Database Database
}

func (c *Config) GetZeroLogLevel() zerolog.Level {
	switch c.LogLevel {
	case "trace":
		return zerolog.TraceLevel
	case "disabled":
		return zerolog.Disabled
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.NoLevel
	}
}

func (t *Timeouts) Verify() ([]string, error) {
	issues := []string{}
	timeoutsRegex, err := regexp.Compile("^[0-9]{1,2}s$")
	if err != nil {
		return nil, err
	}
	if !timeoutsRegex.MatchString(t.Startup.String()) || t.Startup.String() == "0s" {
		issues = append(issues, "Startup timeout should be represented in the form '{int}s', e.g. '30s'")
	}
	if !timeoutsRegex.MatchString(t.Shutdown.String()) || t.Shutdown.String() == "0s" {
		issues = append(issues, "Shutdown timeout should be represented in the form '{int}s', e.g. '30s'")
	}
	return issues, nil
}

func (h *HTTP) Verify() ([]string, error) {
	issues := []string{}
	_, _, err := net.SplitHostPort(h.ListenAddress)
	if err != nil {
		issues = append(issues, "HTTP listen address is not valid")
	}
	httpDomainRegex, err := regexp.Compile("^(?:https?://)?" + domainRegexStr + "$")
	if err != nil {
		return nil, err
	}
	if len(h.CORS.AllowedOrigins) == 0 {
		issues = append(issues, "No allowed origins specified")
	}
	for _, origin := range h.CORS.AllowedOrigins {
		if origin != "*" && !httpDomainRegex.MatchString(origin) {
			issues = append(issues, fmt.Sprintf("The allowed origin %s is invalid", origin))
		}
	}
	return issues, nil
}

func (d *Database) Verify() ([]string, error) {
	issues := []string{}
	mongoHostRegex, err := regexp.Compile("^mongodb://" + domainRegexStr + "$")
	if err != nil {
		return nil, err
	}
	if !mongoHostRegex.MatchString(d.Host) {
		issues = append(issues, "Mongo host is not valid")
	}
	if len(d.Name) < 3 {
		issues = append(issues, "Mongo database name is not long enough")
	}
	if len(d.Username) < 3 {
		issues = append(issues, "Mongo database username is not long enough")
	}
	if len(d.Password) < 3 {
		issues = append(issues, "Mongo database password is not long enough")
	}
	return issues, nil
}

func (c *Config) Verify() ([]string, error) {
	issues := []string{}

	if c.GetZeroLogLevel() == zerolog.NoLevel {
		issues = append(issues, "An invalid log level was specified")
	}

	timeoutIssues, err := c.Timeouts.Verify()
	if err != nil {
		return nil, err
	}
	issues = append(issues, timeoutIssues...)

	httpIssues, err := c.HTTP.Verify()
	if err != nil {
		return nil, err
	}
	issues = append(issues, httpIssues...)

	databaseIssues, err := c.Database.Verify()
	if err != nil {
		return nil, err
	}
	issues = append(issues, databaseIssues...)

	return issues, nil
}
