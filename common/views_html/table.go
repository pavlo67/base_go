package views_html

func Table(labels []string, data [][]string) string {
	var html string

	first := true
	for _, row := range append([][]string{labels}, data...) {
		html += "<tr>"
		for _, cell := range row {
			if first {
				cell = "<b>" + cell + "</b>"
			}
			html += "<td>" + cell + "</td>\n"
		}
		html += "</tr>\n"
		first = false
	}

	return "<table border cellpadding=5 cellspacing=0>" + html + "</table>\n"
}
