package alert

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/leoride/leolife/util"
	"io/ioutil"
	"net/http"
)

//Entity: AlertType
//It defines an alert type.
type AlertType struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Fields      []AlertTypeField `json:"fields"`
}

//Entity: AlertTypeField
//It defines the custom fields of an alert type.
type AlertTypeField struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Label     string `json:"label"`
	Type      string `json:"type"`
	Default   string `json:"default"`
	Mandatory bool   `json:"mandatory"`
}

//Repository: AlertType
//Finds all alert types
func getAlertTypes(tx *sql.Tx) ([]AlertType, error) {
	var alertTypes []AlertType = make([]AlertType, 0)
	var err error = nil

	var rows *sql.Rows
	sqlQ := "SELECT id, name, description FROM alert_type"

	if rows, err = tx.Query(sqlQ); err == nil {
		defer rows.Close()

		for rows.Next() {
			alertType := AlertType{}

			if err = rows.Scan(&alertType.Id, &alertType.Name, &alertType.Description); err == nil {
				alertTypes = append(alertTypes, alertType)
			}
		}
	}

	if err == nil {
		err = rows.Err()
	}
	rows.Close()

	for idx, alertType := range alertTypes {
		var alertTypeFields []AlertTypeField

		if alertTypeFields, err = getAlertTypeFieldByAlertType(tx, alertType.Id); err == nil {
			alertTypes[idx].Fields = alertTypeFields
		}
	}

	return alertTypes, err
}

//Repository: AlertType
//Finds an alert type by id
func getAlertType(tx *sql.Tx, id string) (*AlertType, error) {
	var alertType *AlertType
	var err error = nil

	var row *sql.Row
	var alertTypeFetched AlertType
	sqlQ := "SELECT id, name, description FROM alert_type where id = $1"

	row = tx.QueryRow(sqlQ, id)
	if err = row.Scan(&alertTypeFetched.Id, &alertTypeFetched.Name, &alertTypeFetched.Description); err == nil {
		if &alertTypeFetched.Id != nil {
			alertType = &alertTypeFetched
			var alertTypeFields []AlertTypeField

			if alertTypeFields, err = getAlertTypeFieldByAlertType(tx, alertType.Id); err == nil {
				alertType.Fields = alertTypeFields
			}
		}
	} else {
		if err.Error() == "sql: no rows in result set" {
			err = nil
		}
	}

	return alertType, err
}

//Repository: AlertType
//Inserts a new alert type
func insertAlertType(tx *sql.Tx, at *AlertType) error {
	var err error = nil
	var uuid string

	uuid, err = util.GenerateUuid()
	sqlQ := "INSERT INTO alert_type(id, name, description) VALUES($1, $2, $3)"

	if _, err = tx.Exec(sqlQ, uuid, at.Name, at.Description); err == nil {

		at.Id = uuid

		for idx, atf := range at.Fields {
			if err = insertAlertTypeField(tx, &atf, at); err != nil {
				return err
			}
			at.Fields[idx] = atf
		}
	}

	return err
}

//Repository: AlertType
//Deletes an alert type
func deleteAlertType(tx *sql.Tx, uuid string) error {
	var err error = nil
	var at *AlertType

	if at, err = getAlertType(tx, uuid); err == nil {
        if at != nil {
            for _, atf := range at.Fields {
                if err = deleteAlertTypeField(tx, atf.Id); err != nil {
                    break
                }
            }
        } else {
            err = fmt.Errorf("Resource not found")
        }

		if err == nil {
			sqlQ := "DELETE FROM alert_type WHERE id = $1"
			_, err = tx.Exec(sqlQ, uuid)
		}
	}

	return err
}

//Repository: AlertTypeField
//Finds all alert type fields linked to a specific alert type
func getAlertTypeFieldByAlertType(tx *sql.Tx, id string) ([]AlertTypeField, error) {
	var alertTypeFields []AlertTypeField = make([]AlertTypeField, 0)
	var err error = nil

	var rows *sql.Rows
	sqlQ := "SELECT id, name, label, type, \"default\", mandatory FROM alert_type_field where alert_type_id = $1"

	if rows, err = tx.Query(sqlQ, id); err == nil {
		defer rows.Close()

		for rows.Next() {
			alertTypeField := AlertTypeField{}
			var nilD sql.NullString

			if err = rows.Scan(&alertTypeField.Id, &alertTypeField.Name, &alertTypeField.Label, &alertTypeField.Type, &nilD, &alertTypeField.Mandatory); err == nil {
				if nilD.Valid {
					alertTypeField.Default = nilD.String
				}
				alertTypeFields = append(alertTypeFields, alertTypeField)
			}
		}

		if err == nil {
			err = rows.Err()
		}
	}

	return alertTypeFields, err
}

//Repository: AlertTypeField
//Inserts a new alert type field
func insertAlertTypeField(tx *sql.Tx, atf *AlertTypeField, at *AlertType) error {
	var err error = nil
	var uuid string

	uuid, err = util.GenerateUuid()
	sqlQ := "INSERT INTO alert_type_field(id, alert_type_id, name, label, type, \"default\", mandatory) VALUES($1, $2, $3, $4, $5, $6, $7)"

	if _, err = tx.Exec(sqlQ, uuid, at.Id, atf.Name, atf.Label, atf.Type, atf.Default, atf.Mandatory); err == nil {
		atf.Id = uuid
	}

	return err
}

//Repository: AlertTypeField
//Deletes an alert type field
func deleteAlertTypeField(tx *sql.Tx, uuid string) error {
	var err error = nil

	sqlQ := "DELETE FROM alert_type_field WHERE id = $1"

	_, err = tx.Exec(sqlQ, uuid)

	return err
}

//Listener: AlertType
//Listens to API calls for AlertType resource
func ListenForAlertType(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/alertTypes", util.RestErrorWrapper(alertTypeCreateHandler(db))).Methods("POST")
	r.HandleFunc("/alertTypes", util.RestErrorWrapper(alertTypesHandler(db))).Methods("GET")
	r.HandleFunc("/alertType/{id}", util.RestErrorWrapper(alertTypeHandler(db))).Methods("GET")
	r.HandleFunc("/alertType/{id}", util.RestErrorWrapper(alertTypeDeleteHandler(db))).Methods("DELETE")
}

//Listener: AlertType Create Handler
//Handles requests for inserting an AlertType
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

//Listener: AlertTypes Handler
//Handles requests for getting AlertTypes
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

//Listener: AlertType Handler
//Handles requests for getting an AlertType by Id
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

//Listener: AlertType Handler
//Handles requests for deleting an AlertType by Id
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
