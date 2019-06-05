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
	"github.com/I1820/lanserver/node"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("18.20 at Sep 07 2016 7:20 IR721")

	cfg := config()

	// create mongodb connection
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Database.URL))
	if err != nil {
		logrus.Fatalf("db new client error: %s", err)
	}
	// connect to the mongodb (change database here!)
	ctxc, donec := context.WithTimeout(context.Background(), 10*time.Second)
	defer donec()
	if err := client.Connect(ctxc); err != nil {
		logrus.Fatalf("db connection error: %s", err)
	}
	// is the mongo really there?
	ctxp, donep := context.WithTimeout(context.Background(), 2*time.Second)
	defer donep()
	if err := client.Ping(ctxp, readpref.Primary()); err != nil {
		logrus.Fatalf("db ping error: %s", err)
	}
	db := client.Database("lanserver")

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
