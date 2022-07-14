package mng

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EntityScheme interface {
	GetCollectionName() string
	CreateIndexes(coll *mongo.Collection) ([]string, error)
}

type DbSet[T any] struct {
	Scheme           EntityScheme
	ConnectionString string

	conn     *mongo.Client
	feniksDb *mongo.Database
	coll     *mongo.Collection
}

func (d *DbSet[T]) Open() error {
	uri := d.ConnectionString
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	d.conn = client
	d.feniksDb = client.Database("feniks")
	d.coll = d.feniksDb.Collection(d.Scheme.GetCollectionName())

	return nil
}

func (d *DbSet[T]) Close() error {
	if d.conn == nil {
		return nil
	}
	return d.conn.Disconnect(context.TODO())
}

func (d *DbSet[T]) CreateIndexes() ([]string, error) {
	return d.Scheme.CreateIndexes(d.coll)
}

func (d *DbSet[T]) AddRange(entities []interface{}) error {
	_, err := d.coll.InsertMany(context.TODO(), entities)
	return err
}

func (d *DbSet[T]) GetByQuery(query bson.M) ([]*T, error) {
	ctx := context.TODO()
	cursor, err := d.coll.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	ret := make([]*T, 0)
	err = cursor.All(ctx, &ret)
	return ret, nil
}

func (d *DbSet[T]) GetCollection() *mongo.Collection {
	return d.coll
}
