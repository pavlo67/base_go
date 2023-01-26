package views_html

import (
	"html"
)

func HTMLEditTable(fields []Field, data ValuesString, frontOps map[string]Operator, formID, url string) string { // ,
	// frontOps map[string]Operator, rView auth.ID, publicChanges bool
	//if data == nil {
	//	data = map[string]string{}
	//}
	//if values == nil {
	//	values = map[string]SelectString{}
	//}
	//
	//values["visibility"], data["visibility"] = dataView(user, rView, publicChanges)

	var editHTML, titleHTML, resHTML string
	for _, f := range fields {
		titleHTML, resHTML = FieldEdit(f, data, frontOps, formID)

		//if resHTML == "" && f.Params[NotEmptyKey] == true {
		//	continue
		//}

		if titleHTML != "" {
			titleHTML = "<small>" + titleHTML + ":</small> \n"
		}
		editHTML += `<tr><td style="width: 200px;">` + titleHTML + "</td><td>" + resHTML + "</td></tr>\n"
		//  id="div_` + html.EscapeString(formID+f.Key) + `"

	}

	return `<table width="100%"><form action="` + url + `" method="POST">` + "\n" + editHTML + "\n</form>\n</table>\n"
}

// var ReDigitsOnly = regexp.MustCompile(`^\d+$`)

func FieldEdit(f Field, data ValuesString, frontOps map[string]Operator, formID string) (string, string) {

	if f.Type == "view" || f.Type == "text" {
		return FieldView(f, data, frontOps, "")

	} else if f.Type == "button" || f.Type == "submit" {
		attributesAdd := Attributes{
			"id":    f.Key,
			"value": f.Format,

			// using generalNoFormID to add listeners on html pages
			//"data-form_id": formID,
			//"data-value":   data[f.Key],
		}
		return f.Label, `<input type="` + f.Type + `" ` + f.AttributesHTML(attributesAdd) + `/>`
	}

	var titleHTML = html.EscapeString(f.Label)
	var resHTML string

	attributesAdd := Attributes{
		"id":   formID + f.Key,
		"name": f.Key,
	}
	attributes := f.AttributesHTML(attributesAdd)
	if f.Type == "password" {
		resHTML = `<input style="width:100%" type="password" ` + f.AttributesHTML(attributesAdd) + ` />`
	} else if f.Type == "select" {
		selectStrings, _ := f.Params["select_options"].(SelectOptions)
		resHTML = HTMLSelectEdit(data[f.Key], selectStrings, attributes)
	} else if f.Type == "checkbox" {
		var checked string
		if data[f.Key] != "" {
			checked = " checked"
		}
		resHTML = `<input type="checkbox" ` + f.AttributesHTML(attributesAdd) + checked + `/>`
	} else if frontOp, ok := frontOps[f.Type]; ok {
		resHTML = frontOp.HTMLToEdit(f, data, formID)
	} else {
		var value = html.EscapeString(data[f.Key])
		if f.Type == "hidden" {
			resHTML = `<input type="hidden" ` + attributes + ` value="` + value + `" /> `
			titleHTML = ""
		} else if f.Type == "textarea" {
			rows := f.Format
			if rows == "" {
				rows = f.Params.StringDefault("rows", "1")
			}
			resHTML = `<textarea style="width:100%; padding:3px;" ` + attributes + ` rows=` + rows + `>` + value + `</textarea>`
		} else if f.Format == "number" {
			parameters := ""
			if step, _ := f.Params.String("step"); step != "" {
				parameters += ` step="` + step + `"`
			}
			if min, _ := f.Params.String("min"); min != "" {
				parameters += ` min="` + min + `"`
			}
			if max, _ := f.Params.String("max"); max != "" {
				parameters += ` max="` + max + `"`
			}
			resHTML = `<input type="number"` + parameters + attributes + ` value="` + value + `" />`
		} else if (f.Format == "date") || (f.Format == "time") || (f.Format == "datetime") || (f.Format == "email") || (f.Format == "url") || (f.Format == "color") {
			resHTML = `<input type="` + f.Type + `" ` + attributes + ` value="` + value + `" />`
		} else {
			resHTML = `<input type="` + f.Type + `"style="width:100%" ` + attributes + ` value="` + value + `" />`
		}
	}

	return titleHTML, resHTML
}
