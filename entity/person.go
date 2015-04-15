package entity

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	MALE   Gender = "male"
	FEMALE Gender = "female"

	BOYFRIEND  Relationship = "boyfriend"
	GIRLFRIEND Relationship = "girlfriend"
	DAD        Relationship = "father"
	MOM        Relationship = "mother"
	SIS        Relationship = "sister"
	BRO        Relationship = "brother"
)

//Enumeration type for gender
type Gender string

//Enumeration type for relationship
type Relationship string

//Person is an entity
type Person struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	FirstName string        `json:"firstName" bson:"firstName"`
	LastName  string        `json:"lastName" bson:"lastName"`
	BirthDate string        `json:"birthDate" bson:"birthDate"`
	Gender    Gender        `json:"gender" bson:"gender"`
	Contact   Contact       `json:"contact" bson:"contact"`
	WatchList WatchList     `json:"watchList" bson:"watchList"`
}

//Person's contact information is an entity
type Contact struct {
	Email   string `json:"email" bson:"email"`
	City    string `json:"city" bson:"city"`
	Country string `json:"country" bson:"country"`
}

//Person's relationships
type WatchList struct {
	Relationship Relationship `json:"relationship" bson:"relationship"`
	PersonId     *mgo.DBRef   `json:"personId" bson:"personId"`
}
