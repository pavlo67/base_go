package htmlimporter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/strlib"
	"github.com/pavlo67/punctum/collector/importer"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/processor.old/news"

	"github.com/pkg/errors"
)

type ImporterHTML struct {
	feedURL string
	content string
	httpS   string
	domain  string
	url     string
	params  ImportParams

	title     string
	meta      string
	items     []string
	titles    []string
	itemIndex int

	rePartTitle *regexp.Regexp
}

type ImportParams struct {
	ImportStartRegexp    string   `json:"import_start_regexp"`
	ImportFinishRegexp   string   `json:"import_finish_regexp"`
	ImportSeparateRegexp string   `json:"import_separate_regexp"`
	AcceptableTags       []string `json:"acceptable_tags"`
	TitleRegexp          string   `json:"title_regexp"`
}

type Entity struct {
	mapObject map[string]string
}

var reTitle = regexp.MustCompile(`(?ims)<title.*?>(.*?)</title>`)
var reMeta = regexp.MustCompile(`(?ims)(<meta http-equiv="Content-Type".*?>)`)
var reHead = regexp.MustCompile(`(?ims)<head>.*?</head>`)
var reScript = regexp.MustCompile(`(?ims)<script.*?</script>`)
var reStyle = regexp.MustCompile(`(?ims)<style.*?</style>`)
var reHREF = regexp.MustCompile(`(?i) href=(['"])/`)
var reImage = regexp.MustCompile(`(?i) src=(['"])/`)
var reHTTP = regexp.MustCompile(`(?i)^(https?://)`)
var reDomain = regexp.MustCompile(`(?i)^(https?://.+?)[/\?$]`)
var reURL = regexp.MustCompile(`(?i)^(https?://.+)[/\?$]`)
var re2Slash = regexp.MustCompile(`(?i)<a href=(['"])//`)
var reTags = regexp.MustCompile(`(?ms)<(\w+)[ >]`)

var reAnyTag = regexp.MustCompile(`(?ms)</?[^>]+>`)
var reSpaces = regexp.MustCompile(`(?ms)\s+`)
var reItemHREF = regexp.MustCompile(`(?ims) href=['"]([^'"]+)['"]`)

var reBody = regexp.MustCompile(`(?msi).*?<body[^>]*>`)
var reEndBody = regexp.MustCompile(`(?msi)</body>.*`)

//var reURL = regexp.MustCompile(`^(.+)[/\?$]`)

func (h *ImporterHTML) Init(feedURL, importParams string, testMode bool) error {

	h.feedURL = feedURL
	h.itemIndex = -1
	h.items = []string{}
	//h.done = false
	resp, err := http.Get(feedURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	h.content = string(html)
	if importParams != "" {
		err = json.Unmarshal([]byte(importParams), &h.params)
		if err != nil {
			return errors.Wrapf(err, "can't Unmarshal: %v", importParams)
		}
		if h.params.TitleRegexp != "" {
			h.rePartTitle = regexp.MustCompile(h.params.TitleRegexp)
		}
	}
	arr := reDomain.FindStringSubmatch(feedURL)
	if len(arr) > 1 {
		h.domain = arr[1]
		//log.Println("UserIS domain:", arr[1])
	} else {
		log.Println("can't get domain from feedURL:", feedURL)
	}
	arr = reURL.FindStringSubmatch(feedURL)
	if len(arr) > 1 {
		h.url = arr[1]
		//log.Println("UserIS URL:", arr[1])
	} else {
		log.Println("can't get main url from feedURL:", feedURL)
	}
	arr = reHTTP.FindStringSubmatch(feedURL)
	if len(arr) > 1 {
		h.httpS = arr[1]
		//log.Println("UserIS http protocol:", arr[1])
	} else {
		log.Println("can't get http protocol from feedURL:", feedURL)
	}

	arr = reTitle.FindStringSubmatch(h.content)
	if len(arr) > 1 {
		h.title = arr[1]
		//log.Println("UserIS title:", arr[1])
	}
	arr = reMeta.FindStringSubmatch(h.content)
	if len(arr) > 1 {
		h.meta = arr[1]
		log.Println("UserIS meta:", arr[1])
	}

	h.content = reHead.ReplaceAllString(h.content, "<head>"+h.meta+"</head>")
	h.content = reScript.ReplaceAllString(h.content, "")
	h.content = reStyle.ReplaceAllString(h.content, "")
	h.content = re2Slash.ReplaceAllString(h.content, "<a href=${1}"+h.httpS)
	content := h.content
	if h.params.ImportStartRegexp != "" {
		re := regexp.MustCompile(`(?ism)^.*?(` + h.params.ImportStartRegexp + `)`)
		content = re.ReplaceAllString(content, "${1}")
	} else {
		content = reBody.ReplaceAllString(content, "<body>")
	}
	if h.params.ImportFinishRegexp != "" {
		re := regexp.MustCompile(`(?ism)(` + h.params.ImportFinishRegexp + `).*`)
		content = re.ReplaceAllString(content, "${1}")
	} else {
		content = reEndBody.ReplaceAllString(content, "</body>")
	}
	//log.Println("UserIS content", content)
	allTags := map[string]int{}
	if len(h.params.AcceptableTags) > 0 && h.params.AcceptableTags[0] != "" {
		// another tags.comp must be to drop
		tags := reTags.FindAllStringSubmatch(content, -1)
		for _, t := range tags {
			allTags[t[1]] = 1
		}
		for t := range allTags {
			needDrop := true
			for _, t0 := range h.params.AcceptableTags {
				if t == t0 {
					needDrop = false
					break
				}
			}
			if needDrop {
				//log.Println("UserIS tag for clear:", t)
				reTag := regexp.MustCompile(`(?ims)<` + t + `[^>]*>`)
				content = reTag.ReplaceAllString(content, "")
				reTag = regexp.MustCompile(`</` + t + `>`)
				content = reTag.ReplaceAllString(content, "")
			}
		}
	}
	content = reHREF.ReplaceAllString(content, " href=${1}"+h.domain+"/")
	content = reImage.ReplaceAllString(content, " src=${1}"+h.domain+"/")
	if h.params.ImportSeparateRegexp != "" {
		re := regexp.MustCompile(`(?ism)` + h.params.ImportSeparateRegexp)
		arr := re.Split(content, -1)
		if len(arr) > 1 {
			h.items = arr[1:]
		}
	} else {
		h.items = []string{content}
	}

	return nil
}

func (h *ImporterHTML) Next() (importer.Entity, error) {

	h.itemIndex++
	if h.itemIndex < len(h.items) {
		var title, href string
		content := h.items[h.itemIndex]
		//	get part title
		if h.rePartTitle != nil {
			arr := h.rePartTitle.FindStringSubmatch(content)
			if len(arr) > 1 {
				title = arr[1]
			}
		}
		if title == "" {
			//	get 250 symbols as title
			title = reAnyTag.ReplaceAllString(content, "")
			title = strlib.ReSpaces.ReplaceAllString(title, " ")
		}
		if len(title) > 255 {
			title = title[:250]
		}
		arr := reItemHREF.FindStringSubmatch(content)
		if len(arr) > 1 {
			href = arr[1]
		}

		return &Entity{
			mapObject: map[string]string{
				"title":   title,
				"href":    href,
				"meta":    h.meta,
				"feedURL": h.feedURL,
				"content": h.items[h.itemIndex],
			},
		}, nil
	}
	return nil, importer.ErrNoMoreItems
}

func (h *ImporterHTML) Close() {
}

var reDelimiter = regexp.MustCompile(`(?ims)[\.,\?\s\[\]{}\(\)-\+\*;:'"%\^#<>\\\/=!~]+`)

func (e Entity) OriginalID() string {
	or := reDelimiter.ReplaceAllString(e.mapObject["content"], "")
	if len(or) > 255 {
		return or[:254]
	}

	return or
}

func (e Entity) Original() (interface{}, error) {
	return interface{}(e.mapObject), nil
}

func (e Entity) Object() (*notes.Item, error) {
	o := notes.Item{
		Name:     e.mapObject["title"],
		Content:  e.mapObject["content"],
		GlobalIS: strlib.RandomString(60), // todo!!! temporary

	}
	return &o, nil
}

func (e Entity) FlowItem() (*news.Item, error) {
	flowItem := news.Item{
		SourceURL:  e.mapObject["feedURL"],
		OriginalID: e.OriginalID(),
		Original:   interface{}(""),
		URL:        e.mapObject["href"],
		Title:      e.mapObject["title"],
		Content:    e.mapObject["content"],
	}

	return &flowItem, nil
}

func (e Entity) Files() ([]files.File, error) {
	return nil, basis.ErrNotImplemented
}
