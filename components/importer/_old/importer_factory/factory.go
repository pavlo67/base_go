package importer_factory

import (
	"log"
	"regexp"

	"github.com/mmcdole/gofeed"
	"github.com/pavlo67/punctum/interfaces/importer"
	"github.com/pavlo67/punctum/interfaces/importer/google_plus"
	"github.com/pavlo67/punctum/interfaces/importer/instagram"
	"github.com/pavlo67/punctum/interfaces/importer/rss"
	"github.com/pavlo67/punctum/interfaces/importer/twitter"
)

const NoneImporter = "NoneImporter"

type Factory interface {
	Init(url string, params interface{}) (importer.Operator, error)
}

var reTwitter = regexp.MustCompile(`(?i)(\.|)twitter\.com`)
var reRSS = regexp.MustCompile(`(?i)rss`)
var reXML = regexp.MustCompile(`(?i)\.(xml|atom)`)
var reGooglePlus = regexp.MustCompile(`(?i)(googleapis\.com/plus|plus\.google)`)
var reInstagram = regexp.MustCompile(`(?i)instagram\.com`)

func CheckURLType(url string) string {

	if reTwitter.MatchString(url) {
		return twitterimporter.InterfaceKey
	} else if reInstagram.MatchString(url) {
		return instagramimporter.InterfaceKey
	} else if reGooglePlus.MatchString(url) {
		return google_plus.InterfaceKey
	} else if reRSS.MatchString(url) {
		return rss.InterfaceKey
	} else if reXML.MatchString(url) {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(url)
		if err == nil && feed != nil {
			return "rss"
		} else {
			log.Println("can't check url:", url, "as rss file", err)
		}
	}
	log.Println("can't find url type for:", url, "; has choosen html type")
	return NoneImporter
}
