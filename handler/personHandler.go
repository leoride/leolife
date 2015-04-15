package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/leoride/leolife/entity"
	"github.com/leoride/leolife/util/persistence"
	"github.com/leoride/leolife/util/rest"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

type PersonHandler struct {
	sessionFactory *persistence.SessionFactory
}

func NewPersonHandler(sf *persistence.SessionFactory) PersonHandler {
	return PersonHandler{sf}
}

func (ph *PersonHandler) HandleGetOne(w http.ResponseWriter, r *http.Request) error {
	var person entity.Person
	var err error
	id := mux.Vars(r)["id"]

	session := ph.sessionFactory.GetNewSession()
	defer session.Close()

	validId := bson.IsObjectIdHex(id)

	if !validId {
		err = fmt.Errorf("Invalid ObjectId")
	}

	if err == nil {
		err = session.DB("leolife").C("person").Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&person)

		if err == nil {
			var personJson string
			err, personJson = ph.marshalPerson(person)

			if err == nil {
				w.WriteHeader(200)
				fmt.Fprint(w, personJson)
			}
		}
	}

	return err
}

func (ph *PersonHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) error {
	var page rest.Page
	var err error

	session := ph.sessionFactory.GetNewSession()
	defer session.Close()

	var count int
	count, err = session.DB("leolife").C("person").Count()

	if err == nil {
		var persons []entity.Person
		err = session.DB("leolife").C("person").Find(bson.M{}).All(&persons)

		if err == nil {
			personsI := make([]interface{}, len(persons))

			for index, value := range persons {
				personsI[index] = value
			}

			page.TotalElements = count
			page.PageNumber = 1
			page.PageSize = count
			page.Results = personsI
			page.TotalPages = 1

			var personsJson string
			err, personsJson = rest.MarshalPage(page)

			if err == nil {
				w.WriteHeader(200)
				fmt.Fprint(w, personsJson)
			}
		}
	}

	return err
}

func (ph *PersonHandler) HandleUpdateOne(w http.ResponseWriter, r *http.Request) error {
	var err error
	id := mux.Vars(r)["id"]

	session := ph.sessionFactory.GetNewSession()
	defer session.Close()

	validId := bson.IsObjectIdHex(id)

	if !validId {
		err = fmt.Errorf("Invalid ObjectId")
	}

	if err == nil {
		var body []byte
		body, err = ioutil.ReadAll(r.Body)

		if err == nil {
			var personU entity.Person
			err, personU = ph.unmarshalPerson(body)

			if err == nil {
				if personU.Id == bson.ObjectId("") {
					personU.Id = bson.ObjectIdHex(id)
				}

				if personU.Id != bson.ObjectIdHex(id) {
					err = fmt.Errorf("Invalid ObjectId in JSON")
				}

				if err == nil {
					err = session.DB("leolife").C("person").UpdateId(bson.ObjectIdHex(id), &personU)

					if err == nil {
						var personUJson string
						err, personUJson = ph.marshalPerson(personU)

						if err == nil {
							w.WriteHeader(200)
							fmt.Fprint(w, personUJson)
						}

					}
				}
			}
		}
	}

	return err
}

func (ph *PersonHandler) HandleDeleteOne(w http.ResponseWriter, r *http.Request) error {
	var err error
	id := mux.Vars(r)["id"]

	session := ph.sessionFactory.GetNewSession()
	defer session.Close()

	validId := bson.IsObjectIdHex(id)

	if !validId {
		err = fmt.Errorf("Invalid ObjectId")
	}

	if err == nil {
		err = session.DB("leolife").C("person").RemoveId(bson.ObjectIdHex(id))

		if err == nil {
			w.WriteHeader(204)
		}
	}

	return err
}

func (ph *PersonHandler) HandleCreate(w http.ResponseWriter, r *http.Request) error {
	var err error

	session := ph.sessionFactory.GetNewSession()
	defer session.Close()

	var body []byte
	body, err = ioutil.ReadAll(r.Body)

	if err == nil {
		var person entity.Person
		err, person = ph.unmarshalPerson(body)

		if err == nil {
			person.Id = bson.NewObjectId()
			err = session.DB("leolife").C("person").Insert(&person)

			if err == nil {
				var personJson string
				err, personJson = ph.marshalPerson(person)

				if err == nil {
					w.WriteHeader(201)
					fmt.Fprint(w, personJson)
				}
			}
		}
	}

	return err
}

func (ph *PersonHandler) marshalPerson(p entity.Person) (error, string) {
	var marshal string
	b, err := json.MarshalIndent(p, "", "    ")

	if err == nil {
		marshal = string(b)
	}

	return err, marshal
}

func (ph *PersonHandler) unmarshalPerson(b []byte) (error, entity.Person) {
	var p entity.Person
	err := json.Unmarshal(b, &p)
	return err, p
}
