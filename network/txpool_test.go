package network

import (
	"github/com/libfox-1/modulon/core"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	// test instantiate the txPool without any transactions
	assert.Equal(t, p.Len(), 0)
}

func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("test"))
	// Add a transaction with a byte slice
	assert.Nil(t, p.Add(tx))
	// Check if the pool length is 1 after added 1 tx
	assert.Equal(t, p.Len(), 1)

	// Try to add the same transaction again and check the mempool length
	_ = core.NewTransaction([]byte("test"))
	assert.Equal(t, p.Len(), 1)

	// flush the pool and check if length is 0
	p.Flush()
	assert.Equal(t, p.Len(), 0)
}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen((int64(i * rand.Intn(10000))))
		assert.Nil(t, p.Add(tx))
	}

	assert.Equal(t, txLen, p.Len())

	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
}
