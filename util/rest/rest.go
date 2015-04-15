package rest

import (
	"encoding/json"
	"fmt"
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
