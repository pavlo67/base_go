package encrlib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"log"
	"testing"

	"github.com/btcsuite/btcutil/base58"

	"fmt"

	"github.com/stretchr/testify/require"
)

var keyToSignature = "test data"

func TestECDSASign(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	signature, err := ECDSASign(keyToSignature, *privKey)
	require.NoError(t, err)

	// publKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	publKey := ECDSAPublicKey(*privKey)
	fmt.Print("ADDRESS: ", base58.Encode(publKey), "\n\n")

	ok := ECDSAVerify(keyToSignature, publKey, signature)
	require.True(t, ok)

	publKeyBad := append(privKey.PublicKey.Y.Bytes(), privKey.PublicKey.X.Bytes()...)
	ok = ECDSAVerify(keyToSignature, publKeyBad, signature)
	require.False(t, ok)
	ok = ECDSAVerify(keyToSignature, []byte{}, signature)
	require.False(t, ok)

	ok = ECDSAVerify(keyToSignature+" ", publKey, signature)
	require.False(t, ok)
	ok = ECDSAVerify(keyToSignature[:len(keyToSignature)-1], publKey, signature)
	require.False(t, ok)
	ok = ECDSAVerify(keyToSignature[:len(keyToSignature)-1]+" ", publKey, signature)
	require.False(t, ok)
	ok = ECDSAVerify("", publKey, signature)
	require.False(t, ok)

	ok = ECDSAVerify(keyToSignature, publKey, []byte(string(signature)+" "))
	require.False(t, ok)
	ok = ECDSAVerify(keyToSignature, publKey, signature[:len(keyToSignature)-1])
	require.False(t, ok)
	signatureBad := make([]byte, len(signature))
	for i, s := range signature[:len(signature)-1] {
		signatureBad[i] = s
	}
	signatureBad[len(signatureBad)-1] = byte(uint8(signatureBad[len(signatureBad)-1]) + 1%256)

	ok = ECDSAVerify(keyToSignature, publKey, signatureBad)
	require.False(t, ok)
}

func TestECDSASignPredefinedPair(t *testing.T) {
	privKeyPredefined, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privKeyPredefined)

	publKeyPredefined := ECDSAPublicKey(*privKeyPredefined)
	publicKeyBase58 := base58.Encode(publKeyPredefined)

	privKeySerialization, err := ECDSASerialize(*privKeyPredefined)
	require.NoError(t, err)

	// log.Printf("private key privKeySerialization: %s", privKeySerialization)

	privKey, err := ECDSADeserialize(privKeySerialization)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	signature, err := ECDSASign(keyToSignature, *privKey)
	require.NoError(t, err)

	// publKey := ECDSAPublicKey(*privKey)

	publKey := base58.Decode(publicKeyBase58)
	require.Equal(t, publKeyPredefined, publKey)

	//log.Printf("predefined address: %s", Base58Encode(publKeyPredefined))
	log.Printf("     private key: %s", privKeySerialization)
	log.Printf("         address: %s", publicKeyBase58)
	log.Printf("key to signature: %s", keyToSignature)
	log.Printf("       signature: %s", base58.Encode(signature))

	ok := ECDSAVerify(keyToSignature, publKey, signature)
	require.True(t, ok)
}

func TestECDSASerialize(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	signature, err := ECDSASign(keyToSignature, *privKey)
	require.NoError(t, err)
	require.NotEmpty(t, signature)

	log.Printf("signature 1: %s\n\n", string(signature))

	publKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	ok := ECDSAVerify(keyToSignature, publKey, signature)
	require.True(t, ok)

	serialization, err := ECDSASerialize(*privKey)
	require.NoError(t, err)

	// log.Printf("private key serialization: %s", serialization)

	privKeyRestored, err := ECDSADeserialize(serialization)
	require.NoError(t, err)
	require.NotNil(t, privKeyRestored)

	signatureRestored, err := ECDSASign(keyToSignature, *privKeyRestored)
	require.NoError(t, err)

	log.Printf("signature 2: %s\n\n", string(signatureRestored))

	publKeyRestored := append(privKeyRestored.PublicKey.X.Bytes(), privKeyRestored.PublicKey.Y.Bytes()...)
	ok = ECDSAVerify(keyToSignature, publKeyRestored, signature)
	require.True(t, ok)
}
