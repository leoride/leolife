package alerttype

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	DB_USER     = "tomg"
	DB_PASSWORD = "bnhr7r82"
	DB_NAME     = "leolife"
	DB_HOST     = "raspberrypi"
)

var tx *sql.Tx
var exampleId string
var exampleId2 string

func TestMain(m *testing.M) {
	//Before
	var err error
	r := 0

	db := initDb()
	defer db.Close()

	tx, err = db.Begin()

	if err == nil {
		r = m.Run()
	}

	//After
	_ = tx.Rollback()

	os.Exit(r)
}

func initDb() *sql.DB {
	var db *sql.DB
	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
	if db, err = sql.Open("postgres", dbinfo); err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Println("DB unreachable on application startup")
	}
	return db
}

func TestInsert(t *testing.T) {
	atfs := make([]AlertTypeField, 1)
	atfs[0] = AlertTypeField{"", "test", "test", "string", "test", true}
	at := AlertType{"", "test", "test", atfs}

	if err := insertAlertType(tx, &at); err != nil {
		t.Errorf(err.Error())
		return
	}

	exampleId = at.Id
	exampleId2 = at.Fields[0].Id

	if at.Id == "" {
		t.Errorf("Alert type Id is null")
		return
	}

	if at.Fields[0].Id == "" {
		t.Errorf("Alert type field Id is null")
		return
	}
}

func TestSelect(t *testing.T) {
	var at *AlertType
	var err error

	if at, err = getAlertType(tx, exampleId); err != nil {
		t.Errorf(err.Error())
		return
	}

	if at.Fields[0].Id != exampleId2 {
		t.Errorf("Alert type field Id is not correct")
		return
	}
}
