package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"tag"`
}

type Group struct {
	ID            primitive.ObjectID `bson:"_id"`
	OriginalIndex int                `bson:"originalIndex"`
	Name          string             `bson:"groupName"`
	Tags          []Tag              `bson:"tags"`
}

type Analytics struct {
	ID    primitive.ObjectID `bson:"_id"`
	Tag   string             `bson:"tag"`
	Group string             `bson:"group"`
	Date  primitive.DateTime `bson:"dateUsed"`
}
