package importer_http_rss

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/logger"

	"github.com/pavlo67/workshop/components/dataimporter"
)

func New(l logger.Operator) (dataimporter.Operator, error) {
	if l == nil {
		return nil, errors.New("on New(): nil logger")
	}

	return &rss{l}, nil
}

var _ dataimporter.Operator = &rss{}

type rss struct {
	l logger.Operator
}

//var reHTTP = regexp.MustCompile("(?i)^https?://")

const onGet = "on rss.Get(): "

func (r *rss) Get(feedURL string) (*dataimporter.DataSeries, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)

	if err != nil {
		// return nil, errors.Wrapf(err, onGet+"can't fp.ParseURL(%s)", feedURL)

		r.l.Warn(errors.Wrapf(err, onGet+"can't fp.ParseURL(%s)", feedURL))

		resp, err := http.Get(feedURL)
		if err != nil {
			return nil, errors.Wrapf(err, onGet+"can't http.Get(%s)", feedURL)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, onGet+"can't ioutil.ReadAll(%#v)", resp.Body)
		}

		feed, err = fp.ParseString(string(body))
		if err != nil {
			return nil, errors.Wrapf(err, onGet+"can't .ParseString(%s)", body)
		} else if feed == nil {
			return nil, errors.Errorf(onGet+"no feed obtained with .ParseString(%s)", body)
		}

	} else if feed == nil {
		return nil, errors.Errorf(onGet+"no feed obtained with .ParseURL(%s)", feedURL)
	}

	now := time.Now()

	series := dataimporter.DataSeries{URL: feedURL, CreatedAt: now}

	for _, feedItem := range feed.Items {
		item := &Item{
			sourceTime: now,
			sourceURL:  feedURL,
			feedItem:   feedItem,
		}

		dataItem, err := item.GetData()
		if err != nil {
			return &series, errors.Wrapf(err, onGet+"can't .GetData(%#v)", feedItem)
		}
		if dataItem == nil {
			// ??? wtf
			continue
		}

		series.Data = append(series.Data, *dataItem)
	}

	return &series, nil
}
