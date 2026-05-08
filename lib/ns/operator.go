package ns

import (
	"regexp"
)

var reSpaces = regexp.MustCompile(`\s+`)

type Item struct {
	Proto    string `json:"proto,omitempty"    bson:"proto,omitempty"`
	Host     string `json:"host,omitempty"     bson:"host,omitempty"`
	Path     string `json:"path,omitempty"     bson:"path,omitempty"`
	Fragment string `json:"fragment,omitempty" bson:"fragment,omitempty"`
}

func (item *Item) IsValid() bool {
	if item == nil {
		return false
	}
	return reSpaces.ReplaceAllString(item.Host, "") != "" &&
		reSpaces.ReplaceAllString(item.Path, "") != "" &&
		reSpaces.ReplaceAllString(item.Fragment, "") != ""
}

func (item *Item) URN() URN {
	if item == nil {
		return ""
	}

	host := reSpaces.ReplaceAllString(item.Host, "")
	path := reSpaces.ReplaceAllString(item.Path, "")
	id := reSpaces.ReplaceAllString(item.Fragment, "")

	if len(id) > 0 {
		return URN(host + path + IDDelim + id)
	} else if len(path) > 1 {
		return URN(host + path)
	} else {
		return URN(host)
	}
}

func (item *Item) String() string {
	return string(item.URN())
}

//func FromURLRaw(urlRaw string) Item {
//	urlWithoutProto := reProto.ReplaceAllString(strings.TrimSpace(urlRaw), "")
//	domain := reHostDelim.ReplaceAllString(urlWithoutProto, "")
//
//	// TODO!!! clean more
//
//	return Item{
//		Host: domain,
//		Path:   urlWithoutProto[len(domain):],
//	}
//}
