package records

//var Fields = views_html.Fields{
//	// {"visibility", "тип", "select", "", common.Map{"class": "ut"}},
//	{"id", "", "hidden", "", nil},
//	{"title", "заголовок", "", "", nil},
//	{"summary", "коротко", "", "", nil},
//	{"type", "тип", "", "", nil},
//	{"data", "опис", "textarea", "", common.Map{"rows": 35}},
//
//	// {"embedded", "", "", "", common.Map{"class": "ut"}},
//	// {"files", "Завантажити файл", "file", "", common.Map{"class": "ut"}},
//
//	{"tags", "теми, розділи", "tag-it", "", nil},
//	{"created_at", "запис створено", "view", "datetime", common.Map{views_html.NotEmptyKey: true}},
//	{"updated_at", "востаннє відредаґовано", "view", "datetime", common.Map{views_html.NotEmptyKey: true}},
//
//	{"save", "", "submit", "зберегти запис", nil},
//}
//
//var reTagsSplit = regexp.MustCompile(`\s+;\s+`)
//
//func ItemFromMap(data common.Map) (*Item, error) {
//	tags := data.Strings("tags_prepared")
//	if len(tags) < 1 {
//		tagsStr := strings.TrimSpace(data.StringDefault("tags", ""))
//		tags = reTagsSplit.Split(tagsStr, -1)
//	}
//
//	return &Item{
//		ID: ID(data.StringDefault("id", "")),
//		Description: crud.Description{
//			// URN:          "",
//			Tags: tags,
//			// RelationsMap: nil,
//			// OwnerNSS:     "",
//			// ViewerNSS:    "",
//			// History:      nil,
//		},
//		Record: Record{
//			Content: Content{
//				Title:   data.StringDefault("title", ""),
//				Summary: data.StringDefault("summary", ""),
//				Type:    data.StringDefault("type", ""),
//				Records:    data.StringDefault("data", ""),
//			},
//			// Additions: nil,
//		},
//	}, nil
//}
//
//func MapFromItem(item Item) views_html.ValuesString {
//	return views_html.ValuesString{
//		"id":         string(item.ID),
//		"title":      item.Title,
//		"summary":    item.Summary,
//		"type":       item.Type,
//		"data":       item.Records,
//		"tags":       strings.Join(item.Tags, "; "),
//		"created_at": timelib.String(&item.At),
//		"updated_at": timelib.String(item.UpdatedAt),
//
//		//"visibility": item.,
//		//"embedded": item.,
//		//"files": item.,
//	}
//
//}
