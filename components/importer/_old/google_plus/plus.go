package google_plus

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/files"
	"github.com/pavlo67/punctum/interfaces/flow"
	"github.com/pavlo67/punctum/interfaces/importer"
	"google.golang.org/api/plus/v1"
)

type PLUS struct {
	feedURL   string
	language  string
	items     []*plus.Activity
	itemIndex int

	ApiKey string
	//ApiID      string
	//ApiSecret  string
	//PathToJSON string

}

// Retrieves a token from a local file.
//func tokenFromFile(file string) (*oauth2.Token, error) {
//	f, err := os.Open(file)
//	defer f.Close()
//	if err != nil {
//		return nil, err
//	}
//	tok := &oauth2.Token{}
//	err = json.NewDecoder(f).Decode(tok)
//	log.Println("IS g+ token", tok, file)
//
//	return tok, err
//}
//
//func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
//	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
//	fmt.Printf("Go to the following link in your viewer then type the "+
//		"authorization code: \n%v\n", authURL)
//
//	var authCode string
//	if _, err := fmt.Scan(&authCode); err != nil {
//		log.Fatalf("Unable to read authorization code: %v", err)
//	}
//
//	tok, err := config.Exchange(oauth2.NoContext, authCode)
//	if err != nil {
//		log.Fatalf("Unable to retrieve token from web: %v", err)
//	}
//	return tok
//}
//
//// Saves a token to a file path.
//func saveToken(path string, token *oauth2.Token) {
//	fmt.Printf("Saving credential file to: %s\n", path)
//	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
//	defer f.Close()
//	if err != nil {
//		log.Fatalf("Unable to cache oauth token: %v", err)
//	}
//	json.NewEncoder(f).Encode(token)
//}

func (pc *PLUS) Init(feedURL, dbKey string, testMode bool) error {

	//var config = &oauth2.Config{
	//	ClientID:     pc.ApiID,     // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
	//	ClientSecret: pc.ApiSecret, // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
	//	Endpoint:     google.Endpoint,
	//	Scopes:       []string{urlshortener.UrlshortenerScope},
	//}
	//
	//log.Println("IS g+ credentials", pc.ApiID, pc.ApiSecret)
	////tok := getTokenFromWeb(config)
	////saveToken(pc.PathToJSON + "555", tok)
	//
	//ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
	//	Transport: &transport.APIKey{Key: pc.ApiKey},
	//})
	//
	////token, err := config.Exchange(ctx, "")
	//token, err := tokenFromFile(pc.PathToJSON)
	//if err != nil {
	//	log.Println("IS err1", err)
	//}
	//
	//oauthHttpClient := config.Client(ctx, token)
	////oauthHttpClient :=config.Client(context.Background(), token)
	//plusService, err := plus.New(oauthHttpClient)
	//if err != nil {
	//	log.Println("IS err2", err)
	//}
	//activities, err := plusService.Activities.List(userID, "public").Do()
	//if err != nil {
	//	log.Println("IS err3", err)
	//}
	//log.Println("IS:", activities)
	pc.feedURL = feedURL
	feedURL += "?key=" + pc.ApiKey
	var responseJSON plus.ActivityFeed
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return err
	}
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//log.Println("IS body!!!", body, err)
	err = json.Unmarshal(body, &responseJSON)
	if err != nil {
		return err
	}
	pc.items = responseJSON.Items
	pc.itemIndex = -1
	return nil
}

type Entity struct {
	plus *PLUS
	item *plus.Activity
}

func (pc *PLUS) Next() (entity importer.Entity, err error) {
	pc.itemIndex++

	if pc.itemIndex < len(pc.items) {
		//log.Println("IS resp:", pc.items[pc.itemIndex].Label)
		return &Entity{plus: pc, item: pc.items[pc.itemIndex]}, nil
	}

	return nil, importer.ErrNoMoreItems
}

func (pc *PLUS) Close() {
}

func (e Entity) OriginalID() string {
	return e.item.Id
}

func (e Entity) Original() (interface{}, error) {
	return e.item, nil
}

func (e Entity) Object() (obj *interfaces.Object, err error) {
	return nil, nil
}

var reUTF8Symbols = regexp.MustCompile(`\p{S}+`)

func (e Entity) FlowItem() (*flow.Item, error) {

	item := e.item

	var createdAt time.Time
	if item.Updated != "" {
		//t, err := time.Parse("Mon Jan 02 15:04:05 -0700 2006", item.Updated)
		t, err := time.Parse("2006-01-02T15:04:05Z", item.Updated)
		//2018-05-17T11:14:30.475Z 2006-01-02T15:04:05.475Z
		if err != nil {
			log.Println("can't time parse: ", item.Updated, err)
			createdAt = time.Now()
		}
		createdAt = t
	} else {
		createdAt = time.Now()
	}

	feedURL := ""
	if e.plus != nil {
		feedURL = e.plus.feedURL
	}

	var media = []flow.ItemPicture{}
	for _, m := range item.Object.Attachments {
		if len(m.Thumbnails) > 0 {
			for _, t := range m.Thumbnails {
				media = append(media,
					flow.ItemPicture{
						ImageUrl: t.Image.Url,
					})
			}
		} else {
			media = append(media,
				flow.ItemPicture{
					ImageUrl: m.FullImage.Url,
				})
		}
	}
	if reUTF8Symbols.MatchString(item.Title) {
		item.Title = reUTF8Symbols.ReplaceAllString(item.Title, "***")
		//log.Println("IS clean text:", item.Text)
	}
	if reUTF8Symbols.MatchString(item.Object.Content) {
		item.Object.Content = reUTF8Symbols.ReplaceAllString(item.Object.Content, "***")
		//log.Println("IS clean text:", item.Text)
	}
	flowItem := flow.Item{
		FountURL:   feedURL,
		OriginalID: e.OriginalID(),
		// todo! on serverhttp:  can't write item  to Original.
		// todo!Incorrect string value: '\xF0\x9F\x87\xBA\xF0\x9F...' for column 'original' at row 1
		Original: interface{}(item),
		URL:      item.Url,
		Title:    item.Title,
		//Summary:    item.Description,
		Content:   item.Object.Content,
		CreatedAt: createdAt,
		Media: &flow.ItemMedia{
			Pictures: media,
		},
	}

	return &flowItem, nil
}

func (entity Entity) Files() ([]files.File, error) {
	return nil, interfaces.ErrNotImplemented
}
