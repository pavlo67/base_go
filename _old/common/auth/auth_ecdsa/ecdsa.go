package auth_ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	r "math/rand"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/encrlib"
)

// const Cryptype encrlib.Cryptype = "ecdsa"

const Proto = "ecdsa"

var _ auth.Operator = &authECDSA{}

var errIPToCheckSignature = errors.New("wrong IP to check signature")
var errWrongSignature = errors.New("wrong signature")
var errWrongNumber = errors.New("wrong user's number")
var errEmptyPublicKeyAddress = errors.New("empty public key address")
var errEmptyPrivateKeyGenerated = errors.New("empty private key generated")

//var cnt uint32
//type Session struct {
//	IP        string
//	StartedAt time.Time
//}

type authECDSA struct {
	//sessions map[uint64]Session
	//mutex    *sync.Mutex
	//maxSessionDuration time.Duration
	//numbersLimit       int
}

func New() (auth.Operator, error) {
	r.Seed(time.Now().UnixNano())

	is := &authECDSA{
		//sessions: map[uint64]Session{},
		//mutex:    &sync.Mutex{},
		//maxSessionDuration: maxSessionDuration,
		//numbersLimit:       numbersLimit,
	}

	return is, nil
}

// 	SetCreds creates either session-generated key or new "BTC identity" and returns it
func (is *authECDSA) SetCreds(actor auth.Actor, toSet auth.Creds) (*auth.Creds, error) {

	//toSet := auth.CredsType(toSet[auth.CredsToSet])
	//
	//if toSet == auth.CredsKeyToSignature {
	//	now := time.Now()
	//
	//	is.mutex.Lock() // Lock() -----------------------------------------------------
	//
	//	if is.numbersLimit > 0 && len(is.sessions) >= is.numbersLimit {
	//		for n, s := range is.sessions {
	//			if now.Sub(s.StartedAt) >= is.maxSessionDuration {
	//				delete(is.sessions, n)
	//			}
	//		}
	//	}
	//
	//	cnt++
	//	numberToSend := uint64(cnt)<<32 + uint64(r.Uint32())
	//
	//	is.sessions[numberToSend] = Session{
	//		IP:        creds[auth.CredsIP], // TODO??? check if IP isn't empty
	//		StartedAt: now,
	//	}
	//
	//	is.mutex.Unlock() // Unlock() -------------------------------------------------
	//
	//	return &auth.Creds{auth.CredsKeyToSignature: strconv.FormatUint(numberToSend, 10)}, nil
	//}

	// TODO: modify acceptableIDs if it's necessary

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	} else if privKey == nil {
		return nil, errEmptyPrivateKeyGenerated
	}

	privKeyBytes, err := encrlib.ECDSASerialize(*privKey)
	if err != nil {
		return nil, err
	}

	publKeyBase58 := base58.Encode(encrlib.ECDSAPublicKey(*privKey))
	//nickname := publKeyBase58
	//if toSet.StringDefault(auth.CredsNickname, "") != "" {
	//	nickname = toSet[auth.CredsNickname]
	//}

	credsNew := &auth.Creds{
		auth.CredsNickname:          publKeyBase58, //  nickname,
		auth.CredsPrivateKey:        string(privKeyBytes),
		auth.CredsPublicKeyBase58:   publKeyBase58,
		auth.CredsPublicKeyEncoding: Proto,
	}

	return credsNew, nil
}

const onAuthorize = "on authECDSA.Authenticate(): "

func (is *authECDSA) Authenticate(toAuth auth.Creds) (*auth.Actor, error) {
	if toAuth[auth.CredsPublicKeyEncoding] != Proto {
		return nil, errors.Wrap(auth.ErrEncryptionType, onAuthorize)
	}

	keyToSignature := strings.TrimSpace(toAuth[auth.CredsIP])
	if keyToSignature == "" {
		return nil, errors.Wrapf(errIPToCheckSignature, "toAuth: %#v", toAuth)
	}

	publKeyBase58 := toAuth[auth.CredsPublicKeyBase58]
	if len(publKeyBase58) < 1 {
		return nil, errEmptyPublicKeyAddress
	}
	publKey := base58.Decode(publKeyBase58)

	//numberToSend, err := strconv.ParseUint(keyToSignature, 10, 64)
	//if err != nil {
	//	return nil, errors.Wrap(auth.ErrSignaturedKey, "not a number!")
	//}
	//is.mutex.Lock() // Lock() -----------------------------------------------------
	//
	//session, ok := is.sessions[numberToSend]
	//if !ok {
	//	return nil, errors.Wrap(auth.ErrSignaturedKey, "no appropriate session")
	//	is.mutex.Unlock() // Unlock() ---------------------------------------------
	//}
	//delete(is.sessions, numberToSend)
	//
	//is.mutex.Unlock() // Unlock() -------------------------------------------------
	//
	//if time.Now().Sub(session.StartedAt) > is.maxSessionDuration {
	//	return nil, errors.Wrap(auth.ErrAuthSession, "session is expired")
	//}
	//
	//if session.IP != toAuth[auth.CredsIP] {
	//	return nil, auth.ErrIP
	//}

	signature := []byte(toAuth[auth.CredsSignature])
	if !encrlib.ECDSAVerify(keyToSignature, publKey, signature) {
		return nil, errWrongSignature
	}

	var nickname = publKeyBase58
	//if nicknameReceived := toAuth[auth.CredsNickname]; strings.TrimSpace(nicknameReceived) != "" {
	//	nickname = nicknameReceived
	//}

	return &auth.Actor{
		Identity: &auth.Identity{
			ID:       auth.ID(Proto + "://" + publKeyBase58),
			Nickname: nickname,
		},
	}, nil
}
