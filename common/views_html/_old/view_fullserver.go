// +build fullserver linux windows

package notebook_comp

import (
	"net/http"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/controller"
	"github.com/pavlo67/punctum/confidenter/rights"

	"github.com/pavlo67/punctum/basis/viewshtml"
)

var tplLeftNoUser map[string]string
var htmlLeftTop, htmlLeftBottom string

func initHTML() {
	htmlPublic := `<div class="ul">` + viewshtml.Public +
		`<a href="` + itemsEndpoints["tags"].Path(string(basis.Anyone)) + `">` + viewshtml.Tags + `</a> ` +
		`<a href="` + itemsEndpoints["itemsTop"].Path() + `">записи для загалу</a>` + "</div>\n"

	htmlSearch := `<div class="ul utd">` +
		`<input style="width:171px;" id="to_search">` +
		`<input type="button" id="` + listeners["search"].ID + `" value="знайти">` +
		"</div>\n\n"

	htmlSearchMy := `<div class="ul utd">` +
		`<input style="width:171px;" id="to_search_my"">` +
		`<input type="button" id="` + listeners["searchMy"].ID + `" value="знайти">` +
		"</div>\n\n"

	htmlLeftTop = `<div class="ut">` + viewshtml.Public + `<a href="/">Ой, мамо! Де я???</a></div>` +
		`<div class="ut">` + viewshtml.My + ` Мої сторінки</div>` +
		`<div class="ul">` + viewshtml.My +
		`<a href="` + itemsEndpoints["tagsMy"].ServerPath + `">` + viewshtml.Tags + `</a> ` +
		`<a href="` + itemsEndpoints["itemsMy"].ServerPath + `">записи</a>` + "</div>\n" +
		`<div class="ul">` + viewshtml.My + ` <a href="` + itemsEndpoints["blank"].Path(GenusDefault, string(controller.Owner)) + `">новий запис</a></div>` +
		htmlSearchMy +
		`<div class="ut">` + viewshtml.Public + ` Публічні сторінки` +
		htmlPublic

	htmlLeftBottom = "\n</div>" +
		`<div class="ul">` + viewshtml.Public + ` <a href="` + itemsEndpoints["blank"].Path(GenusDefault, string(basis.Anyone)) + `">новий запис</a></div>` +
		"\n</div>" +
		htmlSearch

	htmlLeftNoUser := `<div class="ut gray">` + viewshtml.No + ` Мої сторінки</div>` +
		`<div class="ul gray">` + viewshtml.No + viewshtml.Tags + "записи</a></div>\n" +
		`<div class="ul gray">` + viewshtml.No + ` новий запис</a></div>` +
		`<div class="ut">` + viewshtml.Public + ` Публічні сторінки` +
		htmlPublic +
		htmlSearch

	tplLeftNoUser = map[string]string{
		"left.comp": htmlLeftNoUser,
		"front":     htmlFront,
	}
}

func notebookTemplator(r *http.Request, user *auth.User) map[string]string {
	if user == nil || user.ID == "" {
		return tplLeftNoUser
	}

	htmlLeft := htmlLeftTop

	for _, gr := range user.Accesses {
		if gr.Right == rights.Member {
			htmlLeft += `<div class="ul">` + viewshtml.Group +
				`<a href="` + itemsEndpoints["tags"].Path(string(gr.IS)) + `">` + viewshtml.Tags + `</a> ` +
				`<a href="` + itemsEndpoints["items.comp"].Path(string(gr.IS)) + `">` + gr.Label + `</a></div>`
		}
	}

	tplLeft := map[string]string{
		"left.comp": htmlLeft + htmlLeftBottom,
		"front":     htmlFront,
	}

	return tplLeft
}
