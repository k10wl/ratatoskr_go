package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"tag"`
}

type Group struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	OriginalIndex int                `bson:"originalIndex"`
	Name          string             `bson:"groupName"`
	Tags          []Tag              `bson:"tags"`
}

type Analytics struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Tag   string             `bson:"tag"`
	Group string             `bson:"group"`
	Date  time.Time          `bson:"dateUsed"`
}
