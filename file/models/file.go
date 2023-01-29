package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type File struct {
	Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name,omitempty"`
	Type    string             `json:"type" bson:"type,omitempty"`
	Url     string             `json:"url" bson:"url,omitempty"`
	OwnerId string             `json:"owner_id" bson:"owner_id,omitempty"`
}
