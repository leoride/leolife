package util

import (
	"encoding/json"
	"fmt"
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
