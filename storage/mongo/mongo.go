package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var pool = map[string]*mongo.Client{}

func Init(ctx context.Context, configs map[string]string) error {
	for name, url := range configs {
		cli, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
		if err != nil {
			return err
		}
		if err := cli.Ping(ctx, readpref.Primary()); err != nil {
			return err
		}
		pool[name] = cli
	}
	return nil
}

func Get(name string) *mongo.Client {
	return pool[name]
}
