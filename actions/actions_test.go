package actions

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// LSTestSuite is a test suite for lanserver component APIs.
type LSTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates lanserver test suite
func (suite *LSTestSuite) SetupSuite() {
	url := os.Getenv("I1820_LANSERVER_DATABASE_URL")
	if url == "" {
		url = "mongodb://127.0.0.1:27017"
	}

	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		suite.NoError(err)
	}
	// connect to the mongodb (change database here!)
	ctxc, donec := context.WithTimeout(context.Background(), 10*time.Second)
	defer donec()
	if err := client.Connect(ctxc); err != nil {
		suite.NoError(err)
	}
	// is the mongo really there?
	ctxp, donep := context.WithTimeout(context.Background(), 2*time.Second)
	defer donep()
	if err := client.Ping(ctxp, readpref.Primary()); err != nil {
		suite.NoError(err)
	}
	db := client.Database("lanserver")

	suite.engine = App(true, db)
}

// Let's test lanserver APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(LSTestSuite))
}
