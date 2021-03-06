package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/I1820/lanserver/actions"
	"github.com/I1820/lanserver/config"
	"github.com/I1820/lanserver/db"
	"github.com/I1820/lanserver/node"
	"github.com/I1820/lanserver/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	// ExitTimeout is a time that application waits for API service to exit
	ExitTimeout = 5 * time.Second
)

func main(cfg config.Config) {
	db, err := db.New(cfg.Database)
	if err != nil {
		logrus.Fatalf("db new client error: %s", err)
	}

	st := store.Device{DB: db}

	app := actions.App(cfg.Debug, st)

	go func() {
		if err := app.Start(":4000"); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)
		}
	}()

	if _, err := node.New(cfg.App.Broker.Addr, cfg.Node.Broker.Addr, st); err != nil {
		logrus.Fatalf("API Service failed with %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), ExitTimeout)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		logrus.Printf("API Service failed on exit: %s", err)
	}
}

// Register server command
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server",
			Short: "Run server to serve the requests",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
