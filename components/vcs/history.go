package vcs

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pavlo67/data/components/ns"
)

type ActionKey string

const ProducedAction ActionKey = "produced_from"
const SavedAction ActionKey = "saved"
const CreatedAction ActionKey = "created"
const UpdatedAction ActionKey = "updated"

type Action struct {
	Actor  ns.URN    `bson:",omitempty" json:",omitempty"`
	Key    ActionKey `bson:",omitempty" json:",omitempty"`
	DoneAt time.Time `bson:",omitempty" json:",omitempty"`
	Error  error     `bson:",omitempty" json:",omitempty"`
}

type History []Action

func (h History) FirstByKey(key ActionKey) int {
	for i, action := range h {
		if action.Key == key {
			return i
		}
	}

	return -1
}

func (h History) SaveAction(action Action) History {
	i := h.FirstByKey(action.Key)
	if i >= 0 {
		h[i] = action
	} else {
		h = append(h, action)
	}
	return h
}

func (h History) CheckOn(hOld History) error {
	if len(hOld) < 1 {
		return nil
	}

	actionLast := hOld[len(hOld)-1]
	for _, actionNew := range h {
		if reflect.DeepEqual(actionLast, actionNew) {
			return nil
		}
	}

	return fmt.Errorf("history (%#v) is inappropriate to the old one (... %#v)", h, actionLast)
}
