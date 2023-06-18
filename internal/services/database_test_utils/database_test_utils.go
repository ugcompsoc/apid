package database_test_utils

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ugcompsoc/apid/internal/config"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func SetupDockerPool() (*dockertest.Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not connect to docker: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("Could not connect to docker: %s", err)
	}
	return pool, err
}

func SetupMongoDocker(c *config.Config, pool *dockertest.Pool) (*dockertest.Resource, error) {
	portProto := strings.Split(c.Database.Host, ":")[2] + "/tcp"
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			"MONGO_INITDB_ROOT_USERNAME=" + c.Database.Username,
			"MONGO_INITDB_ROOT_PASSWORD=" + c.Database.Password,
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017/tcp": {{HostIP: "localhost", HostPort: portProto}},
		},
	}, func(dc *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		dc.AutoRemove = true
		dc.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %s", err)
	}
	return resource, nil
}

func PurgeDockerResources(pool *dockertest.Pool, resources []*dockertest.Resource) {
	for _, r := range resources {
		if err := pool.Purge(r); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}
}
