package views_html

import (
	"strconv"
)

var num int

func HTMLOpenClose(title, content, imgPlus, imgMinus string, visible bool) string {
	num++

	id := strconv.Itoa(num)
	imageID := "plus_minus_img" + id
	contentID := "plus_minus_content" + id

	var htmlContent string

	if visible {
		htmlContent = `<a href=# onclick="openClose('` + imageID + `','` + contentID + `')">` +
			`<img id="` + imageID + `" src="` + imgMinus + `"></a> ` + title + "\n" +
			`<br><div id="` + contentID + `">` + content + `</div>`
	} else {
		htmlContent = `<a href=# onclick="openClose('` + imageID + `','` + contentID + `')">` +
			`<img id="` + imageID + `" src="` + imgPlus + `"></a> ` + title + "\n" +
			`<br><div id="` + contentID + `" style="visibility:hidden;position:absolute;">` + content + `</div>`

	}

	return htmlContent
}
