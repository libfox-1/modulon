package core

import (
	"github/com/libfox-1/modulon/crypto"
	"github/com/libfox-1/modulon/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0, types.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	falsePrivKey := crypto.GeneratePrivateKey()
	b.Validator = falsePrivKey.PublicKey()

	assert.NotNil(t, b.Verify())

	// tamper the height and error should be NotNil
	b.Height = 100
	assert.NotNil(t, b.Verify())
}

func randomBlock(height uint32, PrevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: PrevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	return NewBlock(header, []Transaction{})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(height, prevBlockHash)
	tx := randomTxWithSignature(t)
	b.AddTransaction(tx)
	assert.Nil(t, b.Sign(privKey))

	return b
}
