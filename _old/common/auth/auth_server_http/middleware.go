package auth_server_http

import (
	"net/http"
	"regexp"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server_http"
)

var _ server_http.OnRequestMiddleware = &onRequestMiddleware{}

func OnRequestMiddleware(authJWTOp auth.Operator) (server_http.OnRequestMiddleware, error) {
	if authJWTOp == nil {
		return nil, errors.New("no authJWTOp")
	}

	return &onRequestMiddleware{
		authJWTOp: authJWTOp,
	}, nil
}

type onRequestMiddleware struct {
	authJWTOp auth.Operator
}

var reBearer = regexp.MustCompile(`^\s*Bearer(\s|%[fF]20)*`)

const onOptions = "on onRequestMiddleware.Identity()"

func (orm *onRequestMiddleware) Identity(r *http.Request) (*auth.Identity, error) {
	//if r == nil {
	//	return nil, errors.New("no server_http.Request in RequestOptions(...)")
	//}

	if tokenJWT := r.Header.Get("Authorization"); tokenJWT != "" {
		tokenJWT = reBearer.ReplaceAllString(tokenJWT, "")
		actor, err := orm.authJWTOp.Authenticate(auth.Creds{auth.CredsJWT: tokenJWT})
		if err != nil {
			return nil, errors.CommonError(err, onOptions)
		}
		if actor != nil {
			return actor.Identity, nil
		}
	}

	return nil, nil
}

//// SIGNATURE CHECK
//signature := r.Header.Get("Signature")
//if signature != "" && r.URL != nil {
//	publicKeyAddress := r.Header.Get("Public-Key-Address")
//	numberToSignature := r.Header.Get("Number-To-Signature")
//
//	credsSignature := auth.Creds{
//		Values: map[auth.CredsType]string{
//			auth.CredsPublicKeyBase58:    publicKeyAddress,
//			auth.CredsContentToSignature: r.URL.Path + "?" + r.URL.RawQuery,
//			auth.CredsKeyToSignature:     numberToSignature,
//			auth.CredsSignature:          signature,
//		},
//	}
//
//	user, errs = auth.GetIdentity(credsSignature, authOps, errs)
//	// previous errs is added by auth.GetIdentity()
//}
//var errNoIdentityOpsMap = errors.New("no map[CredsType]identity.UserKey")
