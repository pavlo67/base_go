package views_html

import (
	"html"

	"github.com/gomarkdown/markdown"

	"github.com/pavlo67/common/common"
)

func HTMLViewTable(fields []Field, data ValuesString, frontOps map[string]Operator, formats common.Map) string { // , frontOps map[string]Operator
	if data == nil {
		data = map[string]string{}
	}
	var viewHTML, titleHTML, resHTML string
	for _, f := range fields {
		if f.Type == "hidden" {
			continue
		}
		titleHTML, resHTML = FieldView(f, data, frontOps, formats.StringDefault(f.Key, "")) // , options, frontOps

		//if resHTML == "" && ((f.Params[NotEmptyKey] == true) || (titleHTML == "")) {
		//	continue
		//}

		if titleHTML != "" {
			titleHTML = "<small>" + titleHTML + ":</small>\n"
		}
		viewHTML += "<tr><td>\n" + titleHTML + "</td><td>&nbsp;</td><td>" + resHTML + "\n</td></tr>\n"
	}

	return `<table cellspacing=0">` + viewHTML + "</table>"
	// +`<input id=links_list type=hidden value="` + html.EscapeString(data["tags"]) + `">` + "\n"
}

// view - not editable data field
// text - text label only (no data field linked to!)

// https://pkg.go.dev/github.com/gomarkdown/markdown

const NotEmptyKey = "not_empty"
const NoEscapeKey = "no_escape"

func FieldView(f Field, data ValuesString, frontOps map[string]Operator, format string) (string, string) {

	var types = []string{"password", "button", "hidden"}
	for _, v := range types {
		if v == f.Type {
			return "", ""
		}
	}

	//if frontOp, ok := frontOps[f.ImporterInterfaceKey]; ok {
	//	params := map[string]string{
	//		// "format": f.Info,
	//		"class": class,
	//		"style": "width:100%",
	//	}
	//	return html.EscapeString(f.Label), frontOp.HTMLToView(f, data[f.Key], nil, params)
	//}

	//var class = field.Attributes["class"]
	//if class != "" {
	//	class = ` class="` + html.EscapeString(class) + `"`
	//}

	//if field.Options.StringDefault("format", "") == "datetime" {
	//	resHTML = html.EscapeString(data[field.Key])

	var resHTML string

	if f.Type == "select" {
		selectOptions, _ := f.Params["select_options"].(SelectOptions)
		resHTML = HTMLSelectView(data[f.Key], selectOptions)
	} else if f.Format == "url" {
		var url = html.EscapeString(data[f.Key])
		resHTML = `<a href="` + url + `" target=_blank>` + url + `</a>`

	} else if f.Type == "text" {
		resHTML = html.EscapeString(f.Format)

	} else if f.Type == "checkbox" {
		if data[f.Key] == "on" {
			resHTML = "так"
		} else if f.Params[NotEmptyKey] != true {
			resHTML = "ні"
		}
	} else if f.Params[NotEmptyKey] == true && data[f.Key] == "0" {
		// shows nothing
	} else if f.Params[NoEscapeKey] == true {
		resHTML = data[f.Key]

	} else if format == "md" {
		md := markdown.NormalizeNewlines([]byte(data[f.Key]))
		resHTML = string(markdown.ToHTML(md, nil, nil))
		// log.Printf("%s --> \n%s", md, resHTML)

	} else {

		resHTML = html.EscapeString(data[f.Key])

	}
	return html.EscapeString(f.Label), resHTML
}
