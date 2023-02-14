package importer_old

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis/program"

	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/confidenter"
	"github.com/pavlo67/punctum/interfaces/convertor"
	"github.com/pavlo67/punctum/interfaces/flow"
	"github.com/pavlo67/punctum/interfaces/objects"
)

const InterfaceKey = "importer"

type ImportType string

const RSSType ImportType = "rss"

var ErrNoFount = errors.New("no fount is reachable")
var ErrNoMoreItems = errors.New("no more items")
var ErrBadItemID = errors.New("bad item id")
var ErrBadItem = errors.New("bad item")
var ErrNilItem = errors.New("item is nil")

type TmpObject struct {
	ID       int64
	Author   string
	Title    string
	Content  string
	FileName string
}

type Authors struct {
	Name  string
	Count int64
}

//// Deprecated: changed to EntityNew!!!
//type Entity interface {
//	// OriginalID gets the value from the imported entity
//	OriginalID() string
//
//	// Object forms an interfaces.Object from the imported entity
//	Object() (obj *interfaces.Object, err error)
//
//	// FlowItem forms an flow.Item from the imported entity
//	FlowItem() (obj *flow.Item, err error)
//}

type Entity interface {
	// OriginalID gets the value from the imported entity
	OriginalID() string

	convertor.Operator
}

// Operator is general interface for external data import
type Operator interface {

	// Init opens import session with selected data fount
	Init(fount, dbKey string, testMode bool) error

	// Next gets the next data entity from the fount
	Next() (entity Entity, err error)

	Close()
}

func Task(userIdentity *confidenter.Identity, importerOp Operator, objectsOp objects.Operator, flowOp flow.Operator, key, source, importTo, dbKey, fountID, login, tags string, rView confidenter.IdentityString, testMode bool) (int, error) {

	var countNew int

	err := importerOp.Init(source, dbKey, testMode)
	if err != nil {
		return 0, errors.Wrap(err, "can't init "+key+" importer")
	}
	liImets := 0
	for {
		liImets++
		if liImets%100 == 0 {
			fmt.Println(liImets)
		}
		entity, err := importerOp.Next()
		if err == ErrNoMoreItems {
			break
		}
		if err != nil {
			log.Printf("error reading imported item: %s", err)
			continue
		}
		if importTo == "flow" {
			item, err := entity.FlowItem()
			if err != nil || item == nil {
				log.Printf("error converting rss entity (%#v) into Flow.Item: %s", entity, err)
				continue
			}
			if fountID != "" {
				item.FountIS = confidenter.IdentityString(program.Domain() + "/fount/" + fountID)
			}
			if rView != "" {
				item.RView = rView
			} else {
				item.RView = userIdentity.String()
			}
			item.ROwner = userIdentity.String()
			ok, err := flowOp.IsNew(*item)
			if err != nil {
				log.Printf("error checking flowOp.IsNew(%s): %s", item.OriginalID, err)
			}
			if !ok {
				continue
			}
			runes := []rune(item.Summary)
			if len(runes) > 255 {
				item.Summary = string(runes[0:250]) + "..."
			}

			_, err = flowOp.Create(userIdentity, *item)
			if err != nil {
				log.Printf("error flowOp.Create(%v): %s", item, err)
			}
			countNew++
		} else {
			o, err := entity.Object()
			if err != nil {
				log.Printf("can't get Object() for imported row: %v; %v", entity, err)
				continue
			}
			o.Genus = key
			if o.Author == "" {
				o.Author = source
			}
			o.ROwner = userIdentity.String()
			o.RView = o.ROwner
			o.Visibility = interfaces.Private

			objects.AddTags(userIdentity, o, tags)

			//if login != "" {
			//	// write page content to file
			//	var path string
			//	path, _, err = filer.pathToFileOrDir(userIdentity, pathRepository, "")
			//	filename := validator.RandomString(12) + `.html`
			//	fullpath := path + `/` + filename
			//	err := ioutil.WriteFile(fullpath, []byte(o.Content), 0644)
			//	if err != nil {
			//		log.Printf("error writing imported data as a file (%s): %s", fullpath, err)
			//	} else {
			//		objects.AddFile(userIdentity, o, filename, fullpath, fileInfoOp)
			//		//o.Content = ""
			//	}
			//}

			_, err = objectsOp.Create(userIdentity, o)
			if err != nil {
				me, ok := errors.Cause(err).(*mysql.MySQLError)
				if !ok || me.Number != 1062 {
					log.Printf("can't create import object for imported object: %v; %v", o, err)
				}

				continue
			}
			countNew++
		}
	}
	return countNew, nil
}
