package auth_jwt

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/encrlib"
	"github.com/pavlo67/common/common/rbac"
)

var _ auth.Operator = &authJWT{}

type authJWT struct {
	privKey rsa.PrivateKey
	builder jwt.Builder
}

// TODO!!! add expiration time

const onNew = "on auth_jwt.New()"

func New(pathToStore string) (auth.Operator, error) {
	privKey, err := encrlib.NewRSAPrivateKey(pathToStore)
	if err != nil {
		return nil, errors.Wrap(err, onNew)
	}

	signerOpts := (&jose.SignerOptions{}).WithType("Token") // signerOpts.WithType("Token")
	signingKey := jose.SigningKey{Algorithm: jose.RS256, Key: privKey}
	rsaSigner, err := jose.NewSigner(signingKey, signerOpts)
	if err != nil {
		return nil, errors.Wrapf(err, onNew+": can't jose.NewSigner(%#v, %#v)", signingKey, signerOpts)
	}

	return &authJWT{privKey: *privKey, builder: jwt.Signed(rsaSigner)}, nil
}

type JWTCreds struct {
	*jwt.Claims
	Nickname string       `json:",omitempty"`
	GroupID  common.IDStr `json:",omitempty"`

	// couldn't use rbac.Roles type because it has unappropriate .MarshalJSON() method
	Roles rbac.Roles
}

// 	SetCreds ignores all input parameters, creates new "BTC identity" and returns it
func (authOp *authJWT) SetCreds(actor auth.Actor, toSet auth.Creds) (*auth.Creds, error) {

	var userID auth.ID
	if actor.Identity != nil {
		userID = actor.Identity.ID
		if userIDAnother := auth.ID(toSet[auth.CredsID]); userIDAnother != "" && actor.Identity.Roles.Has(rbac.RoleAdmin) {
			userID = userIDAnother
		}
	}

	jc := JWTCreds{
		Claims: &jwt.Claims{
			// Issuer:   "issuer1",
			// Subject:  "subject1",
			// Audience: jwt.Audience{"aud1", "aud2"},
			ID:       string(userID),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Expiry:   jwt.NewNumericDate(time.Date(2017, 1, 1, 0, 8, 0, 0, time.UTC)),
		},

		Nickname: toSet[auth.CredsNickname],
		GroupID:  common.IDStr(toSet[auth.CredsGroupID]),
	}

	if rolesJSON := toSet[auth.CredsRolesJSON]; rolesJSON != "" {
		if err := json.Unmarshal([]byte(rolesJSON), &jc.Roles); err != nil {
			return nil, fmt.Errorf("on authJWT.SetCreds() with json.Unmarshal(%s): %s", rolesJSON, err)
		}
	}

	// add claims to the Builder
	builder := authOp.builder.Claims(jc)

	rawJWT, err := builder.CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("on authJWT.SetCreds() with builder.CompactSerialize(): %s", err)
	}

	delete(toSet, auth.CredsToSet)

	toSet[auth.CredsJWT] = rawJWT

	return &toSet, nil
}

func (authOp *authJWT) Authenticate(toAuth auth.Creds) (*auth.Actor, error) {
	credsJWT := toAuth[auth.CredsJWT]
	if strings.TrimSpace(credsJWT) == "" {
		return nil, nil
	}

	// l.Infof("length = %d: '%s'", len(credsJWT), credsJWT)

	parsedJWT, err := jwt.ParseSigned(credsJWT)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse Token: %s", credsJWT)
	}

	res := JWTCreds{}
	err = parsedJWT.Claims(&authOp.privKey.PublicKey, &res)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get claims: %#v", parsedJWT)
	}

	return &auth.Actor{
		Identity: &auth.Identity{
			ID:       auth.ID(res.ID),
			Nickname: res.Nickname,
			Roles:    res.Roles,
		},
	}, nil
}
