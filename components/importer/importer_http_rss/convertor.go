package importer_http_rss

import (
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"

	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/dataimporter"
	"github.com/pavlo67/workshop/components/tagger"
)

var _ data.Convertor = &Item{}

type Item struct {
	sourceTime time.Time
	sourceURL  string
	feedItem   *gofeed.Item
}

func (item *Item) GetData() (*data.Item, error) {
	if item == nil || item.feedItem == nil {
		return nil, nil
	}

	feedItem := item.feedItem

	originalID := feedItem.GUID
	if originalID == "" {
		originalID = feedItem.Link
	}

	sourceTime := item.sourceTime
	// if feedItem.PublishedParsed != nil {
	//	sourceTime = *feedItem.PublishedParsed
	// }

	// origin, _ := json.Marshal(feedItem)

	key := dataimporter.Identity(item.sourceURL, originalID).Key()
	history := []crud.Action{
		{Key: crud.CreatedAction, DoneAt: sourceTime},
		{Key: dataimporter.ActionKey, DoneAt: time.Now(), ActorKey: &key},
	}

	var embedded []data.Item

	if feedItem.Image != nil {
		embedded = append(embedded, data.Item{
			Data:    crud.Data{TypeKey: crud.HRefImageTypeKey},
			URL:     feedItem.Image.URL,
			Title:   feedItem.Image.Title,
			History: history,
		})
	}

	if len(feedItem.Enclosures) > 0 {
		for _, p := range feedItem.Enclosures {
			embedded = append(embedded, data.Item{
				Data:    crud.Data{TypeKey: crud.HRefTypeKey},
				URL:     p.URL,
				Title:   p.Type + ": " + p.Length,
				History: history,
			})
		}
	}

	var items []tagger.Tag
	for _, c := range feedItem.Categories {
		items = append(items, tagger.Tag{Label: c})
	}

	var dataItem = data.Item{
		Key:      key,
		URL:      feedItem.Link,
		Data:     crud.Data{TypeKey: crud.StringTypeKey, Content: []byte(feedItem.Content)},
		Title:    feedItem.Title,
		Summary:  feedItem.Description,
		Embedded: embedded,
		Tags:     items,
		History:  history,
	}

	sourceKey := dataimporter.SourceKey(history)
	if sourceKey != nil {
		dataItem.Key = *sourceKey
	}

	return &dataItem, nil

}

func (*Item) SaveData(data.Item) error {
	return common.ErrNotImplemented
}

// language := feed.Language

// type Census struct {
// 	Label        string                `json:"title,omitempty"`
// 	Description    string                `json:"description,omitempty"`
// 	Contentus        string                `json:"content,omitempty"`
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
