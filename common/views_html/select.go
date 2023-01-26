package views_html

import (
	"html"
)

type SelectOptions [][2]string

const InlineFields = "inline"

func HTMLSelectEdit(selected string, selectStrings SelectOptions, attributes string) string {
	body := ""
	var option string
	for i := 0; i < len(selectStrings); i++ {
		body += "<option"
		if selectStrings[i][1] != "" {
			option = selectStrings[i][1]
			body += ` value="` + html.EscapeString(selectStrings[i][1]) + `"`
		} else {
			option = selectStrings[i][0]
		}
		if option == selected {
			body += " selected"
		}
		body += ">" + html.EscapeString(selectStrings[i][0]) + "</option>\n"
	}
	return `<select ` + attributes + `>` + body + "</select>\n"
}

func HTMLSelectView(selected string, selectStrings SelectOptions) string {
	for i := 0; i < len(selectStrings); i++ {
		option := selectStrings[i][0]
		if selectStrings[i][1] != "" {
			option = selectStrings[i][1]
		}
		if option == selected {
			return selectStrings[i][0]
		}
	}
	return ""
}

//// select string validation ---------------------------------------------------------------------------------
//
//var _ validator.Operator = &SelectStringValidator{}
//
//type SelectStringValidator struct {
//	data  string
//	label string
//	value string
//	errs  basis.Errors
//}
//
//func NewSelectString(data, label string, values SelectOptions, trim bool) SelectStringValidator {
//	if trim {
//		data = strings.TrimSpace(data)
//	}
//	value := ""
//	errs := basis.Errors{validator.BadValue}
//	for _, v := range values {
//		if v[1] == "" {
//			if v[0] == data {
//				value = data
//				errs = nil
//			}
//		} else {
//			if v[1] == data {
//				value = data
//				errs = nil
//			}
//		}
//	}
//
//	return SelectStringValidator{label: label, data: data, value: value, errs: errs}
//}
//
//func (v SelectStringValidator) Label() string {
//	return v.label
//}
//
//func (v SelectStringValidator) Errs() basis.Errors {
//	return v.errs
//}
//
//func (v SelectStringValidator) Value() string {
//	return v.value
//}
