package rest

import (
	"encoding/json"
	"fmt"
	"time"
)

type RestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PageOptions struct {
	PageNumber int
	PageSize   int
	OrderBy    string
	OrderDir   string
}

type Page struct {
	Results       []interface{} `json:"results"`
	PageNumber    int           `json:"pageNumber"`
	PageSize      int           `json:"pageSize"`
	TotalElements int           `json:"totalElements"`
	TotalPages    int           `json:"totalPages"`
}

func MarshalPage(p Page) (error, string) {
	var marshal string
	b, err := json.MarshalIndent(p, "", "    ")

	if err == nil {
		marshal = string(b)
	}

	return err, marshal
}

type SimpleDate time.Time

func (t *SimpleDate) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(*t).Format("01/02/2006"))
	return []byte(stamp), nil
}

func (jsT *SimpleDate) UnmarshalJSON(b []byte) error {
	var err error = nil
	var t time.Time

	t, err = time.Parse("01/02/2006", "08/15/1989") //string(b))

	if err == nil {
		*jsT = SimpleDate(t)
	}

	return err
}
