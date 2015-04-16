package handler

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/leoride/leolife/entity"
	"github.com/leoride/leolife/util/persistence"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

type AccountHandler struct {
	sessionFactory *persistence.SessionFactory
}

func NewAccountHandler(sf *persistence.SessionFactory) AccountHandler {
	session := sf.GetNewSession()
	defer session.Close()

	entity.SetupAccountDBValidation(session.DB("leolife").C("account"))

	return AccountHandler{sf}
}

func (ah *AccountHandler) HandleCreate(w http.ResponseWriter, r *http.Request) error {
	var err error

	session := ah.sessionFactory.GetNewSession()
	defer session.Close()

	var body []byte
	body, err = ioutil.ReadAll(r.Body)

	if err == nil {
		var account entity.Account
		err, account = ah.unmarshalAccount(body)

		if err == nil {
			account.Id = bson.NewObjectId()
			err = account.Validate()

			if err == nil {

				//Hash password + set Validaded to false before insert!
				hasher := sha1.New()
				hasher.Write([]byte(account.Password))
				account.Password = fmt.Sprintf("%x", hasher.Sum(nil))
				account.Validaded = false

				err = session.DB("leolife").C("account").Insert(&account)

				//Remove password before to return
				account.Password = ""

				if err == nil {
					var accountJson string
					err, accountJson = ah.marshalAccount(account)

					if err == nil {
						w.WriteHeader(201)
						fmt.Fprint(w, accountJson)
					}
				}
			}
		}
	}

	return err
}

func (ah *AccountHandler) HandleAuthentication(w http.ResponseWriter, r *http.Request) error {
	var err error

	session := ah.sessionFactory.GetNewSession()
	defer session.Close()

	var body []byte
	body, err = ioutil.ReadAll(r.Body)

	if err == nil {
		var params map[string]string

		err = json.Unmarshal(body, &params)

		if err == nil {
			username := params["username"]
			password := params["password"]

			hasher := sha1.New()
			hasher.Write([]byte(password))
			password = fmt.Sprintf("%x", hasher.Sum(nil))

			var dbAccount entity.Account

			err = session.DB("leolife").C("account").Find(bson.M{"username": username}).One(&dbAccount)

			if err == nil {
				if (dbAccount.Username == username) && (dbAccount.Password == password) {

					if dbAccount.Validaded == false {
						err = fmt.Errorf("Inactive account")
					} else {
						token := jwt.New(jwt.SigningMethodHS256)
						token.Claims["username"] = username
						token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

						var tokenString string
						tokenString, err = token.SignedString([]byte("leoride_leolife_rox"))

						if err == nil {
							w.Header().Add("auth-token", tokenString)
						}
					}
				} else {
					err = fmt.Errorf("Wrong credentials")
				}

				if err == nil {
					w.WriteHeader(204)
				}
			}
		}
	}

	return err
}

func (ah *AccountHandler) marshalAccount(a entity.Account) (error, string) {
	var marshal string
	b, err := json.MarshalIndent(a, "", "    ")

	if err == nil {
		marshal = string(b)
	}

	return err, marshal
}

func (ah *AccountHandler) unmarshalAccount(b []byte) (error, entity.Account) {
	var a entity.Account
	err := json.Unmarshal(b, &a)
	return err, a
}
