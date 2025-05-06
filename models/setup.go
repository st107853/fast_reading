package models

import (
	"context"
	"log"

	"github.com/go-redis/redis"
	"github.com/st107853/fast_reading/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBName = "library"
const CollName = "books"

var conf, _ = config.LoadConfig(".")

type Logger struct {
	*mongo.Client // The database access interface
}

var DB Logger

func ConnectToMongoDB() error {
	ctx := context.TODO()
	mongoconn := options.Client().ApplyURI(conf.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		return err
	}

	// Check the connection
	err = mongoclient.Ping(ctx, nil)
	if err != nil {
		return err
	}

	DB = Logger{mongoclient}
	log.Println("Connected to MongoDB!")
	return nil
}

func ConnectToRedis() error {
	redisclient := redis.NewClient(&redis.Options{
		Addr:     conf.RedisUri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if _, err := redisclient.Ping().Result(); err != nil {
		panic(err)
	}
	err := redisclient.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	log.Println("Redis client connected successfully...")
	return nil
}
