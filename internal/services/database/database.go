package database

import (
	"context"
	"fmt"

	"github.com/ugcompsoc/apid/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Datastore struct {
	Client   *mongo.Client
	Database *mongo.Database
}

/*
 *	Database Setup
 */

func NewDatastore(config *config.Config) (*Datastore, error) {
	datastore := &Datastore{}
	if err := datastore.connect(config); err != nil {
		return nil, err
	}
	return datastore, nil
}

func (ds *Datastore) connect(config *config.Config) error {
	credential := options.Credential{
		Username: config.Database.Username,
		Password: config.Database.Password,
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.Database.Host).SetAuth(credential).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return fmt.Errorf("failed to connect to/create session with database host: %w", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("failed to ping the database host: %w", err)
	}

	ds.Client = client
	ds.Database = client.Database(config.Database.Name)

	return err
}
