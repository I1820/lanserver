package store

import (
	"context"

	"github.com/I1820/lanserver/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Collection = "devices"

type Device struct {
	DB *mongo.Database
}

func (d Device) Get(ctx context.Context) ([]model.Device, error) {
	var results = make([]model.Device, 0)

	cur, err := d.DB.Collection(Collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var result model.Device

		if err := cur.Decode(&result); err != nil {
			return nil, err
		}

		results = append(results, result)
	}
	if err := cur.Close(ctx); err != nil {
		return nil, err
	}

	return results, nil
}

func (d Device) Show(ctx context.Context, deviceID string) (model.Device, error) {
	var dev model.Device

	result := d.DB.Collection(Collection).FindOne(ctx, bson.M{
		"deveui": deviceID,
	})

	if err := result.Decode(&dev); err != nil {
		return dev, err
	}

	return dev, nil
}

func (d Device) Insert(ctx context.Context, dev model.Device) error {
	if _, err := d.DB.Collection(Collection).InsertOne(ctx, dev); err != nil {
		return err
	}
	return nil
}

func (d Device) Destroy(ctx context.Context, deviceID string) (model.Device, error) {
	var dev model.Device

	result := d.DB.Collection(Collection).FindOneAndDelete(ctx, bson.M{
		"deveui": deviceID,
	})

	if err := result.Decode(&dev); err != nil {
		return dev, err
	}

	return dev, nil
}

func (d Device) Update(ctx context.Context, deviceID string, field string, value interface{}) (model.Device, error) {
	var dev model.Device

	res := d.DB.Collection(Collection).FindOneAndUpdate(ctx, bson.M{
		"deveui": deviceID,
	}, bson.M{
		"$set": bson.M{
			field: value,
		},
	}, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if err := res.Decode(&dev); err != nil {
		return dev, err
	}

	return dev, nil
}
