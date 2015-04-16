package rest

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)

func RestHandlerWrapper(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "Request received at:", r.RequestURI)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		var err error = f(w, r)

		if err != nil {
			log.Println("Error:", err)
			w.WriteHeader(500)

			restErr := RestError{"RestError", err.Error()}
			restErrB, _ := json.MarshalIndent(restErr, "", "    ")
			fmt.Fprint(w, string(restErrB))
		}
	}

}

func RestAuthHandlerWrapper(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {

	f2 := func(w http.ResponseWriter, r *http.Request) error {
		var err error

		tokenStr := r.Header.Get("auth-token")

		if tokenStr == "" {
			err = fmt.Errorf("No auth-token found in header")
		} else {
			var token *jwt.Token
			token, err = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return []byte("leoride_leolife_rox"), nil })

			if (err == nil) && (token.Valid) {
				err = f(w, r)

			} else {
				err = fmt.Errorf("invalid token")
			}
		}

		return err
	}

	return RestHandlerWrapper(f2)

}
