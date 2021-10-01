package views_html

import "html"

type SelectString [][2]string

func HTMLSelect(general string, values SelectString, selected string) string {
	body := ""
	var option string
	for i := 0; i < len(values); i++ {
		body += "<option"
		if values[i][1] != "" {
			option = values[i][1]
			body += ` value="` + html.EscapeString(values[i][1]) + `"`
		} else {
			option = values[i][0]
		}
		if option == selected {
			body += " selected"
		}
		body += ">" + html.EscapeString(values[i][0]) + "</option>\n"
	}
	return `<select ` + general + `>` + body + "</select>\n"
}

func HTMLSelectView(values SelectString, selected string) string {
	for i := 0; i < len(values); i++ {
		option := values[i][0]
		if values[i][1] != "" {
			option = values[i][1]
		}
		if option == selected {
			return values[i][0]
		}
	}
	return ""
}
