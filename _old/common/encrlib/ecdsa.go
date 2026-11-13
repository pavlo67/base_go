package encrlib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"io"
	"math/big"
)

func ECDSAPublicKey(privKey ecdsa.PrivateKey) []byte {
	return append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
}

func ECDSASign(data string, privKey ecdsa.PrivateKey) ([]byte, error) {
	h := md5.New()
	io.WriteString(h, data)

	r, s, err := ecdsa.Sign(rand.Reader, &privKey, h.Sum(nil))
	if err != nil {
		return nil, err
	}

	return append(r.Bytes(), s.Bytes()...), nil
}

func ECDSAVerify(data string, publKey, signature []byte) bool {
	h := md5.New()
	io.WriteString(h, data)
	dataSum := h.Sum(nil)

	// build key and verify dataSum
	sigLen := len(signature)

	s := big.Int{}
	s.SetBytes(signature[(sigLen / 2):])

	r := big.Int{}
	r.SetBytes(signature[:(sigLen / 2)])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(publKey)
	x.SetBytes(publKey[:(keyLen / 2)])
	y.SetBytes(publKey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: &x, Y: &y}

	return ecdsa.Verify(&rawPubKey, dataSum, &r, &s)
}

func ECDSASerialize(privKey ecdsa.PrivateKey) ([]byte, error) {
	privKey.Curve = nil

	return json.Marshal(privKey)

	//var encoded bytes.Buffer
	// gob.Register(privKey.Curve)
	//enc := gob.NewEncoder(&encoded)
	//err := enc.Encode(privKey)
	//if err != nil {
	//	return nil, err
	//}
	//return encoded.Bytes(), nil
}

func ECDSADeserialize(data []byte) (*ecdsa.PrivateKey, error) {
	//decoder := gob.NewDecoder(bytes.NewReader(data))
	//err := decoder.Decode(privKey)
	//if err != nil {
	//	return err
	//}

	privKey := ecdsa.PrivateKey{}

	err := json.Unmarshal(data, &privKey)
	if err != nil {
		return nil, err
	}

	// it's not necessary if privKey.Curve would be serialized instead "nilled"
	privKey.Curve = elliptic.P256()

	return &privKey, nil
}
