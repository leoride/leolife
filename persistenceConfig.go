package main

import (
	"fmt"
	"github.com/leoride/leolife/util/persistence"
)

var sessionFactory *persistence.SessionFactory

func SetupPersistenceConfig() {
	dbPort := 27017
	dbName := "leolife"

	sessionFactory = persistence.NewSessionFactory(fmt.Sprint(":", dbPort, "/", dbName))
}
