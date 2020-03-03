package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/surajjain36/log_server/misc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Mongo Struct
type Mongo struct {
	db *mongo.Database
}

//NewMongo is a function which establishes connection to MongoDB.
func NewMongo(conf *misc.MongoConfig) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.DB)))

	if err != nil {
		log.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	err = client.Ping(ctx, readpref.Primary())
	log.Println(err)
	if err == nil {
		return &Mongo{
			db: client.Database(conf.DB),
		}, nil
	}

	return nil, err
}

//Insert is a generic function to do insert call on DB
func (mgo *Mongo) Insert(collection string, data interface{}) (interface{}, error) {
	var err error
	var res interface{}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if res, err = mgo.db.Collection(collection).InsertOne(ctx, data); err != nil {
		log.Println("Error while inserting record: ", err)
	}
	return res, err
}

//Aggregate : Generic function to get aggregated data from DB
func (mgo *Mongo) Aggregate(collection string, pipeline interface{}, result interface{}) error {
	var err error
	var cursor *mongo.Cursor
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	opts := options.Aggregate()
	opts.SetAllowDiskUse(true)
	opts.SetBatchSize(5)
	if cursor, err = mgo.db.Collection(collection).Aggregate(ctx, pipeline, opts); err == nil && cursor != nil {
		defer cursor.Close(ctx)
		err = cursor.All(ctx, result)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}
