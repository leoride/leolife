package alerttype

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/leoride/leolife/util"
	"io/ioutil"
	"net/http"
)

//Listener: AlertType
//Listens to API calls for the AlertType resource
func ListenForAlertType(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/alertTypes", util.RestErrorWrapper(alertTypeCreateHandler(db))).Methods("POST")
	r.HandleFunc("/alertTypes", util.RestErrorWrapper(alertTypesHandler(db))).Methods("GET")
	r.HandleFunc("/alertType/{id}", util.RestErrorWrapper(alertTypeHandler(db))).Methods("GET")
	r.HandleFunc("/alertType/{id}", util.RestErrorWrapper(alertTypeUpdateHandler(db))).Methods("PUT")
	r.HandleFunc("/alertType/{id}", util.RestErrorWrapper(alertTypeDeleteHandler(db))).Methods("DELETE")
}

func alertTypeCreateHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) *util.RestError {
	return func(w http.ResponseWriter, r *http.Request) *util.RestError {
		var (
			restErr   *util.RestError
			err       error
			tx        *sql.Tx
			alertType AlertType
			jsonB     []byte
		)

		if tx, err = db.Begin(); err == nil {
			if jsonB, err = ioutil.ReadAll(r.Body); err == nil {
				if err = json.Unmarshal(jsonB, &alertType); err == nil {
					if err = insertAlertType(tx, &alertType); err == nil {
						if jsonB, err = json.MarshalIndent(alertType, "", "    "); err == nil {
							w.Header().Set("content-type", "application/json")
							w.WriteHeader(201)
							fmt.Fprint(w, string(jsonB))
						}
					}
				}
			}
		}

		if err != nil {
			restErr = new(util.RestError)
			restErr.Code = 0
			restErr.Status = 500
			restErr.DeveloperMessage = err.Error()
			restErr.Message = "An error has ocured while processing a request to insert an alert type"

			tx.Rollback()
		} else {
			err = tx.Commit()
		}

		return restErr
	}
}

func alertTypesHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) *util.RestError {
	return func(w http.ResponseWriter, r *http.Request) *util.RestError {

		var (
			restErr    *util.RestError
			err        error
			tx         *sql.Tx
			alertTypes []AlertType
			jsonB      []byte
		)

		if tx, err = db.Begin(); err == nil {
			if alertTypes, err = getAlertTypes(tx); err == nil {
				if jsonB, err = json.MarshalIndent(alertTypes, "", "    "); err == nil {
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(200)
					fmt.Fprint(w, string(jsonB))
				}
			}
		}

		if err != nil {
			restErr = new(util.RestError)
			restErr.Code = 0
			restErr.Status = 500
			restErr.DeveloperMessage = err.Error()
			restErr.Message = "An error has ocured while processing a request to retrieve alert types"

			tx.Rollback()
		} else {
			err = tx.Commit()
		}

		return restErr
	}
}

func alertTypeHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) *util.RestError {
	return func(w http.ResponseWriter, r *http.Request) *util.RestError {

		var (
			restErr   *util.RestError
			err       error
			tx        *sql.Tx
			alertType *AlertType
			jsonB     []byte
		)

		vars := mux.Vars(r)

		if tx, err = db.Begin(); err == nil {
			id := vars["id"]
			if alertType, err = getAlertType(tx, id); err == nil {
				if alertType != nil {
					if jsonB, err = json.MarshalIndent(alertType, "", "    "); err == nil {
						w.Header().Set("content-type", "application/json")
						w.WriteHeader(200)
						fmt.Fprint(w, string(jsonB))
					}
				} else {
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(404)
				}
			}
		}

		if err != nil {
			restErr = new(util.RestError)
			restErr.Code = 0
			restErr.Status = 500
			restErr.DeveloperMessage = err.Error()
			restErr.Message = "An error has ocured while processing a request to retrieve an alert type by id"

			tx.Rollback()
		} else {
			err = tx.Commit()
		}

		return restErr
	}
}

func alertTypeUpdateHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) *util.RestError {
	return func(w http.ResponseWriter, r *http.Request) *util.RestError {

		var (
			restErr   *util.RestError
			err       error
			tx        *sql.Tx
			alertType AlertType
			jsonB     []byte
		)

		vars := mux.Vars(r)

		if tx, err = db.Begin(); err == nil {
			id := vars["id"]
			if jsonB, err = ioutil.ReadAll(r.Body); err == nil {
				if err = json.Unmarshal(jsonB, &alertType); err == nil {
					alertType.Id = id
					if err = updateAlertType(tx, &alertType); err == nil {
						if jsonB, err = json.MarshalIndent(alertType, "", "    "); err == nil {
							w.Header().Set("content-type", "application/json")
							w.WriteHeader(200)
							fmt.Fprint(w, string(jsonB))
						}
					}
				}
			}
		}

		if err != nil {
			restErr = new(util.RestError)
			restErr.Code = 0
			restErr.Status = 500
			restErr.DeveloperMessage = err.Error()
			restErr.Message = "An error has ocured while processing a request to update an alert type by id"

			tx.Rollback()
		} else {
			err = tx.Commit()
		}

		return restErr
	}
}

func alertTypeDeleteHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) *util.RestError {
	return func(w http.ResponseWriter, r *http.Request) *util.RestError {

		var (
			restErr *util.RestError
			err     error
			tx      *sql.Tx
		)

		vars := mux.Vars(r)

		if tx, err = db.Begin(); err == nil {
			id := vars["id"]
			if err = deleteAlertType(tx, id); err == nil {
				w.Header().Set("content-type", "application/json")
				w.WriteHeader(200)
			}
		}

		if err != nil {
			restErr = new(util.RestError)
			restErr.Code = 0
			restErr.Status = 500
			restErr.DeveloperMessage = err.Error()
			restErr.Message = "An error has ocured while processing a request to delete an alert type by id"

			tx.Rollback()
		} else {
			err = tx.Commit()
		}

		return restErr
	}
}
