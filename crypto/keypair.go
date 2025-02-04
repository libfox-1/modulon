package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"

	"github/com/libfox-1/modulon/types"
)

type PrivateKey struct {
	Key *ecdsa.PrivateKey
}

func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.Key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R: r,
		S: s,
	}, nil
	//return &Signature{r, s}, nil
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		Key: key,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())

	return types.AddressFromBytes(h[len(h)-20:])
}

func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &k.Key.PublicKey,
	}
}

type Signature struct {
	R, S *big.Int
}

func (sig Signature) Verify(pubkey PublicKey, data []byte) bool {
	return ecdsa.Verify(pubkey.Key, data, sig.R, sig.S)
}
