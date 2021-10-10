package auth

import "github.com/pavlo67/common/common/auth"

type Actor struct {
	*auth.Identity
	auth.Creds
}
