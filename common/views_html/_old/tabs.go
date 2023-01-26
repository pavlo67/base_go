package viewshtml

import (
	"strconv"
	"strings"
)

type Tab struct {
	HTMLTitle   string
	HTMLContent string
	Active      bool
}

func HTMLTabs(groupID string, tabs ...Tab) string {
	groupID += "_"

	var htmlContent string
	var htmlTitleItems []string
	for i, tab := range tabs {
		id := groupID + strconv.Itoa(i)

		htmlTitleItem := `<a href=# onclick="tabSelect('` + id + `','` + groupID + `', 'title_` + id + `','title_` + groupID + `')">` + tab.HTMLTitle + "</a>"

		if tab.Active {
			htmlTitleItems = append(htmlTitleItems, `<span id="title_`+id+`" style="font-weight:bold">`+htmlTitleItem+"</span>")
			htmlContent += `<div id="` + id + `" class="visible">` + tab.HTMLContent + "</div>\n"
		} else {
			htmlTitleItems = append(htmlTitleItems, `<span id="title_`+id+`" style="font-weight:normal">`+htmlTitleItem+"</span>")
			htmlContent += `<div id="` + id + `" class="hidden">` + tab.HTMLContent + "</div>\n"
		}
	}

	return strings.Join(htmlTitleItems, " &nbsp; ") + "\n<p>" + htmlContent
}
