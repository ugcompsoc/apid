package database

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	db_utils "github.com/ugcompsoc/apid/internal/services/database_test_utils"

	"github.com/stretchr/testify/assert"
	"github.com/ugcompsoc/apid/internal/config"

	"github.com/containers/common/pkg/retry"
	"github.com/ory/dockertest/v3"
)

var dockerPool *dockertest.Pool
var dockerDefaultResource *dockertest.Resource
var dockerResourcesToPurge []*dockertest.Resource
var ds *Datastore

func TestMain(m *testing.M) {
	// Recover from panic so that container can be purged
	defer func() {
		db_utils.PurgeDockerResources(dockerPool, dockerResourcesToPurge)
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()
	// Setup
	var err error
	dockerPool, err = db_utils.SetupDockerPool()
	if err != nil {
		log.Fatalf("could not connect to docker %s", err)
	}
	config := &config.Config{}
	port, err := db_utils.GetFreePort()
	if err != nil {
		log.Fatalf("issue resolving a free port: %s", err)
	}
	config.Database.Host = "mongodb://localhost:" + strconv.Itoa(port)
	config.Database.Name = "test_database"
	config.Database.Username = "test_username"
	config.Database.Password = "test_password"
	dockerDefaultResource, err = db_utils.SetupMongoDocker(config, dockerPool)
	if err != nil {
		log.Fatalf("could not create mongo docker %s", err)
	}
	dockerResourcesToPurge = append(dockerResourcesToPurge, dockerDefaultResource)
	err = retry.RetryIfNecessary(context.Background(), func() error {
		var err error
		ds, err = NewDatastore(config)
		return err
	}, &retry.Options{
		MaxRetry: 10,
		Delay:    100 * time.Millisecond,
	})
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	// Run tests
	exitCode := m.Run()

	// Teardown
	// purge all resources
	db_utils.PurgeDockerResources(dockerPool, dockerResourcesToPurge)

	// Exit
	os.Exit(exitCode)
}

func TestNewDatastore(t *testing.T) {
	t.Parallel()

	// Setup
	defaultConfig := &config.Config{}
	port, err := db_utils.GetFreePort()
	assert.NoError(t, err)
	defaultConfig.Database.Host = "mongodb://localhost:" + strconv.Itoa(port)
	defaultConfig.Database.Name = "test_database"
	defaultConfig.Database.Username = "test_username"
	defaultConfig.Database.Password = "test_password"
	resource, err := db_utils.SetupMongoDocker(defaultConfig, dockerPool)
	assert.NoError(t, err, "could not create mongo docker")
	dockerResourcesToPurge = append(dockerResourcesToPurge, resource)

	t.Run("can connect and ping", func(t *testing.T) {
		err = retry.RetryIfNecessary(context.Background(), func() error {
			var err error
			_, err = NewDatastore(defaultConfig)
			return err
		}, &retry.Options{
			MaxRetry: 10,
			Delay:    100 * time.Millisecond,
		})
		assert.NoError(t, err, "could not connect to docker")
	})

	t.Run("failed to connect to database if the host is wrong", func(t *testing.T) {
		config := &config.Config{}
		_, err := NewDatastore(config)
		assert.True(t, strings.Contains(err.Error(), "failed to connect to/create session with database host"), "could not connect to docker")
	})

	t.Run("cant ping database if the credentials are wrong", func(t *testing.T) {
		config := defaultConfig
		err = retry.RetryIfNecessary(context.Background(), func() error {
			var err error
			config.Database.Password = ""
			_, err = NewDatastore(config)
			return err
		}, &retry.Options{
			MaxRetry: 10,
			Delay:    100 * time.Millisecond,
		})
		assert.True(t, strings.Contains(err.Error(), "failed to ping the database host"), "could not connect to docker")
	})
}
