package actions

import (
	"testing"

	"github.com/I1820/lanserver/config"
	"github.com/I1820/lanserver/db"
	"github.com/I1820/lanserver/store"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

// LSTestSuite is a test suite for lanserver component APIs.
type LSTestSuite struct {
	suite.Suite
	engine *echo.Echo
}

// SetupSuite initiates lanserver test suite
func (suite *LSTestSuite) SetupSuite() {
	cfg := config.New()

	db, err := db.New(cfg.Database)
	suite.NoError(err)

	st := store.Device{DB: db}

	suite.engine = App(true, st)
}

// Let's test lanserver APIs!
func TestService(t *testing.T) {
	suite.Run(t, new(LSTestSuite))
}
