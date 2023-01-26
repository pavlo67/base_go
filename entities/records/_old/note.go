package _old

const maxBriefLength = 1000
const GenusKey = "note"

//var buttonAttributes = viewshtml.Attributes{
//	"class": "ut button",
//}
//
//var noteCreateFields = append(dataFields, viewshtml.Field{"save", "Зберегти запис", "button", "", nil, buttonAttributes})
//var noteUpdateFields = append(dataFields, viewshtml.Field{"save", "Зберегти зміни", "button", "", nil, buttonAttributes})
//
//func Genus(ctrlOp groups.Operator, endpoints map[string]config.Endpoint, listeners map[string]config.Listener, pxPreview int) *genera.Genus {
//	return &genera.Genus{
//		ID:         GenusKey,
//		Key:        GenusKey,
//		Name:       "нотатка",
//		NamePlural: "нотатки",
//		Translator: &noteTranslator{
//			ctrlOp:         ctrlOp,
//			endpoints:      endpoints,
//			listeners:      listeners,
//			pxPreview:      pxPreview,
//			maxBriefLength: maxBriefLength,
//		},
//	}
//}
//
//const txtNoNote = "Нема такої нотатки"
//
//// genera.Translator -----------------------------------------------------------------------------------------
//
//var _ genera.Translator = &noteTranslator{}
//
//type noteTranslator struct {
//	ctrlOp         groups.Operator
//	endpoints      map[string]config.Endpoint
//	listeners      map[string]config.Listener
//	pxPreview      int
//	maxBriefLength int
//}
//
//func (gt *noteTranslator) DataFromObject(userIS auth.ID, o *notes.Item, dataDefault genera.DataRaw) (genera.DataRaw, genera.EditForm, basis.Errors) {
//
//	editForm := genera.EditForm{
//		Fields: noteUpdateFields,
//	}
//
//	data := genera.DataRaw{}
//	for k, v := range dataDefault {
//		if vStr, ok := v.([]string); ok && len(vStr) > 0 {
//			data[k] = vStr[0]
//		} else if vStr, ok := v.(string); ok {
//			data[k] = vStr
//		}
//	}
//	data["genus"] = GenusKey
//
//	if o == nil {
//		return data, editForm, nil
//	}
//
//	data["id"] = o.ID
//	data["title"] = o.Name
//	data["genus"] = o.Genus
//	data["content"] = o.Content
//
//	linksList, err := json.Marshal(o.Links)
//	if err != nil {
//		return nil, editForm, basis.Errors{errors.Wrapf(err, "can't marshal object.tags.comp: %#v for object.id: %s", o.Links, o.ID)}
//	}
//	data["links"] = string(linksList)
//
//	tags := ""
//	filesList := []notes.Item{}
//	for _, l := range o.Links {
//		switch l.Type {
//
//		case links.TypeTag:
//			tags += l.Name + "; "
//
//		case files.LinkType:
//			filesList = append(filesList, l)
//		}
//	}
//	data["tags"] = tags
//	if len(filesList) > 0 {
//		files, err := json.Marshal(filesList)
//		if err != nil {
//			log.Println(err)
//		}
//		data["fileslocal"] = string(files)
//	}
//
//	if o.UpdatedAt != nil {
//		data["updated_at"] = o.UpdatedAt.Format("02.01.2006 15:04:05")
//	}
//
//	return data, editForm, nil
//}
//
//func (gt *noteTranslator) ObjectFromData(userIS auth.ID, oOld *notes.Item, dataRaw genera.DataRaw, linksList []notes.Item) (o *notes.Item, index interface{}, errs basis.Errors) {
//	if userIS == nil {
//		return nil, nil, basis.Errors{basis.ErrBadIdentity}
//	}
//
//	data := dataRaw.StringsMap()
//
//	var visibility string
//	var rView auth.ID
//	var managers = rights.Managers{}
//
//	visibility_ := data["visibility"]
//
//	if visibility_ == string(basis.Anyone) {
//		rView = basis.Anyone
//		visibility = things_old.Public
//	} else if visibility_ == string(controller.Owner) {
//		rView = userIS.String()
//		visibility = things_old.Private
//	} else if items_comp.ReCommonEdit.MatchString(visibility_) {
//		rView = auth.ID(items_comp.ReCommonEdit.ReplaceAllString(visibility_, ""))
//		managers[rights.Change] = rView
//		visibility = things_old.InGroup
//	} else {
//		rView = auth.ID(visibility_)
//		managers[rights.Change] = userIS.String()
//		visibility = things_old.InGroup
//	}
//
//	o = &notes.Item{
//		ID:         data["id"],
//		Genus:      data["genus"],
//		Author:     data["author"],
//		Name:       data["title"],
//		Content:    data["content"],
//		Links:      linksList,
//		Visibility: visibility,
//		RView:      rView,
//		Managers:   managers,
//	}
//
//	runes := []rune(o.Content)
//	if len(runes) > gt.maxBriefLength {
//		o.Brief = string(runes[0:gt.maxBriefLength]) + "..."
//	} else {
//		o.Brief = o.Content
//	}
//
//	arr := strings.Split(o.Name, ":")
//	if len(arr) > 1 {
//		o.Links = append(o.Links, notes.Item{Type: "author", Name: arr[0]})
//	}
//
//	notes.AddTags(userIS, o, data["tags"])
//
//	return o, nil, nil
//}
//
//func (gt *noteTranslator) View(userIS auth.ID, o *notes.Item, linkedObjects []notes.Item, context *genera.Context) genera.DataView {
//	if o == nil {
//		return map[string]string{
//			"caput":   "Перегляд",
//			"titulus": "Перегляд",
//			"corpus":  txtNoNote,
//		}
//	}
//
//	var i, htmlIndex, htmlContent, htmlShareTags, htmlShare, linksTitle, htmlLinked string
//	canChange := groups.OneOf(userIS, gt.ctrlOp, o.ROwner, o.Managers[rights.Change])
//	canDelete := groups.OneOf(userIS, gt.ctrlOp, o.ROwner, o.Managers[rights.Delete])
//
//	if o.Content != "" || canChange || canDelete {
//		htmlIndex = "<tr><td>" + items_comp.HTMLAuthor(userIS, o) + "</td></tr>\n"
//	}
//
//	if o.Content != "" {
//		mtext := mt.Read(o.Content)
//		i, htmlContent = mtext.HTML(0, 0)
//		if i != "" {
//			htmlContent = i + "<p>&nbsp;<p>" + htmlContent
//		}
//
//		htmlShareTags = items_comp.HTMLFBTags(o)
//		htmlShare = "<br>&nbsp;<table cellpadding=0 cellspacing=0><tr><td valign=top>" + items_comp.HTMLFBShare(joiner.SystemDomain()+gt.endpoints["view"].Path(o.ID), o.RView == basis.Anyone) +
//			"<td>&nbsp;<td>" + items_comp.HTMLTwitterShare(o.Name+" "+o.Brief, joiner.SystemDomain()+gt.endpoints["view"].Path(o.ID)) +
//			"</tr></table>" +
//			items_comp.HTMLFBSDK
//
//		linksTitle = "<p><b>Повʼязані записи, підрозділи</b>\n<p>"
//	}
//
//	if canChange {
//		htmlIndex += `<tr><td>- <a href="` + gt.endpoints["edit"].Path(o.ID) + `">редаґування</a></td></tr>` + "\n"
//	}
//	if canDelete {
//		htmlIndex += `<tr><td>- <a href="#" id="` + gt.listeners["deleteItem"].ID + `">вилучити запис</a>` +
//			`<input type="hidden" id="id" value="` + o.ID + `"></td></tr>` + "\n"
//	}
//
//	htmlIndex += `<tr><td>&nbsp;<br>` + items_comp.HTMLTags(o.Links, o.RView, "", "<br>- ") + "</td></tr>\n"
//
//	data := map[string]string{
//		"caput":      o.Name,
//		"titulus":    o.Name,
//		"share_tags": htmlShareTags,
//		"index":      htmlIndex,
//	}
//
//	if len(linkedObjects) > 0 {
//		htmlLinked = linksTitle + items_comp.HTMLIndex(userIS.String(), linkedObjects) + "\n<p>"
//	}
//
//	data["corpus"] = "\n" +
//		htmlContent + "\n<p>" +
//		items_comp.HTMLFiles(o.Links, gt.pxPreview) +
//		htmlShare +
//		htmlLinked
//
//	return data
//}
//
//var frontOps = map[string]viewshtml.Operator{
//	"file": items_comp.FileFrontOp,
//}
//
//func (gt *noteTranslator) NewItem(user *auth.User, o *notes.Item, context *genera.Context) genera.DataView {
//	data := map[string]string{
//		"genus": GenusKey,
//	}
//
//	rView := user.Identity().String()
//	if o != nil {
//		rView = o.RView
//	}
//
//	return map[string]string{
//		"caput":   "Нова нотатка",
//		"titulus": "Нова нотатка",
//		"corpus":  items_comp.HTMLEdit(user, noteCreateFields, data, nil, frontOps, rView, false),
//	}
//
//}
//
//func (gt *noteTranslator) Edit(user *auth.User, o *notes.Item, context *genera.Context) genera.DataView {
//	if o == nil {
//		return map[string]string{
//			"caput":   "Редаґування",
//			"titulus": "Редаґування",
//			"corpus":  txtNoNote,
//		}
//	}
//
//	userIS := user.Identity()
//
//	var htmlIndex string
//	if groups.OneOf(userIS, gt.ctrlOp, o.ROwner, o.Managers[rights.Delete]) {
//		htmlIndex += `<tr><td>- <a href="#" id="` + gt.listeners["deleteItem"].ID + `">вилучити запис</a>` +
//			`<input type="hidden" id="id" value="` + o.ID + `"></td></tr>` + "\n"
//	}
//
//	htmlIndex += `<tr><td>&nbsp;<br>` + items_comp.HTMLTags(o.Links, o.RView, "", "<br>- ") + "</td></tr>\n"
//
//	responseData := map[string]string{
//		"caput":   o.Name,
//		"titulus": o.Name,
//		"index":   htmlIndex,
//	}
//
//	if !groups.OneOf(userIS, gt.ctrlOp, o.ROwner, o.Managers[rights.Change]) {
//		responseData["corpus"] = rights.ErrNoRights.Error()
//		return responseData
//	}
//
//	dataRaw, _, err := gt.DataFromObject(userIS, o, nil)
//	if err != nil {
//		log.Println(err)
//		responseData["corpus"] = basis.ErrCantPerform.Error()
//		return responseData
//
//	}
//
//	publicChanges := o.Managers != nil && o.Managers[rights.Change] == o.RView
//	responseData["corpus"] = items_comp.HTMLEdit(user, noteUpdateFields, dataRaw.StringsMap(), nil, frontOps, o.RView, publicChanges)
//	return responseData
//
//}
