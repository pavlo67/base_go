package auth

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/rbac"
)

func Auth(cfgService config.Config, authOp auth.Operator, role rbac.Role) (actor *auth.Actor, err error) {

	var actors []auth.Actor

	if err := cfgService.Value("actors", &actors); err != nil {
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
