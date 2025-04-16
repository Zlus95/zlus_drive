package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password"`
	Name         string             `json:"name" bson:"name"`
	LastName     string             `json:"lastName" bson:"lastName"`
	StorageLimit int                `json:"storageLimit" bson:"storageLimit"`
	UsedStorage  int                `json:"usedStorage" bson:"usedStorage"`
}
