package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	tr_a := NewLocalTransport("A")
	tr_b := NewLocalTransport("B")

	tr_a.Connect(tr_b)
	tr_b.Connect(tr_a)

	//assert.Equal(t, tr_a.peers[tr_b.Addr()], tr_b)
	//assert.Equal(t, tr_b.peers[tr_a.Addr()], tr_a)
	//assert.Equal(t, 1, 1)
}

func TestSendMessage(t *testing.T) {
	tr_a := NewLocalTransport("A")
	tr_b := NewLocalTransport("B")

	tr_a.Connect(tr_b)
	tr_b.Connect(tr_a)

	msg := []byte("hello world!")
	assert.Nil(t, tr_a.SendMessage(tr_b.Addr(), msg))

	rpc := <-tr_b.Consume()
	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, n, len(msg))

	assert.Equal(t, buf, msg)
	assert.Equal(t, rpc.From, tr_a.Addr())
}
