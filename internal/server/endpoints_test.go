package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ugcompsoc/apid/internal/services/database"
	db_utils "github.com/ugcompsoc/apid/internal/services/database_test_utils"

	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"

	"github.com/ugcompsoc/apid/internal/config"

	"github.com/containers/common/pkg/retry"
)

var dockerPool *dockertest.Pool
var dockerDefaultResource *dockertest.Resource
var dockerResourcesToPurge []*dockertest.Resource
var ds *database.Datastore

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
		log.Fatalf("could not connect to docker: %s", err)
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
		ds, err = database.NewDatastore(config)
		return err
	}, &retry.Options{
		MaxRetry: 3,
		Delay:    2 * time.Second,
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

func TestRootV2Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v2", s.RootV2Get)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status 200 from endpoint")
		assert.Equal(t, "{\"message\":\"Root V2\"}", w.Body.String(), "unexpected response")
	})
}

func TestMiscV2HealthcheckGet(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/healthcheck", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{
			Datastore: ds,
		}
		engine.GET("/v2/healthcheck", s.MiscV2HealthcheckGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status 200 from endpoint")
		assert.Equal(t, "{\"errors\":[]}", w.Body.String(), "expected empty errors array")
	})

	t.Run("expect database issue", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/healthcheck", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{
			Datastore: ds,
		}
		err = s.Datastore.Client.Disconnect(context.TODO())
		assert.NoError(t, err, "could not disconnect from database")
		engine.GET("/v2/healthcheck", s.MiscV2HealthcheckGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status 200 from endpoint")
		assert.Equal(t, "{\"errors\":[\"cannot ping database\"]}", w.Body.String(), "expected empty errors array")
	})
}

func TestBrewV2Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/brew", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v2/brew", s.MiscV2BrewGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTeapot, w.Code, "expected status 418 from endpoint")
		assert.Equal(t, "{\"error\":\"I refuse to brew coffee because I am, permanently, a teapot.\"}", w.Body.String(), "unexpected response")
	})
}

func TestMiscV2Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/ping", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v2/ping", s.MiscV2PingGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status 200 from endpoint")
		assert.Equal(t, "{\"message\":\"Pong!\"}", w.Body.String(), "unexpected response")
	})
}
