package db

import (
	"context"
	"fmt"
	"time"

	"github.com/I1820/lanserver/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	// ConnectionTimeout database connection timeout
	ConnectionTimeout = 10 * time.Second
	// PingTimeout database ping timeout
	PingTimeout = 2 * time.Second
)

// New creates a new mongodb connection and tests it
func New(cfg config.Database) (*mongo.Database, error) {
	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.URL))
	if err != nil {
		return nil, fmt.Errorf("db new client error: %s", err)
	}

	// connect to the mongodb
	ctxc, donec := context.WithTimeout(context.Background(), ConnectionTimeout)
	defer donec()

	if err := client.Connect(ctxc); err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	// is the mongo really there?
	ctxp, donep := context.WithTimeout(context.Background(), PingTimeout)
	defer donep()

	if err := client.Ping(ctxp, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	return client.Database(cfg.Name), nil
}
