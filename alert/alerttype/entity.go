package alerttype

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
