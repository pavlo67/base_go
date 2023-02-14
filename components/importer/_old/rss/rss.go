package rss

import (
	"time"

	"regexp"

	"github.com/mmcdole/gofeed"
	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/files"
	"github.com/pavlo67/punctum/interfaces/flow"
	"github.com/pavlo67/punctum/interfaces/importer"
	"github.com/pavlo67/punctum/interfaces/links"
)

// only to check compatibility with the interface
var testOperator importer.Operator = &RSS{}

// RSS implements importer.Operator interface ----------------------------------
type RSS struct {
	feedURL   string
	language  string
	items     []*gofeed.Item
	itemIndex int
}

// Init opens import session with selected data fount
func (r *RSS) Init(feedURL, dbKey string, testMode bool) error {
	r.feedURL = feedURL
	r.items = nil

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return err
	} else if feed == nil {
		return importer.ErrNoFount
	}

	r.itemIndex = -1
	r.language = feed.Language
	r.items = feed.Items

	return nil
}

// Next gets the next data entity from the fount
func (r *RSS) Next() (entity importer.Entity, err error) {
	r.itemIndex++

	if r.itemIndex < len(r.items) {
		return &Entity{rss: r, item: r.items[r.itemIndex]}, nil
	}

	return nil, importer.ErrNoMoreItems
}

func (r *RSS) Close() {
}

// Entity implements importer.Operator interface ----------------------------

// Entity implements importer.Entity interface.
type Entity struct {
	rss  *RSS
	item *gofeed.Item
}

// OriginalID gets the value from the imported entity
func (entity Entity) OriginalID() string {
	item := entity.item

	if item == nil {
		return ""
	}

	if item.GUID == "" {
		return item.Link
	}

	proto := reHTTP.FindString(item.GUID)

	if len(proto) == 0 {
		return item.GUID
	}

	feedURL := ""
	if entity.rss != nil {
		feedURL = entity.rss.feedURL
	}

	return feedURL + "#" + item.GUID

}

// Original gets a full original representation of the imported entity
func (entity Entity) Original() (interface{}, error) {
	return entity.item, nil
}

// type Item struct {
// 	Label        string                `json:"title,omitempty"`
// 	Description    string                `json:"description,omitempty"`
// 	Content        string                `json:"content,omitempty"`
// 	Link        string                `json:"link,omitempty"`
// 	Updated        string                `json:"updated,omitempty"`
// 	UpdatedParsed    *time.Time            `json:"updatedParsed,omitempty"`
// 	Published    string                `json:"published,omitempty"`
// 	PublishedParsed    *time.Time            `json:"publishedParsed,omitempty"`
// 	Author        *Person                `json:"author,omitempty"`
// 	GUID        string                `json:"guid,omitempty"`
// 	Image        *Image                `json:"image,omitempty"`
// 	Categories    []string            `json:"categories,omitempty"`
// 	Enclosures    []*Enclosure            `json:"enclosures,omitempty"`
// 	DublinCoreExt    *ext.DublinCoreExtension    `json:"dcExt,omitempty"`
// 	ITunesExt    *ext.ITunesItemExtension    `json:"itunesExt,omitempty"`
// 	Extensions    ext.Extensions            `json:"extensions,omitempty"`
// 	Custom        map[string]string        `json:"custom,omitempty"`
// }

// Object forms an interfaces.Object from the imported entity
func (entity Entity) Object() (obj *interfaces.Object, err error) {
	if entity.item == nil {
		return nil, importer.ErrNilItem
	}

	item := entity.item

	//language := ""
	//if entity.rss != nil {
	//	language = entity.rss.language
	//}

	createdLinks := []interfaces.Link{}
	if item.Author != nil {
		//email, err := url.Parse(item.Author.Email)
		//if err != nil {
		//	email = nil
		//}

		createdLinks = append(createdLinks, interfaces.Link{
			Type: "author",
			//Name:    []interfaces.Text{{Text: item.Author.Name, Language: language}},
			//Whereto: email,
			Name: item.Author.Name,
			To:   item.Author.Email,
		})
	}
	if item.Link != "" {
		//URL, err := url.Parse(item.Link)
		//if err != nil {
		//	URL = nil
		//}

		createdLinks = append(createdLinks, interfaces.Link{
			Type: "url",
			//Name:    []interfaces.Text{{Text: item.Link, Language: language}},
			//Whereto: URL,
			Name: item.Link,
			To:   item.Link,
		})
	}
	if item.Image != nil {
		//URL, err := url.Parse(item.Image.URL)
		//if err != nil {
		//	URL = nil
		//}
		createdLinks = append(createdLinks, interfaces.Link{
			Type: "image",
			//Name:    []interfaces.Text{{Text: item.Label, Language: language}},
			//Whereto: URL,
			Name: item.Title,
			To:   item.Image.URL,
		})
	}
	for _, category := range item.Categories {
		createdLinks = append(createdLinks, interfaces.Link{
			Type: links.TypeTag,
			//Name: []interfaces.Text{{Text: category, Language: language}},
			Name: category,
		})
	}

	return &interfaces.Object{
		//Name:    []interfaces.Text{{Text: item.Label, Language: language}},
		//Summary: []interfaces.Text{{Text: item.Description, Language: language}},
		Name:    item.Title,
		Content: item.Description + " " + item.Content,
		Links:   createdLinks,
	}, nil

}

var reHTTP = regexp.MustCompile("(?i)^https?://")

// FlowItem forms an flow.Item from the imported entity
func (entity Entity) FlowItem() (*flow.Item, error) {
	if entity.item == nil {
		return nil, importer.ErrNilItem
	}

	item := entity.item

	var createdAt time.Time
	if item.PublishedParsed != nil {
		createdAt = *item.PublishedParsed
	} else {
		createdAt = time.Now()
	}

	feedURL := ""
	if entity.rss != nil {
		feedURL = entity.rss.feedURL
	}

	flowItem := flow.Item{
		FountURL:   feedURL,
		OriginalID: entity.OriginalID(),
		Original:   interface{}(item),
		URL:        item.Link,
		Title:      item.Title,
		Summary:    item.Description,
		Content:    item.Content,
		CreatedAt:  createdAt,
	}
	if len(item.Enclosures) > 0 {
		flowItem.Media = &flow.ItemMedia{}
		for _, p := range item.Enclosures {
			flowItem.Media.Pictures = append(flowItem.Media.Pictures, flow.ItemPicture{ImageUrl: p.URL, HREFUrl: "#"})
		}
	}

	return &flowItem, nil

}

func (entity Entity) Files() ([]files.File, error) {
	return nil, interfaces.ErrNotImplemented
}
