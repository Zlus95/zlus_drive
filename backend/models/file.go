package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`
	Name      string              `bson:"name"`
	Size      int64               `bson:"size"`
	OwnerID   primitive.ObjectID  `bson:"ownerId"`
	Path      string              `bson:"path"`
	MimeType  string              `bson:"mimeType"`
	CreatedAt time.Time           `bson:"createdAt"`
	IsFolder  bool                `bson:"isFolder"`
	Parent    *primitive.ObjectID `bson:"parent,omitempty" json:"parent,omitempty"`
	Children  []string            `bson:"сhildren" json:"сhildren"`
}
