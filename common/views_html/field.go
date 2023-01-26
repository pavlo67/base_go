package views_html

import (
	"fmt"
	"html"

	"github.com/pavlo67/common/common"
)

type ValuesString map[string]string
type Attributes map[string]string

type Field struct {
	Key    string     `bson:"key"              json:"key"`
	Label  string     `bson:"label,omitempty"  json:"label,omitempty"`
	Type   string     `bson:"type,omitempty"   json:"type,omitempty"`
	Format string     `bson:"format,omitempty" json:"format,omitempty"`
	Params common.Map `bson:"params,omitempty" json:"params,omitempty"`
	// in particular
	//    .Params["attributes_html"]
	//    .Params["select_options"]
	//    .Params["rows"]                                 - for textarea
	//    .Params["min"], .Params["max"], .Params["step"] - for number
}

type Fields []Field

func (fields Fields) SetFormat(key, format string) error {
	for i, f := range fields {
		if f.Key == key {
			fields[i].Format = format
			return nil
		}
	}
	return fmt.Errorf("field (%s) not found to set format", key)
}

func (f *Field) AttributesHTML(attributesAdd Attributes) string {
	if f == nil {
		return ""
	}

	attributes, _ := f.Params["attributes_html"].(Attributes)
	if attributes == nil {
		attributes = attributesAdd
	} else {
		for k, v := range attributesAdd {
			attributes[k] = v
		}
	}

	var attributesHTML string
	for k, v := range attributes {
		attributesHTML += " " + html.EscapeString(k) + `="` + html.EscapeString(v) + `"`
	}

	return attributesHTML
}

//func AttributeHTML(key, value string) string {
//	return key + `="` + html.EscapeString(value) + `"`
//}
