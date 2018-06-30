package grifts

import (
	"log"

	"github.com/gobuffalo/envy"
	"github.com/mongodb/mongo-go-driver/bson"
	mgo "github.com/mongodb/mongo-go-driver/mongo"

	grift "github.com/markbates/grift/grift"
)

var _ = grift.Desc("mongo", "Creates mongo database and collections")
var _ = grift.Add("mongo", func(c *grift.Context) error {
	// Create mongodb connection
	url := envy.Get("DB_URL", "mongodb://127.0.0.1")
	client, err := mgo.NewClient(url)
	if err != nil {
		return err
	}
	if err := client.Connect(c); err != nil {
		return err
	}
	log.Printf("DB url: %s\n", url)

	db := client.Database("lanserver")

	// Device collection
	cd := db.Collection("devices")
	name, err := cd.Indexes().CreateOne(
		c,
		mgo.IndexModel{
			Keys: bson.NewDocument(
				bson.EC.Int32("deveui", 1),
			),
			Options: bson.NewDocument(
				bson.EC.Boolean("unique", true),
			),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("DB [devices] index: %s\n", name)

	return nil
})
