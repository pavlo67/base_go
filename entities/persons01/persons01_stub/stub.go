package persons01_stub

import (
	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/selectors"
	"github.com/pavlo67/data/components/vcs"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/db"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/entities/persons01"
)

var _ persons01.Operator = &personsStub{}

type personsStub struct {
	personItems []persons01.Item
}

const onNew = "on persons_stub.New(): "

func New(personItems []persons01.Item) (persons01.Operator, db.Cleaner, error) {
	personsStub := personsStub{
		personItems: personItems,
	}

	return &personsStub, &personsStub, nil
}

const onSave = "on personsStub.Save(): "

var currentID int

func (personsOp *personsStub) Save(personsItem persons01.Item, _ auth.Actor) (persons01.ID, vcs.History, error) {
	if personsItem.ID == "" {
		currentID++
		personsItem.ID = crud.NewIDInt64(int64(currentID))

		personsOp.personItems = append(personsOp.personItems, personsItem)
		return personsItem.ID, nil, nil
	}

	for i, pi := range personsOp.personItems {
		if pi.ID == personsItem.ID {
			urnOriginal := pi.Description.URN
			personsOp.personItems[i] = personsItem
			personsOp.personItems[i].Description.URN = urnOriginal // UNCHANGED!!!
			return pi.ID, nil, nil
		}
	}

	return "", nil, errors.Wrapf(common.ErrNotFound, onSave+"no person with the same ID as %#v", personsItem)
}

const onRead = "on personsStub.Read(): "

func (personsOp *personsStub) Read(id persons01.ID, _ auth.Actor) (*persons01.Item, error) {
	for _, pi := range personsOp.personItems {
		if pi.ID == id {
			return &pi, nil
		}
	}

	return nil, errors.Wrapf(common.ErrNotFound, onSave+"no person with the ID = %#v", id)
}

const onRemove = "on personsStub.Remove(): "

func (personsOp *personsStub) Remove(id persons01.ID, _ auth.Actor) error {
	for i, pi := range personsOp.personItems {
		if pi.ID == id {
			personItemsOld := personsOp.personItems

			personsOp.personItems = personsOp.personItems[:i]
			if i <= len(personItemsOld)-2 {
				personsOp.personItems = append(personsOp.personItems, personItemsOld[1+1:]...)
			}

			return nil
		}
	}

	return errors.Wrapf(common.ErrNotFound, onSave+"no person with the ID = %#v", id)
}

const onList = "on personsStub.List(): "

func (personsOp *personsStub) List(*selectors.Term, auth.Actor) ([]persons01.Item, error) {
	return personsOp.personItems, nil
}

var _ db.Cleaner = &personsStub{}

func (personsOp *personsStub) Clean() error {
	personsOp.personItems = nil

	return nil
}
