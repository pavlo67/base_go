package auth

import (
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/encrlib"
)

const testIP = "1.2.3.4"

func OperatorTestScenarioPassword(t *testing.T, authOp Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	testCreds := []Creds{
		{
			CredsNickname: "nickname" + strconv.FormatInt(time.Now().Unix(), 10),
			CredsPassword: "password" + strconv.FormatInt(time.Now().Unix(), 10),
		},
	}

	for i, tc := range testCreds {
		nickname := tc[CredsNickname]
		password := tc[CredsPassword]

		t.Log(i, "\n")

		// .SetCredsByKey() ------------------------------------------

		userCreds, err := authOp.SetCreds(Actor{}, tc)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		// t.Logf("           creds: %#v\n\n", *userCreds)

		require.Equal(t, nickname, (*userCreds)[CredsNickname])

		// .Authenticate() ok -----------------------------------------

		userCreds = &Creds{
			// CredsIP:       testIP,
			CredsNickname: nickname,
			CredsPassword: password,
		}

		actor, err := authOp.Authenticate(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, actor)

		// t.Logf("%#v", *actor)

		require.NotNil(t, actor.Identity)
		require.Equal(t, nickname, actor.Identity.Nickname)
		require.NotEmpty(t, actor.Identity.ID)

		// .Authenticate() err ----------------------------------------

		userCreds = &Creds{
			// CredsIP:       testIP,
			CredsNickname: nickname,
			CredsPassword: password + "1",
		}

		actor, err = authOp.Authenticate(*userCreds)
		require.Error(t, err)
		require.Nil(t, actor)
	}
}

func OperatorTestScenarioToken(t *testing.T, operator Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	testCreds := []Creds{
		{
			CredsNickname: "nickname" + strconv.FormatInt(time.Now().Unix(), 10),
		},
	}

	for i, tc := range testCreds {
		t.Log(i)

		// .SetCredsByKey() ------------------------------------------

		userCreds, err := operator.SetCreds(Actor{}, tc)
		require.NoError(t, err)
		require.NotNil(t, userCreds)

		log.Printf("           creds: %#v", *userCreds)
		require.Equal(t, tc[CredsNickname], (*userCreds)[CredsNickname])

		// .Authenticate() -----------------------------------------

		// (*userCreds)[CredsIP] = testIP

		actor, err := operator.Authenticate(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, actor)
		require.NotNil(t, actor.Identity)
		require.Equal(t, tc[CredsNickname], actor.Identity.Nickname)
	}
}

func OperatorTestScenarioPublicKey(t *testing.T, operator Operator) {
	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	testCreds := []Creds{{}}

	for i, tc := range testCreds {
		t.Log(i)

		// .SetCredsByKey() ------------------------------------------

		userCreds, err := operator.SetCreds(Actor{}, tc)
		require.NoError(t, err)
		require.NotNil(t, userCreds)
		require.NotEmpty(t, (*userCreds)[CredsPublicKeyBase58])
		require.NotEmpty(t, (*userCreds)[CredsPublicKeyBase58], (*userCreds)[CredsNickname])

		log.Printf("            creds: %#v", userCreds)

		// require.Equal(t, tc[CredsNickname], userCreds.StringDefault(CredsNickname, ""))
		// nickname := (*userCreds)[CredsNickname]

		// .InitAuth() -----------------------------------

		privKeySerialization := []byte((*userCreds)[CredsPrivateKey])
		privKey, err := encrlib.ECDSADeserialize(privKeySerialization)
		require.NoError(t, err)
		require.NotNil(t, privKey)

		publicKeyBase58 := (*userCreds)[CredsPublicKeyBase58]
		log.Printf("public key base58: %s", publicKeyBase58)

		credsToSet := Creds{CredsToSet: string(CredsKeyToSignature)} // CredsIP: testIP,
		sessionCreds, err := operator.SetCreds(Actor{}, credsToSet)
		require.NoError(t, err)
		require.NotNil(t, sessionCreds)

		// ---------------------------------------------------------------------

		keyToSignature := testIP

		//log.Printf(" key to signature: %s", keyToSignature)
		//require.True(t, len(keyToSignature) > 0)

		signature, err := encrlib.ECDSASign(keyToSignature, *privKey)
		require.NoError(t, err)
		require.True(t, len(signature) > 0)

		log.Printf("      private key: %s", privKeySerialization)
		log.Printf("        signature: %s", base58.Encode(signature))

		publKey := base58.Decode(publicKeyBase58)
		ok := encrlib.ECDSAVerify(keyToSignature, publKey, signature)
		require.True(t, ok)

		// .Authenticate() -----------------------------------------

		(*userCreds)[CredsIP] = testIP
		(*userCreds)[CredsKeyToSignature] = keyToSignature
		(*userCreds)[CredsSignature] = string(signature)

		actor, err := operator.Authenticate(*userCreds)

		require.NoError(t, err)
		require.NotNil(t, actor)
		require.NotNil(t, actor.Identity)
		// require.Equal(t, nickname, actor.Creds[CredsNickname])
		require.NotEmpty(t, actor.Identity.ID)
	}
}
