package persons_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/entities/persons"
)

var _ persons.Operator = &personsStub{}

type personsStub struct {
	personItems []persons.Item01
}

const onNew = "on persons_stub.New(): "

func New(personItems []persons.Item01) (persons.Operator, db.Cleaner, error) {
	personsStub := personsStub{
		personItems: personItems,
	}

	return &personsStub, &personsStub, nil
}

const onSave = "on personsStub.Save(): "

var currentID int

func (personsOp *personsStub) Save(personsItem persons.Item01, _ *auth.Identity) (persons.ID, error) {
	if personsItem.ID == nil {
		currentID++
		personsItem.ID = currentID
		personsOp.personItems = append(personsOp.personItems, personsItem)
		return currentID, nil
	}

	for i, pi := range personsOp.personItems {
		if pi.ID == personsItem.ID {
			personsOp.personItems[i] = personsItem
			return pi.ID, nil
		}
	}

	return "", errors.Wrapf(common.ErrNotFound, onSave+"no person with the same ID as %#v", personsItem)
}

const onRead = "on personsStub.Read(): "

func (personsOp *personsStub) Read(id persons.ID, _ *auth.Identity) (*persons.Item01, error) {
	for _, pi := range personsOp.personItems {
		if pi.ID == id {
			return &pi, nil
		}
	}

	return nil, errors.Wrapf(common.ErrNotFound, onSave+"no person with the ID = %#v", id)
}

const onRemove = "on personsStub.Remove(): "

func (personsOp *personsStub) Remove(id persons.ID, _ *auth.Identity) error {
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

func (personsOp *personsStub) List(*selectors.Term, *auth.Identity) ([]persons.Item01, error) {
	return personsOp.personItems, nil
}

var _ db.Cleaner = &personsStub{}

func (personsOp *personsStub) Clean() error {
	personsOp.personItems = nil

	return nil
}
