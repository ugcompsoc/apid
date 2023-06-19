// Credit
// Thanks to TUD NetSoc for this brilliant entry point.
// https://github.com/netsoc/webspaced

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/ugcompsoc/apid/docs"
	"github.com/ugcompsoc/apid/internal/config"
	"github.com/ugcompsoc/apid/internal/server"
)

var srv *server.Server

func init() {
	// Config defaults
	viper.SetDefault("log_level", zerolog.TraceLevel)

	viper.SetDefault("timeouts.startup", 30*time.Second)
	viper.SetDefault("timeouts.shutdown", 30*time.Second)

	viper.SetDefault("http.listen_address", ":8080")
	viper.SetDefault("http.cors.allowed_origins", []string{"*"})

	// Config file loading
	viper.SetConfigType("yml")
	viper.SetConfigName("apid")
	viper.AddConfigPath("/run/config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")

	// Set logger defaults
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	// Config from flags
	pflag.StringP("log_level", "l", "info", "log level")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Err(err).Msg("failed to bind pflags to config")
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Err(err).Msg("failed to read configuration")
	}
}

func reload() {
	if srv != nil {
		stop()
		srv = nil
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to parse configuration")
	}

	logLevel := cfg.GetZeroLogLevel()
	if logLevel == zerolog.NoLevel {
		log.Fatal().Msg("cannot start server without log level specified")
	}
	zerolog.SetGlobalLevel(logLevel)

	log.Trace().Any("config", cfg).Msg("got config")

	srv = server.NewServer(cfg)

	log.Info().Msg("starting server")
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeouts.Startup)
		defer cancel()
		if err := srv.Start(ctx); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()
}

func stop() {
	log.Info().Msg("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), srv.Config.Timeouts.Shutdown)
	defer cancel()
	if err := srv.Stop(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to stop server")
	}

	log.Info().Msg("stopped server sucessfully")
}

// @title           UG CompSoc APId
// @version         1.0
// @description     Webservices APIv2 for account and IAAS management.
// @termsOfService  https://compsoc.ie/terms
// @contact.name   	UG CompSoc Admin Team
// @contact.url    	https://compsoc.ie/support
// @contact.email  	compsoc@socs.nuigalway.ie
// @license.name  	MIT
// @license.url   	https://github.com/ugcompsoc/apid/blob/main/LICENSE
// @BasePath  		/
func main() {
	// Recover from panic
	defer func() {
		if err := recover(); err != nil {
			log.Fatal().Any("error", err).Msg("server panic")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Trace().Str("file", e.Name).Msg("Config changed, reloading")
		reload()
	})
	viper.WatchConfig()
	reload()

	<-sigs
	stop()
}
