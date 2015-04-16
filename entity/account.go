package entity

import (
	"fmt"
	"github.com/leoride/leolife/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Account is an entity
type Account struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Username  string        `json:"username" bson:"username"`
	Password  string        `json:"password,omitempty" bson:"password,omitempty"`
	Email     string        `json:"email" bson:"email"`
	Validaded bool          `json:"validaded" bson:"validaded"`
}

func SetupAccountDBValidation(c *mgo.Collection) error {
	var err error

	index1 := mgo.Index{Key: []string{"username"}, Unique: true, Sparse: true}
	index2 := mgo.Index{Key: []string{"email"}, Unique: true, Sparse: true}

	err = c.EnsureIndex(index1)
	if err != nil {
		return err
	}
	err = c.EnsureIndex(index2)
	if err != nil {
		return err
	}

	return nil
}

func (a Account) Validate() error {
	var err error

	if a.Id == bson.ObjectId("") {
		err = fmt.Errorf("ObjectId must exist")
	} else if !util.IsEmailValid(a.Email) {
		err = fmt.Errorf("Email is not valid")
	} else if a.Password == "" {
		err = fmt.Errorf("Password must exist")
	} else if a.Username == "" {
		err = fmt.Errorf("Username must exist")
	}

	return err
}
