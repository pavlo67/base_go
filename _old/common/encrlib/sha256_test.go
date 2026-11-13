package encrlib

import (
	"testing"

	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/sha256_crypt"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/strlib"
)

var passwordMinLength = 6
var testPassword = strlib.RandomString(passwordMinLength)
var testPasswordBad = strlib.RandomString(passwordMinLength - 1)

const testSalt = "$5$1234"

func TestGetEncodedPassword(t *testing.T) {

	crypt1 := crypt.SHA256.New()

	hash1, err := crypt1.Generate([]byte(testPasswordBad), []byte(testSalt+"a"))
	require.NoError(t, err)
	require.NotEmpty(t, hash1)

	hash2, err := crypt1.Generate([]byte(testPasswordBad), []byte(testSalt+"a"))
	require.NoError(t, err)
	require.NotEmpty(t, hash2)

	// test ok

	crypt2 := crypt.SHA256.New()

	err = crypt2.Verify(hash1, []byte(testPasswordBad))
	require.NoError(t, err)

	err = crypt2.Verify(hash2, []byte(testPasswordBad))
	require.NoError(t, err)

	// test wrong password

	err = crypt2.Verify(hash1, []byte(testPasswordBad+"!"))
	require.Error(t, err)

	err = crypt2.Verify(hash1, []byte(testPasswordBad[1:]))
	require.Error(t, err)

	err = crypt2.Verify(hash2, []byte(testPasswordBad+"!"))
	require.Error(t, err)

	err = crypt2.Verify(hash2, []byte(testPasswordBad[1:]))
	require.Error(t, err)

}
