package auth_stub

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
)

var _ auth.Operator = &authstub{}

type authstub struct {
	crypter crypt.Crypter
	actors  []auth.Actor
}

const onNew = "on authstub.New()"

func New(defaultActors []auth.Actor) (auth.Operator, error) {

	authOp := authstub{crypter: crypt.SHA256.New()}

	for _, actor := range defaultActors {
		if actor.Creds == nil {
			actor.Creds = auth.Creds{}
		}

		var err error
		actor.Creds[auth.CredsPasshash], err = authOp.crypter.Generate([]byte(actor.Creds[auth.CredsPassword]), nil)
		if err != nil {
			return nil, errors.Wrap(err, onNew)
		}
		delete(actor.Creds, auth.CredsPassword)

		authOp.actors = append(authOp.actors, actor)
	}

	return &authOp, nil
}

const onSetCreds = "on authstub.SetCreds()"

func (authOp *authstub) SetCreds(actor auth.Actor, toSet auth.Creds) (*auth.Creds, error) {
	if passwordToSet := strings.TrimSpace(toSet[auth.CredsPassword]); passwordToSet != "" {
		var err error
		toSet[auth.CredsPasshash], err = authOp.crypter.Generate([]byte(passwordToSet), nil)
		if err != nil {
			return nil, errors.Wrap(err, onSetCreds)
		}
		delete(actor.Creds, auth.CredsPassword)

	}
	delete(toSet, auth.CredsPassword)

	var idToFind auth.ID

	if actor.Identity != nil {
		idToFind = actor.Identity.ID
		if idToFindAnother := auth.ID(toSet[auth.CredsID]); idToFindAnother != "" {
			if actor.Identity.Roles.Has(rbac.RoleAdmin) {
				idToFind = idToFindAnother
			} else {
				return nil, fmt.Errorf(onSetCreds+": actor (H%#v) can't set alien creds non having admin role", actor.Identity)
			}
		}
	}

	i := -1

	if idToFind == "" {
		i = len(authOp.actors)
		actor := auth.Actor{Identity: &auth.Identity{ID: auth.ID(strconv.Itoa(i))}}
		authOp.actors = append(authOp.actors, actor)
	} else {
		for iToCheck, actorToSet := range authOp.actors {
			if actorToSet.Identity != nil && actorToSet.Identity.ID == idToFind {
				i = iToCheck
			}
		}
		if i < 0 {
			return nil, errors.Wrapf(auth.ErrNoUser, "no user with ID = %s", idToFind)
		}
	}

	actorToSet := authOp.actors[i]

	if nicknameToSet := strings.TrimSpace(toSet[auth.CredsNickname]); nicknameToSet != "" {
		actorToSet.Identity.Nickname = nicknameToSet
		delete(toSet, auth.CredsNickname)
	}
	if roleStr, ok := toSet[auth.CredsRole]; ok {
		role := rbac.Role(roleStr)
		if role == rbac.RoleAdmin && !actor.Roles.Has(rbac.RoleAdmin) {
			return nil, fmt.Errorf(onSetCreds+": actor (H%#v) can't set admin role non having own admin role", actor.Identity)
		}

		// TODO: multiple roles
		actorToSet.Identity.Roles = rbac.Roles{rbac.Role(role)}
	}

	if actorToSet.Creds == nil {
		actorToSet.Creds = auth.Creds{}
	}
	for k, v := range toSet {
		actorToSet.Creds[k] = v
	}

	authOp.actors[i] = actorToSet

	return &auth.Creds{auth.CredsNickname: actorToSet.Identity.Nickname}, nil // auth.CredsRole: actorToSet.Identity.Roles
}

const onAuthenticate = "on authstub.Authenticate()"

func (authOp *authstub) Authenticate(toAuth auth.Creds) (*auth.Actor, error) {
	nickname := toAuth[auth.CredsNickname]

	// l.Infof("ACTORS: %#v", authOp.actors)

	l.Infof("TO AUTH: %s / %#v", nickname, toAuth)

	for _, actor := range authOp.actors {
		if actor.Identity != nil && actor.Identity.Nickname == nickname {
			if err := authOp.crypter.Verify(actor.Creds[auth.CredsPasshash], []byte(toAuth[auth.CredsPassword])); err == nil {
				return &actor, nil
			}
		}
	}

	return nil, auth.ErrNotAuthenticated
}
