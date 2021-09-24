package types

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/data/elements/ns"
	"github.com/pavlo67/data/elements/vcs"
)

type RelationKey string

type Relation01 struct {
	Key RelationKey
	ns.NSS
	Note string
}

type Relations01Map map[string]Relation01

type Description01 struct {
	URN          ns.URN         `json:",omitempty" bson:",omitempty"`
	Tags         []string       `json:",omitempty" bson:",omitempty"`
	RelationsMap Relations01Map `json:",omitempty" bson:",omitempty"`
	OwnerNSS     ns.NSS         `json:",omitempty" bson:",omitempty"`
	ViewerNSS    ns.NSS         `json:",omitempty" bson:",omitempty"`
	History      vcs.History    `json:",omitempty" bson:",omitempty"`
	CreatedAt    time.Time      `json:",omitempty" bson:",omitempty"`
	UpdatedAt    *time.Time     `json:",omitempty" bson:",omitempty"`
}

var Description01FieldsToSave = []string{"urn", "tags", "relations_map", "owner_nss", "viewer_nss", "history"}
var Description01FieldsToRead = append(Description01FieldsToSave, "created_at", "updated_at")

func (descr *Description01) FoldToSave() ([]interface{}, error) {
	if descr == nil {
		return nil, errors.New("nil persons.Item to be folded")
	}

	var urnBytes, tagsBytes, relationsMapBytes, historyBytes []byte
	var err error

	if len(descr.URN) > 0 {
		urnBytes = []byte(descr.URN)
	}
	// TODO!!! append to descr.History

	// tagsBytes = []byte{} // to satisfy NOT NULL constraint
	if len(descr.Tags) > 0 {
		if tagsBytes, err = json.Marshal(descr.Tags); err != nil {
			return nil, errors.Wrapf(err, "can't marshal .Tags (%#v)", descr.Tags)
		}
	}

	// relationsMapBytes = []byte{} // to satisfy NOT NULL constraint
	if len(descr.RelationsMap) > 0 {
		if relationsMapBytes, err = json.Marshal(descr.RelationsMap); err != nil {
			return nil, errors.Wrapf(err, "can't marshal .RelationsMap (%#v)", descr.RelationsMap)
		}
	}

	// historyBytes = []byte{} // to to satisfy NOT NULL constraint
	if len(descr.History) > 0 {
		historyBytes, err = json.Marshal(descr.History)
		if err != nil {
			return nil, errors.Wrapf(err, "can't marshal .History(%#v)", descr.History)
		}
	}

	return []interface{}{urnBytes, tagsBytes, relationsMapBytes, descr.OwnerNSS, descr.ViewerNSS, historyBytes}, nil
}

func (descr *Description01) UnfoldReaded(urnBytes, tagsBytes, relationsMapBytes, historyBytes []byte) error {
	if descr == nil {
		return errors.New("nil Description01 to be unfolded")
	}

	descr.URN = ns.URN(urnBytes)
	// TODO!!! append to descr.History

	if len(tagsBytes) > 0 {
		if err := json.Unmarshal(tagsBytes, &descr.Tags); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Tags (%s)", tagsBytes)
		}
	}

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
