package views_html

func HTMLMessage(msgs ...string) string {
	var html string

	for _, msg := range msgs {
		html += "<li>" + msg + "</li>\n"
	}

	return html
}
