package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/I1820/lanserver/actions"
	"github.com/I1820/lanserver/config"
	"github.com/I1820/lanserver/db"
	"github.com/I1820/lanserver/node"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	cfg := config.New()

	db, err := db.New(cfg.Database.URL, "lanserver")
	if err != nil {
		logrus.Fatalf("db new client error: %s", err)
	}

	app := actions.App(cfg.Debug, db)
	go func() {
		if err := app.Start(":4000"); err != http.ErrServerClosed {
			logrus.Fatalf("API Service failed with %s", err)
		}
	}()

	if _, err := node.New(cfg.App.Broker.Addr, cfg.Node.Broker.Addr, db); err != nil {
		logrus.Fatalf("API Service failed with %s", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("18.20 As always ... left me alone")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		logrus.Printf("API Service failed on exit: %s", err)
	}
}
