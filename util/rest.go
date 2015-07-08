package util

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
)

//Util: RestError
//RestError object used for rest errors
type RestError struct {
	Status           int    `json:"status"`
	Code             int    `json:"code"`
	Property         string `json:"property,omitempty"`
	Message          string `json:"message"`
	DeveloperMessage string `json:"developerMessage"`
	MoreInfo         string `json:"moreInfo,omitempty"`
}

func (e *RestError) Error() string {
	return e.DeveloperMessage
}

//Util: RestErrorWrapper
//Rest error wrapper for rest listeners
func RestErrorWrapper(f func(w http.ResponseWriter, r *http.Request) *RestError) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err *RestError

		if err = f(w, r); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(err.Status)
			errJson, _ := json.MarshalIndent(err, "", "    ")
			fmt.Fprint(w, string(errJson))
		}
	}
}

func RestAuthErrorWrapper(f func(w http.ResponseWriter, r *http.Request) *RestError) func(w http.ResponseWriter, r *http.Request) {
	return RestErrorWrapper(func(w http.ResponseWriter, r *http.Request) *RestError {
		var err error
		var restErr *RestError

		if tokenStr := r.Header.Get("auth-token"); tokenStr == "" {
			err = fmt.Errorf("No auth-token found in header")
		} else {
			var token *jwt.Token
			token, err = jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return []byte("leoride_leolife_rox"), nil })

			if (err == nil) && (token.Valid) {
				restErr = f(w, r)
			} else {
				err = fmt.Errorf("invalid token")
			}
		}

		if err != nil {
			restErr = new(RestError)
			restErr.Code = 1
			restErr.Status = 500
			restErr.DeveloperMessage = err.Error()
			restErr.Message = "An error has ocured while handling authorization for the API call"
		}

		return restErr
	})
}
