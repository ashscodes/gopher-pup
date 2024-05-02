package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A Location represents a Person's location.
type Location struct {
	City    string `bson:"city,omitempty" json:"city,omitempty"`
	Country string `bson:"country,omitempty" json:"country,omitempty"`
}

// A Patch represents a Json Patch structure - http://jsonpatch.com/
type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// A Person represents a user.
type Person struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Firstname string             `bson:"firstname,omitempty" json:"firstname,omitempty"`
	Lastname  string             `bson:"lastname,omitempty" json:"lastname,omitempty"`
	Location  *Location          `json:"location,omitempty"`
}

type QueryFilter struct {
	Name       string
	ParentPath string
}

// People represents a slice of Person
type People []Person
