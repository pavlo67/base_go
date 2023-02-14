package entities

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/vcs"
)

type RelationKey string

type Relation struct {
	Key RelationKey
	ns.NSS
	Note string
}

type RelationsMap map[string]Relation

type Description struct {
	Tags         []string     `json:",omitempty" bson:",omitempty"`
	RelationsMap RelationsMap `json:",omitempty" bson:",omitempty"`
	OwnerNSS     ns.NSS       `json:",omitempty" bson:",omitempty"`
	ViewerNSS    ns.NSS       `json:",omitempty" bson:",omitempty"`
}

type Point struct {
	URN       ns.URN      `json:",omitempty" bson:",omitempty"`
	History   vcs.History `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time   `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time  `json:",omitempty" bson:",omitempty"`
}

var Description01FieldsBasis = []string{"urn", "tags", "relations_map", "owner_nss", "viewer_nss", "history"}
var Description01FieldsToUpdate = append(Description01FieldsBasis, "updated_at")
var Description01FieldsToRead = append(Description01FieldsBasis, "updated_at", "created_at")

func (descr *Description) FoldToSavePg(onInsert bool) ([]interface{}, vcs.History, string, error) {
	if descr == nil {
		return nil, nil, "", errors.New("on FoldToSavePg(): nil persons.Item to be folded")
	}

	var relationsMapBytes []byte
	var err error

	// relationsMapBytes = []byte{} // to satisfy NOT NULL constraint
	if len(descr.RelationsMap) > 0 {
		if relationsMapBytes, err = json.Marshal(descr.RelationsMap); err != nil {
			return nil, nil, "", errors.Wrapf(err, "on FoldToSavePg(): can't marshal .RelationsMap (%#v)", descr.RelationsMap)
		}
	}

	var urnBytes []byte
	if len(descr.URN) > 0 {
		urnBytes = []byte(descr.URN)
	}

	if onInsert {
		return []interface{}{urnBytes, pq.Array(descr.Tags), relationsMapBytes, descr.OwnerNSS, descr.ViewerNSS, ""}, nil, "", nil
	}

	historyChanged, historyChangedStr, historyOriginalStr, err := vcs.ModifyHistory(descr.History, nil)
	if err != nil {
		return nil, nil, "", errors.Wrap(err, "on FoldToSavePg()")
	}

	return []interface{}{urnBytes, pq.Array(descr.Tags), relationsMapBytes, descr.OwnerNSS, descr.ViewerNSS, historyChangedStr, time.Now().UTC()}, historyChanged, historyOriginalStr, nil

}

func (descr *Description) UnfoldReaded(urnBytes, relationsMapBytes, historyBytes []byte) error {
	if descr == nil {
		return errors.New("nil Description to be unfolded")
	}

	descr.URN = ns.URN(urnBytes)
	// TODO!!! append to descr.History

	if len(relationsMapBytes) > 0 {
		if err := json.Unmarshal(relationsMapBytes, &descr.RelationsMap); err != nil {
			return errors.Wrapf(err, "can't unmarshal .RelationsMap (%s)", relationsMapBytes)
		}
	}

	if len(historyBytes) > 0 {
		if err := json.Unmarshal(historyBytes, &descr.History); err != nil {
			return errors.Wrapf(err, "can't unmarshal .History (%s)", historyBytes)
		}
	}

	return nil
}

func (testDescription Description) TestIfEqual(t *testing.T, descriptionToCheck Description) {
	require.Equal(t, testDescription.URN, descriptionToCheck.URN)

	if len(testDescription.Tags) > 0 {
		require.Equal(t, testDescription.Tags, descriptionToCheck.Tags)
	} else {
		require.Equal(t, 0, len(descriptionToCheck.Tags))
	}
	if len(testDescription.RelationsMap) > 0 {
		require.Equal(t, testDescription.RelationsMap, descriptionToCheck.RelationsMap)
	} else {
		require.Equal(t, 0, len(descriptionToCheck.RelationsMap))
	}

	require.Equal(t, testDescription.ViewerNSS, descriptionToCheck.ViewerNSS)
	require.Equal(t, testDescription.OwnerNSS, descriptionToCheck.OwnerNSS)

	require.True(t, len(descriptionToCheck.History) >= len(testDescription.History))
	if len(testDescription.History) > 0 {
		require.Equal(t, testDescription.History, descriptionToCheck.History[:len(testDescription.History)])
	}
}

func (testDescription Description) ChangeForTest() Description {
	testDescription.URN += "_changed"
	testDescription.Tags = append(testDescription.Tags, "changed_tag")
	if testDescription.RelationsMap == nil {
		testDescription.RelationsMap = RelationsMap{}
	}
	testDescription.RelationsMap["changed"] = Relation{
		Key:  "chg",
		NSS:  "qwer",
		Note: "wqer qwer",
	}
	testDescription.OwnerNSS += "_changed"
	testDescription.ViewerNSS += "_changed"

	return testDescription
}

var TestDescription = Description{
	URN:  "urn1",
	Tags: []string{"famous", "writer"},
	RelationsMap: RelationsMap{"r": Relation{
		Key:  "r1key",
		NSS:  "nss_r1",
		Note: "wetr wert eryry",
	}},
	OwnerNSS:  "owner_nss",
	ViewerNSS: "viever_nss",
	// History:      nil,
}
