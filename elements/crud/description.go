package crud

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/pavlo67/data/elements/ns"
	"github.com/pavlo67/data/elements/vcs"
)

type Description struct {
	URN       ns.URN      `json:",omitempty" bson:",omitempty"`
	Tags      []string    `json:",omitempty" bson:",omitempty"`
	OwnerNSS  ns.NSS      `json:",omitempty" bson:",omitempty"`
	ViewerNSS ns.NSS      `json:",omitempty" bson:",omitempty"`
	History   vcs.History `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time   `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time  `json:",omitempty" bson:",omitempty"`
}

func (descr *Description) UnfoldFromJSON(tagsBytes, urnBytes, historyBytes []byte) error {
	if descr == nil {
		return errors.New("nil Description to be unfolded")
	}

	if len(tagsBytes) > 0 {
		if err := json.Unmarshal(tagsBytes, &descr.Tags); err != nil {
			return errors.Wrapf(err, "can't unmarshal .Tags (%s)", tagsBytes)
		}
	}

	descr.URN = ns.URN(urnBytes)

	// TODO!!! append to descr.History

	if len(historyBytes) > 0 {
		if err := json.Unmarshal(historyBytes, &descr.History); err != nil {
			return errors.Wrapf(err, "can't unmarshal .History (%s)", historyBytes)
		}
	}

	return nil
}

func (descr *Description) FoldIntoJSON() (tagsBytes, urnBytes, historyBytes []byte, err error) {
	if descr == nil {
		return nil, nil, nil, errors.New("nil persons.Item to be folded")
	}

	tagsBytes = []byte{} // to satisfy NOT NULL constraint
	if len(descr.Tags) > 0 {
		if tagsBytes, err = json.Marshal(descr.Tags); err != nil {
			return nil, nil, nil, errors.Wrapf(err, "can't marshal .Tags (%#v)", descr.Tags)
		}
	}

	if len(descr.URN) > 0 {
		urnBytes = []byte(descr.URN)
	}

	// TODO!!! append to descr.History

	historyBytes = []byte{} // to to satisfy NOT NULL constraint
	if len(descr.History) > 0 {
		historyBytes, err = json.Marshal(descr.History)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "can't marshal .History(%#v)", descr.History)
		}
	}

	return tagsBytes, urnBytes, historyBytes, nil
}
