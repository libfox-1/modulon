package core

import (
	"github/com/libfox-1/modulon/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	// assert if the validator exists
	assert.NotNil(t, bc.validator)
	//assert if the blockcahin height is 0 at instantiation
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	// asserting the genesis block
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}
	// Adding blocks and testing correct height
	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	// Because we start at 0 but we count from 1 test this
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	// Testing to add a new block into a Block slot that is already filled (see loop)
	assert.NotNil(t, bc.AddBlock(randomBlock(89, types.Hash{})))
}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

func TestAddBlockTooHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	// Gets the genesis block hash "0" and adds the first non genesis block "1"
	assert.Nil(t, bc.AddBlock(randomBlockWithSignature(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	// Should fail if we add a block that's too high
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 3, types.Hash{})))
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return BlockHasher{}.Hash(prevHeader)
}
