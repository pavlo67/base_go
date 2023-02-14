package auth_persons

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/data/components/selectors"

	"github.com/pavlo67/data/entities/persons"
)

var _ auth.Operator = &authPersons{}

type authPersons struct {
	personsOp             persons.Operator
	maxPersonsToAuthCheck int
}

const onNew = "on authPersons.New()"

func New(personsOp persons.Operator, maxPersonsToAuthCheck int) (auth.Operator, error) {
	if personsOp == nil {
		return nil, errors.New(onNew + ": no persons.Operator")
	}
	if maxPersonsToAuthCheck < 1 {
		maxPersonsToAuthCheck = 1
	}

	return &authPersons{
		personsOp:             personsOp,
		maxPersonsToAuthCheck: maxPersonsToAuthCheck,
	}, nil
}

const onSetCreds = "on authPersons.SetCredsByKey()"

func (authOp *authPersons) SetCreds(authID auth.ID, toSet auth.Creds) (*auth.Creds, error) {
	if authID == "" {
		// TODO: set .Allowed = false and verify email

		identity := auth.Identity{
			Nickname: toSet[auth.CredsNickname],
		}

		// TODO!!! hash password

		person := persons.Item{Identity: identity}
		if err := person.SetCreds(toSet); err != nil {
			return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.SetCreds(%#v)", toSet)
		}

		if _, err := authOp.personsOp.Save(person, auth.IdentityWithRoles(rbac.RoleAdmin)); err != nil {
			return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.Add(%#v, nil)", identity)
		}

		return &toSet, nil
	}

	//if item.Creds, err = hashCreds(item.Creds, itemOld.Creds); err != nil {
	//	return nil, errata.CommonError(err, onChange)
	//}

	return nil, common.ErrNotImplemented

	//credsTypeToSet := auth.CredsType(toSet[auth.CredsToSet])
	//delete(toSet, auth.CredsToSet)
	//
	//credsToSet, ok := toSet[credsTypeToSet]
	//if !ok {
	//	return nil, errors.Errorf(onSetCreds+"no creds to set in %#v", toSet)
	//}
	//
	//selector := selectors.Binary(selectors.Eq, persons.UserKeyFieldName, selectors.Value{string(personKey)})
	//
	//items, err := authOp.personsOp.List(selector, nil)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.List(selector = %#v, nil)", *selector)
	//}
	//if len(items) < 1 {
	//	return nil, errors.Errorf(onSetCreds+"no persons with key %s)", personKey)
	//} else if len(items) > 1 {
	//	return nil, errors.Errorf(onSetCreds+"too many persons with key %s)", personKey)
	//}
	//
	//if credsTypeToSet == auth.CredsEmail {
	//	// TODO: verify and/or another actions with some other creds types
	//}
	//
	//items[0].Creds[credsTypeToSet] = credsToSet
	//
	//_, err = authOp.personsOp.Add(items[0], nil)
	//if err != nil {
	//	return nil, errors.Wrapf(err, onSetCreds+"can't .personsOp.Add(%#v, nil)", items[0])
	//
	//}
	//
	//return &items[0].Creds, nil
}

const onAuthenticate = "on authPersons.Authenticate()"

var reEmail = regexp.MustCompile("@")

func (authOp *authPersons) Authenticate(toAuth auth.Creds) (*auth.Identity, error) {
	nickname := strings.TrimSpace(toAuth[auth.CredsNickname])
	if nickname == "" {
		return nil, errors.CommonError(common.NoCredsKey, common.Map{"no nickname in creds": toAuth})
	}

	password := strings.TrimSpace(toAuth[auth.CredsPassword])

	//if login := toAuth.StringDefault(auth.CredsLogin, ""); login != "" {
	//	if reEmail.MatchString(login) {
	//		selector = selectors.Binary(selectors.Eq, persons.EmailFieldName, selectors.Value{login})
	//	} else {
	//		selector = selectors.Binary(selectors.Eq, persons.NicknameFieldName, selectors.Value{login})
	//	}
	//} else if email, ok := toAuth[auth.CredsEmail]; ok {
	//	selector = selectors.Binary(selectors.Eq, persons.EmailFieldName, selectors.Value{email})
	//} else if nickname, ok := toAuth[auth.CredsNickname]; ok {
	//	selector = selectors.Binary(selectors.Eq, persons.NicknameFieldName, selectors.Value{nickname})
	//} else {
	//	return nil, nil
	//	// return nil, errors.New(onAuthorize + "no login to auth")
	//}
	//selector = logic.AND(
	//	selector,
	//	selectors.Binary(selectors.Gt, persons.VerifiedFieldName, selectors.Value{0}),
	//)

	var selector selectors.Item
	//selector, err := authOp.personsOp.HasNickname(nickname)
	//if err != nil {
	//	return nil, errata.CommonError(err, onAuthenticate)
	//}

	items, err := authOp.personsOp.List(&selector, auth.IdentityWithRoles(rbac.RoleAdmin))
	if err != nil {
		return nil, errors.CommonError(err, fmt.Sprintf(onAuthenticate+": can't .personsOp.List(selector = %#v, nil)", selector))
	}

	for _, item := range items {

		// TODO!!! remove this check adding selector above
		if item.Nickname != nickname {
			continue
		}

		if item.CheckCreds(auth.CredsPassword, password) {
			return &item.Identity, nil
		}

	}

	return nil, errors.CommonError(common.NoCredsKey, common.Map{onAuthenticate + ": wrong creds": toAuth})

	//maxPersonsToAuthCheck := authOp.maxPersonsToAuthCheck
	//if len(items) < authOp.maxPersonsToAuthCheck {
	//	maxPersonsToAuthCheck = len(items)
	//}
	//for i := 0; i < maxPersonsToAuthCheck; i++ {
	//
	//	//// TODO: use selector.AND (commented at the moment)
	//	//if !items[i].Allowed {
	//	//	continue
	//	//}
	//
	//	item := items[i]
	//
	//	if authOp.personsOp.CheckPassword(toAuth[auth.CredsPassword], item.Creds[auth.CredsPasshash]) {
	//		person := item.User
	//		person.Creds = auth.Creds{
	//			auth.CredsNickname: item.Creds[auth.CredsNickname],
	//		}
	//
	//		return &person, nil
	//	}
	//}

}
