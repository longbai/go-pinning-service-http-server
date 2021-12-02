package model

import (
	"context"
	"log"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type db struct {
	*qmgo.Database
	cli *qmgo.Client
}

var store *db

func IdGen() string {
	return qmgo.NewObjectID().Hex()
}

func OpenDb(addr string) {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: addr})
	if err != nil {
		log.Fatal(err)
	}
	d := client.Database("pin")
	store = &db{d, client}
}

func PinAdd(ctx context.Context, status *PinStatus) error {
	_, err := store.Collection("pin").InsertOne(ctx, status)
	return err
}

func PinGet(ctx context.Context, reqId string) (*PinStatus, error) {
	one := PinStatus{}
	err := store.Collection("pin").Find(ctx, bson.M{"requestid": reqId}).One(&one)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &one, nil
}

func PinCount(ctx context.Context) (int64, error) {
	filter := bson.M{}
	count, err := store.Collection("pin").Find(ctx, filter).Count()
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func PinList(ctx context.Context, offset, limit int64) ([]PinStatus, error) {
	batch := []PinStatus{}
	filter := bson.M{}

	err := store.Collection("pin").Find(ctx, filter).Skip(offset).Limit(limit).All(&batch)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return batch, nil
}

func PinUpdate(ctx context.Context, status *PinStatus) error {
	return store.Collection("pin").UpdateOne(ctx, bson.M{"requestid": status.Requestid}, bson.M{"$set": status})
}

func PinDelete(ctx context.Context, reqId string) error {
	return store.Collection("pin").Remove(ctx, bson.M{"requestid": reqId})
}

func CloseDb() {
	if store != nil {
		err := store.cli.Close(context.Background())
		log.Println(err)
	}
}
