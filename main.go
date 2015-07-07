package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/leoride/leolife/alert"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	DB_USER     = "tomg"
	DB_PASSWORD = "bnhr7r82"
	DB_NAME     = "leolife"
	DB_HOST     = "raspberrypi"
)

func main() {
	db := initDb()
	defer db.Close()
	r := mux.NewRouter()
	srV1 := r.PathPrefix("/api/v1/").Subrouter()

	//do stuff
	alert.ListenForAlertType(srV1, db)
	//done doing stuff

	http.Handle("/", r)
	log.Panic(http.ListenAndServe(":8080", nil))
}

//Initialize connection pool for database
func initDb() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("DB unreachable on application startup")
	}

	return db
}
