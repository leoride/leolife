package alerttype

import (
	"database/sql"
	"fmt"
	"github.com/leoride/leolife/util"
)

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

func updateAlertType(tx *sql.Tx, at *AlertType) error {
	var err error = nil

	sqlQ := "UPDATE alert_type set name = $1, description = $2 WHERE id = $3"

	if _, err = tx.Exec(sqlQ, at.Name, at.Description, at.Id); err == nil {

		for idx, atf := range at.Fields {
			if atf.Id == "" {
				if err = insertAlertTypeField(tx, &atf, at); err != nil {
					return err
				}
				at.Fields[idx] = atf
			} else {
				if err = updateAlertTypeField(tx, &atf, at); err != nil {
					return err
				}
				at.Fields[idx] = atf
			}
		}

		//watch for atf to delete
		if currAtfs, err := getAlertTypeFieldByAlertType(tx, at.Id); err == nil {
			for _, currAtf := range currAtfs {
				stillExists := false

				for _, atf := range at.Fields {
					if currAtf.Id == atf.Id {
						stillExists = true
					}
				}

				if !stillExists {
					err = deleteAlertTypeField(tx, currAtf.Id)
				}
			}
		}
	}

	return err
}

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

func updateAlertTypeField(tx *sql.Tx, atf *AlertTypeField, at *AlertType) error {
	var err error = nil

	sqlQ := "UPDATE alert_type_field set alert_type_id = $1, name = $2, label = $3, type = $4, \"default\" = $5, mandatory = $6 where id = $7"

	_, err = tx.Exec(sqlQ, at.Id, atf.Name, atf.Label, atf.Type, atf.Default, atf.Mandatory, atf.Id)

	return err
}

func deleteAlertTypeField(tx *sql.Tx, uuid string) error {
	var err error = nil

	sqlQ := "DELETE FROM alert_type_field WHERE id = $1"

	_, err = tx.Exec(sqlQ, uuid)

	return err
}
