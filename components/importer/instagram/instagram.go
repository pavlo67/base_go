package instagramimporter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/collector/importer"
	"github.com/pavlo67/punctum/notebook/notes"
	"github.com/pavlo67/punctum/processor.old/news"

	"github.com/pkg/errors"
)

type Instagram struct {
	ID        string
	Secret    string
	Token     string
	feedURL   string
	items     []InstagramItem
	itemIndex int
}
type InstagramItem struct {
	UserID       string
	UserName     string
	FullName     string
	UserImageURL string
	ItemID       string
	ItemURL      string
	ItemImageURL string
	ItemCaption  string
}

//type I1 struct {
//	EntryData I2 `json:"entry_data"`
//}
//type I2 struct {
//	ProfilePage []I3 `json:"ProfilePage"`
//}
//type I3 struct {
//	Graphql interface{} `json:"graphql"`
//}
//

type InstagramJSON struct {
	EntryData EntryData `json:"entry_data"`
}
type EntryData struct {
	ProfilePage []ProfilePage `json:"ProfilePage"`
}
type ProfilePage struct {
	Graphql Qraphql `json:"graphql"`
}
type Qraphql struct {
	User User `json:"user"`
}
type User struct {
	ID                       string                   `json:"id"`
	UserName                 string                   `json:"username"`
	FullName                 string                   `json:"full_name"`
	UserImageURL             string                   `json:"profile_pic_url"`
	EdgeOwnerToTimelineMedia EdgeOwnerToTimelineMedia `json:"edge_owner_to_timeline_media"`
}
type EdgeOwnerToTimelineMedia struct {
	Edges []Edge `json:"edges"`
}
type Edge struct {
	Node Node `json:"node"`
}
type Node struct {
	ID                 string             `json:"id"`
	DisplayURL         string             `json:"display_url"`
	ShortCode          string             `json:"shortcode"`
	EdgeMediaToCaption EdgeMediaToCaption `json:"edge_media_to_caption"`
}
type EdgeMediaToCaption struct {
	Edges []CaptionEdge `json:"edges"`
}
type CaptionEdge struct {
	Node CaptionNode `json:"node"`
}
type CaptionNode struct {
	Text string `json:"text"`
}

type Entity struct {
	instagram *Instagram
	item      InstagramItem
	feedURL   string
}

//var reUserID = regexp.MustCompile(`(?i)instagram\.com/([^/]+)/`)
var reJSON = regexp.MustCompile(`(?ism)window._sharedData = (.*?);</script>`)

func (i *Instagram) Init(feedURL, dbKey string, testMode bool) error {

	i.feedURL = feedURL
	i.items = nil
	resp, err := http.Get(feedURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//log.Println(string(html))
	var userInstagramMedias InstagramJSON
	if reJSON.MatchString(string(html)) {
		arr := reJSON.FindStringSubmatch(string(html))
		//log.Println("UserIS JSON: ", arr[1])
		err = json.Unmarshal([]byte(arr[1]), &userInstagramMedias)
		if err != nil {
			return errors.Wrapf(err, "no find instagram json content on page: "+feedURL)
		}
		//log.Println("UserIS struct: ", userInstagramMedias)

		//var i1 I1
		//err = json.Unmarshal([]byte(arr[1]), &i1)
		//if err != nil {
		//	return nil.Wrapf(err, "no find instagram json content on page: " + feedURL)
		//}
		//log.Println("UserIS struct1: ", i1)

	} else {
		return errors.New("no find instagram json content on page: " + feedURL)
	}
	for _, profile := range userInstagramMedias.EntryData.ProfilePage {
		for _, edge := range profile.Graphql.User.EdgeOwnerToTimelineMedia.Edges {
			var item InstagramItem
			item.ItemID = edge.Node.ID
			item.ItemURL = edge.Node.ShortCode
			item.ItemImageURL = edge.Node.DisplayURL
			for _, caption := range edge.Node.EdgeMediaToCaption.Edges {
				item.ItemCaption += caption.Node.Text + " "
			}
			item.UserID = profile.Graphql.User.ID
			item.UserName = profile.Graphql.User.UserName
			item.FullName = profile.Graphql.User.FullName
			i.items = append(i.items, item)
		}
	}
	//log.Println("UserIS Items:", i.items.comp)

	//var userID string
	//arr := reUserID.FindStringSubmatch(feedURL)
	//if len(arr)>1 {
	//	userID = arr[1]
	//} else {
	//	return nil.New("can't get instagram user Label from url: " + feedURL)
	//}
	//api := instagram.New(i.Label, i.Secret, i.SessionIDs, true)
	//params := url.Info{}
	//params.Set("name", userID)
	//r,err := api.GetUserSearch(params)
	//if err != nil {
	//	return err
	//}
	//
	//log.Println("UserIS confidenter.comp:", r.Users)
	//res, err := api.GetUserRecentMedia(userID, params)
	//if err != nil {
	//	return err
	//}
	//
	////doneChan := make(chan bool)
	//
	//mediaChan, errChan := api.IterateMedia(res, nil)
	//var linker_server_http int
	//for media := range mediaChan {
	//	//media.Images
	//	log.Println("UserIS user:", media.User.Username)
	//	log.Println("UserIS media:", media.Images)
	//	i.items.comp = append(i.items.comp, media)
	//	linker_server_http++
	//	if linker_server_http > 2 {
	//		break
	//	}
	//	//if media.User.Username != "ladygaga" {
	//	//
	//	//}
	//	//if isDone(media) {
	//	//	close(doneChan) // Signal to iterator to quit
	//	//	break
	//	//}
	//}
	//if err := <-errChan; err != nil {
	//	return err
	//}
	i.itemIndex = -1
	return nil
}

// Next gets the next data entity from the fount
func (i *Instagram) Next() (entity importer.Entity, err error) {

	i.itemIndex++
	if i.itemIndex < len(i.items) {
		return &Entity{instagram: i, item: i.items[i.itemIndex], feedURL: i.feedURL}, nil
	}

	return nil, importer.ErrNoMoreItems
}

func (i *Instagram) Close() {
}

func (entity Entity) OriginalID() string {
	return entity.item.ItemID
}

// Original gets a full original representation of the imported entity
func (entity Entity) Original() (interface{}, error) {
	return entity.item, nil
}

func (entity Entity) Object() (obj *notes.Item, err error) {
	return nil, nil
}

func (entity Entity) Files() ([]files.File, error) {
	return nil, basis.ErrNotImplemented
}

var reHashTag = regexp.MustCompile(`#([^ ]+)`)
var reUTF8Symbols = regexp.MustCompile(`\p{S}+`)

func (entity Entity) FlowItem() (*news.Item, error) {

	arr := reHashTag.FindAllStringSubmatch(entity.item.ItemCaption, -1)
	var hashTags []string
	for _, t := range arr {
		hashTags = append(hashTags, t[1])
	}
	entity.item.ItemCaption = reHashTag.ReplaceAllString(entity.item.ItemCaption, "")
	if reUTF8Symbols.MatchString(entity.item.ItemCaption) {
		entity.item.ItemCaption = reUTF8Symbols.ReplaceAllString(entity.item.ItemCaption, "***")
	}
	if strings.TrimSpace(entity.item.ItemCaption) == "" {
		entity.item.ItemCaption = "no comment"
	}
	flowItem := news.Item{
		SourceURL:  entity.feedURL,
		OriginalID: entity.OriginalID(),
		Title:      entity.item.ItemCaption,
		URL:        "https://www.instagram.com/p/" + entity.item.ItemURL,
		Media: &news.ItemMedia{
			HashTags: hashTags,
			Pictures: []news.ItemPicture{
				{
					ImageUrl: entity.item.ItemImageURL,
				},
			},
		},
	}
	return &flowItem, nil
}
