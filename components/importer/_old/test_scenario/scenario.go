package importer_test

import (
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/interfaces/importer"
)

type ImporterTestCase struct {
	Operator importer.Operator
	Fount    string
	DBKey    string
}

func TestImporterWithCases(t *testing.T, testCases []ImporterTestCase) {
	var err error

	for _, tc := range testCases {
		err = tc.Operator.Init(tc.Fount, tc.DBKey, false)
		if err != nil {
			t.Fatalf("can't init importer.Operator (%+v): %v", tc.Operator, err)
		}

		for {
			entity, err := tc.Operator.Next()
			if err == importer.ErrNoMoreItems {
				break
			}
			if err != nil {
				t.Fatalf("can't get next item: %v", err)
			}

			object, err := entity.Object()
			if err != nil {
				t.Fatalf("can't .Object(): %s", err)
			}

			flowItem, err := entity.FlowItem()
			if err != nil {
				t.Fatalf("can't .FlowItem(): %s", err)
			}

			//log.Println("ID:", entity.OriginalID(), "\nItem:", flowItem, "\nObj:", object)
			log.Println("/nID:", entity.OriginalID(), "\nItem:", flowItem.Content)

			if object != nil {
				f, err := os.Create(`obj.html`)
				if err != nil {
					t.Fatalf("can't create file for write: %s", err)
				}
				_, err = f.Write([]byte(object.Content))
				if err != nil {
					t.Fatalf("error write to file: %s", err)
				}
			}

		}
	}
}
