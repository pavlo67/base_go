package auth

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/rbac"
)

func Auth(env config.Envs, authOp Operator, role rbac.Role) (actor *Actor, err error) {

	var actors []Actor

	if err := env.Value("actors", &actors); err != nil {
		return nil, err
	}

	for _, actor := range actors {
		if actor.Identity != nil && actor.Identity.Roles.Has(role) {
			actorAuthenticated, err := authOp.Authenticate(actor.Creds)

			if err != nil {
				return nil, err
			}
			return actorAuthenticated, nil
		}
	}

	return nil, errors.Wrapf(common.ErrNotFound, fmt.Sprintf("actor with role %s isn't found", role))
}
