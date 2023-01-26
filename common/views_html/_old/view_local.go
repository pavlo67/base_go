// +build local

package notebook_comp

import (
	"net/http"

	"github.com/pavlo67/punctum/basis/viewshtml"
	"github.com/pavlo67/punctum/confidenter/auth"
	"github.com/pavlo67/punctum/confidenter/controller"
)

var tplLeftNoUser, tplLeft map[string]string

func initHTML() {
	htmlSearchMy := `<div class="ul utd">` +
		`<input style="width:171px;" id="to_search_my"">` +
		`<input type="button" id="` + listeners["searchMy"].ID + `" value="знайти">` +
		"</div>\n\n"

	htmlLeft := `<div class="ut">` + viewshtml.My + ` Мої сторінки</div>` +
		`<div class="ul">` + viewshtml.My +
		`<a href="` + itemsEndpoints["tagsMy"].ServerPath + `">` + viewshtml.Tags + `</a> ` +
		`<a href="` + itemsEndpoints["itemsMy"].ServerPath + `">записи</a>` + "</div>\n" +

		`<div class="ul">` + viewshtml.My + ` <a href="` + itemsEndpoints["blank"].Path(GenusDefault, string(controller.Owner)) + `">новий запис</a></div>` +
		htmlSearchMy

	htmlLeftNoUser := `<div class="ut gray">` + viewshtml.No + ` Мої сторінки</div>` +
		`<div class="ul gray">` + viewshtml.No + viewshtml.Tags + "записи</a></div>\n" +
		`<div class="ul gray">` + viewshtml.No + ` новий запис</a></div>`

	tplLeftNoUser = map[string]string{
		"left.comp": htmlLeftNoUser,
		"front":     htmlFront,
	}

	tplLeft = map[string]string{
		"left.comp": htmlLeft,
		"front":     htmlFront,
	}

}

func notebookTemplator(r *http.Request, user *auth.User) map[string]string {
	if user == nil || user.ID == "" {
		return tplLeftNoUser
	}

	return tplLeft
}
