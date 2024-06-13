package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyPair_Sign_Verify_Success(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("Testing Message")

	sig, err := privKey.Sign(msg)
	assert.Nil(t, err)

	assert.True(t, sig.Verify(pubKey, msg))

}

func TestKeyPair_Sign_Verify_Fail(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey()
	msg := []byte("Testing Message")

	sig, err := privKey.Sign(msg)

	// No error, Signing Works
	assert.Nil(t, err)

	falsePrivKey := GeneratePrivateKey()
	falsePubKey := falsePrivKey.PublicKey()

	// Tries to Verify with a different Public Key => False
	assert.False(t, sig.Verify(falsePubKey, msg))

	// Tries to verify a different message => False
	assert.False(t, sig.Verify(pubKey, []byte("Different Test Message")))

}
