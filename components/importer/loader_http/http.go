package loader_http

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/files"
	"github.com/pavlo67/common/common/httplib"
	"github.com/pavlo67/common/common/selectors"
	"github.com/pavlo67/data/components/loader"
)

var _ loader.Operator = &loaderHTTP{}

type loaderHTTP struct {
	pathToStoreDefault string
}

const onNew = "on loaderHTTP.New(): "

func New(pathToStoreDefault string) (loader.Operator, db.Cleaner, error) {
	if strings.TrimSpace(pathToStoreDefault) == "" {
		pathToStoreDefault = "./"
	}
	pathToStoreDefaultFinal, err := filelib.Dir(pathToStoreDefault)
	if err != nil {
		return nil, nil, errors.Wrapf(err, onNew+"can't filelib.GetDir('%s', './')", pathToStoreDefault)
	}

	flOp := loaderHTTP{
		pathToStoreDefault: pathToStoreDefaultFinal,
	}

	return &flOp, nil, nil
}

type toPrepare struct {
	url      string
	fileType string
	fileName string
	priority int
}

const onLoad = "on loaderHTTP.Load(): "

func (flOp *loaderHTTP) Load(urlToLoad, pathToStore string, priority Priority) (*files.Item, error) {
	if strings.TrimSpace(pathToStore) == "" {
		pathToStore = flOp.pathToStoreDefault
	}

	pathToStoreFinal, err := filelib.SubDirUnique(pathToStore)
	if err != nil {
		return nil, errors.Wrapf(err, onLoad+"can't filelib.SubDirUnique('%s')", pathToStore)
	}

	if priority == nil {
		priority = PriorityDefault(urlToLoad, false)
	}

	var fileIndex int

	fileName, fileType, err := httplib.DownloadFile(urlToLoad, pathToStoreFinal, fileIndex, 0644)
	// TODO!!! postpone errors
	if err != nil {
		return nil, err
	}
	fileIndex++

	filesToPrepare := []toPrepare{{urlToLoad, fileType, fileName, 1}}

	for len(filesToPrepare) > 0 {
		fileToPrepare := filesToPrepare[0]
		filesToPrepare = filesToPrepare[1:]

		var posterior []toPrepare

		posterior, fileIndex, err = flOp.PreparePosterior(fileToPrepare, pathToStoreFinal, fileIndex, priority)
		// TODO!!! postpone errors
		if err != nil {
			return nil, err
		}

		if len(posterior) > 0 {
			filesToPrepare = append(filesToPrepare, posterior...)
			sort.Slice(filesToPrepare, func(i, j int) bool { return filesToPrepare[j].priority < filesToPrepare[i].priority })
		}
	}

	now := time.Now()

	return &files.Item{
		Path: pathToStoreFinal,
		Origin: flow.Origin{
			Source: urlToLoad,
			Time:   &now,
		},
	}, nil
}

const onPreparePosterior = "on loaderHTTP.PreparePosterior(): "

func (flOp *loaderHTTP) PreparePosterior(fileToPrepare toPrepare, pathToStore string, fileIndex int, priority Priority) ([]toPrepare, int, error) {

	var filesToPrepare []toPrepare

	file, err := os.Open(fileToPrepare.fileName)
	if err != nil {
		return nil, fileIndex, errors.Wrapf(err, "can't os.Open(%s)", fileToPrepare.fileName)
	}

	var nodeScript, nodeStyle bool

	htmlConverted := ""
	err = nil

	z := html.NewTokenizer(file)

LOOP:
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			err = z.Err()
			if err == io.EOF {
				err = nil
			}
			break LOOP
		case html.CommentToken, html.DoctypeToken:
			continue
		}

		token := z.Token()

		data := strings.ToLower(token.Data)
		switch data {
		case "script":
			if tt == html.StartTagToken {
				nodeScript = true
			} else if tt == html.EndTagToken {
				nodeScript = false
			}
			continue
			//  case "style":
			//	if tt == html.StartTagToken {
			//		nodeStyle = true
			//	} else if tt == html.EndTagToken {
			//		nodeStyle = false
			//	}
			//	continue
		}
		if nodeScript || nodeStyle {
			continue
		}

		for _, attr := range token.Attr {
			key := strings.ToLower(attr.Key)
			if key == "href" || key == "src" {
				l.Info(attr.Val)
			}
		}

		htmlConverted += token.String()
	}

	return filesToPrepare, fileIndex, ioutil.WriteFile(fileToPrepare.fileName, []byte(htmlConverted), 0644)
}

const onClean = "on loaderHTTP.Clean(): "

func (flOp *loaderHTTP) Clean(term *selectors.Term, _ *crud.RemoveOptions) error {
	return nil
}
