package db

import (
	"fmt"
	"github.com/pingcap/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
)

const (
	MONGO_DB_CONN_POOL	= 128
)

var (
	DB *MongoClient
)

type (
	MongoClient struct {
		client 	*mongo.Client
		colls  	map[string]*mongo.Collection
		uri		string
	}
)

func NewMongoClient(uri string) error {
	opts := options.Client().ApplyURI(uri)
	opts.SetMaxPoolSize(MONGO_DB_CONN_POOL)
	opts.ReadPreference = readpref.Nearest()

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return errors.Trace(err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		return errors.Trace(err)
	}
	DB = &MongoClient{
		client: nil,
		colls:  nil,
		uri:    "",
	}
	DB.client = client
	DB.uri = uri

	return nil
}

//db.usercomplains.find({_serverDate: {$gte: new Date("2020-12-01"), $lt: new Date("2020-12-02")}},{user_message:1})
func (mc *MongoClient) collection(collname string) *mongo.Collection {
	if db := mc.client.Database("gamecluster"); db != nil {
		return db.Collection(collname)
	}
	return nil
}

func (mc *MongoClient) findAndUpdate(collname string, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {
	coll := mc.collection(collname)
	if coll == nil {
		return nil, errors.Errorf("no coll %s", collname)
	}
	return coll.FindOneAndUpdate(context.Background(), filter, update, opts...), nil
}

func (mc *MongoClient) FindAndCreateUser(deviceId string, ) (*User, error) {
	cu := mc.collection("users")
	if cu == nil {
		return nil, errors.Errorf("no users coll")
	}

	u := &User{}
	err := cu.FindOne(context.Background(), bson.M{"deviceid":deviceId}).Decode(u)
	if err != nil {
		cc := mc.collection("counters")
		if cc == nil {
			return nil, errors.Errorf("no counters coll")
		}
		c := &Counter{}
		opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
		if err := cc.FindOneAndUpdate(context.Background(), bson.M{"type":"users"}, bson.M{"$inc":bson.M{"count":1}}, opts).Decode(c); err != nil {
			return nil, errors.Errorf("counter error")
		}

		name := fmt.Sprintf("name-%d", c.Count)
		diamond := 10
		if _, err := cu.InsertOne(context.Background(), bson.M{
			"_id":c.Count,
			"deviceid":deviceId,
			"diamond": diamond,
			"name": name,
		}); err != nil {
			return nil, errors.Trace(err)
		} else {
			u.UId = c.Count
			u.Diamond = diamond
			u.Name = name
		}
	}
	return u, nil
}