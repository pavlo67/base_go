package twitterimporter

import (
	"log"
	"regexp"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/collector/importer"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/processor.old/news"
)

type Twitter struct {
	feedURL   string
	language  string
	itemIndex int
	items     []twitter.Tweet

	Key         string
	KeySecret   string
	Token       string
	TokenSecret string
}

type Entity struct {
	twitter *Twitter
	item    twitter.Tweet
}

var reTwitterUser = regexp.MustCompile(`.*/`)
var reTwitterUser2 = regexp.MustCompile(`\?.*`)

// Prepare opens import session with selected data fount
func (t *Twitter) Init(feedURL, dbKey string, testMode bool) error {

	t.feedURL = reTwitterUser2.ReplaceAllString(feedURL, "")
	twitterUser := reTwitterUser.ReplaceAllString(feedURL, "")
	twitterUser = reTwitterUser2.ReplaceAllString(twitterUser, "")
	//log.Println("UserIS: scan twitter user:", twitterUser)
	//log.Println("UserIS: token:", t.SessionIDs)

	//config := &oauth2.Config{}
	//token := &oauth2.SessionIDs{AccessToken: t.SessionIDs}
	//// http.Client will automatically authorize Requests
	//httpClient := config.Client(context.TODO(), token)

	config := oauth1.NewConfig(t.Key, t.KeySecret)
	token := oauth1.NewToken(t.Token, t.TokenSecret)
	//http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	// Twitter client
	client := twitter.NewClient(httpClient)

	//user, resp, err := client.Users.Show(&twitter.UserShowParams{
	//	ScreenName: twitterUser,
	//})
	//log.Println("UserIS: U:", user, "\nR:", resp, "\nE:", err, "\n\n" )
	twits, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: twitterUser,
	})
	//log.Println("UserIS: T0:", twits[0].User, twits[0].Label, twits[0].FullText, twits[0].IDStr, twits[0].Lang, twits[0].Text, twits[0].SavedAt  )
	//for i := range twits {
	//	log.Println("UserIS: T0:", twits[i].Label, twits[i].Text, "\nE:")
	//	for _, h := range twits[i].Entities.Hashtags {
	//		log.Println("hashtag:", h.Text)
	//	}
	//	log.Println("MEDIA:", twits[i].Entities.Media)
	//	for _, m := range twits[i].Entities.Media {
	//		log.Println("mediaURL:", m.MediaURL)
	//		log.Println("mediaURLhref:", m.URLEntity.ExpandedURL)
	//	}
	//}
	//log.Fatal("!!! stop for test")
	if err != nil {
		return err
	} else if twits == nil {
		return importer.ErrNoFount
	}
	t.items = twits
	t.itemIndex = -1
	return nil
}

// Next gets the next data entity from the fount
func (t *Twitter) Next() (entity importer.Entity, err error) {

	t.itemIndex++
	if t.itemIndex < len(t.items) {
		return &Entity{twitter: t, item: t.items[t.itemIndex]}, nil
	}

	return nil, importer.ErrNoMoreItems
}

func (t *Twitter) Close() {
}

func (entity Entity) OriginalID() string {

	return entity.item.IDStr
}

// Original gets a full original representation of the imported entity
func (entity Entity) Original() (interface{}, error) {
	return entity.item, nil
}

func (entity Entity) Object() (obj *notes.Item, err error) {
	return nil, nil
}

var reUTF8Symbols = regexp.MustCompile(`\p{S}+`)

func (entity Entity) FlowItem() (*news.Item, error) {

	item := entity.item

	var createdAt time.Time
	if item.CreatedAt != "" {
		t, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", item.CreatedAt)
		if err != nil {
			log.Println("can't time parse: ", item.CreatedAt, err)
			createdAt = time.Now()
		}
		createdAt = t
	} else {
		createdAt = time.Now()
	}

	feedURL := ""
	if entity.twitter != nil {
		feedURL = entity.twitter.feedURL
	}

	var hashTags = []string{}
	for _, h := range item.Entities.Hashtags {
		hashTags = append(hashTags, h.Text)
	}
	var media = []news.ItemPicture{}
	for _, m := range item.Entities.Media {
		media = append(media,
			news.ItemPicture{
				ImageUrl: m.MediaURL,
				HREFUrl:  m.URLEntity.ExpandedURL,
			})
	}
	if reUTF8Symbols.MatchString(item.Text) {
		item.Text = reUTF8Symbols.ReplaceAllString(item.Text, "***")
		//log.Println("UserIS clean text:", item.Text)
	}
	flowItem := news.Item{
		SourceURL:  feedURL,
		OriginalID: entity.OriginalID(),
		// todo! on serverhttp_jschmhr:  can't write item  to Original.
		// todo!Incorrect string value: '\xF0\x9F\x87\xBA\xF0\x9F...' for column 'original' at row 1
		//Original:   interface{}(item),
		URL:   feedURL + "/status/" + entity.OriginalID(),
		Title: item.Text,
		//Summary:    item.Description,
		//Contentus:    item.Contentus,
		CreatedAt: createdAt,
		Media: &news.ItemMedia{
			HashTags: hashTags,
			Pictures: media,
		},
	}

	return &flowItem, nil
}

func (entity Entity) Files() ([]files.File, error) {
	return nil, basis.ErrNotImplemented
}
