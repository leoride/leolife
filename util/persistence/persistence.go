package persistence

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

//Session Factory structure
type SessionFactory struct {
	session mgo.Session
}

//Constructor for Session Factory
func NewSessionFactory(datasource string) *SessionFactory {

	mongoSession, err := mgo.Dial(datasource)

	if err != nil {
		log.Fatal("An error has occurred while trying to connect to the database: ", err)
		os.Exit(1)
	}

	sf := SessionFactory{*mongoSession}
	return &sf
}

func (sf *SessionFactory) GetNewSession() *mgo.Session {
	return sf.session.Clone()
}
