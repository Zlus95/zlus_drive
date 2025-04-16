package config

import "go.mongodb.org/mongo-driver/mongo"

var (
	UserCollection *mongo.Collection
)

func InitCollections(db *mongo.Database) {
	UserCollection = db.Collection("users")
}
