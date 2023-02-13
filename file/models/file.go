package models

type File struct {
	Name    string `json:"name" bson:"name,omitempty"`
	Ext     string `json:"ext" bson:"ext,omitempty"`
	OwnerId string `json:"owner_id" bson:"owner_id,omitempty"`
}
