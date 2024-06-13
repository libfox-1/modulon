package network

import (
	"fmt"
	"time"

	"core"
	"crypto"

	"github.com/sirupsen/logrus"
)

var defaultBlockTime = 5 * time.Second

// ServerOpts contains options for configuring the Server.
type ServerOpts struct {
	RPCDecodeFunc RPCDecodeFunc      // Function to decode RPC messages
	RPCProcessor  RPCProcessor       // Processor for handling RPC messages
	Transports    []Transport        // Network transports
	PrivateKey    *crypto.PrivateKey // Validator node's private key (if applicable)
	BlockTime     time.Duration      // Time interval between block creation
}

// Server represents a node in the network.
type Server struct {
	ServerOpts
	memPool     *TxPool       // Transaction pool
	isValidator bool          // Whether the node is a validator
	rpcCh       chan RPC      // Channel for receiving RPC messages
	quitCh      chan struct{} // Channel for shutting down the server
}

// NewServer creates and initializes a new Server with the given options.
func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == 0 {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	s := &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	return s
}

// Start begins the server's operation.
func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.BlockTime)
	defer ticker.Stop()

	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error("Failed to decode RPC:", err)
				continue
			}
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				logrus.Error("Failed to process message:", err)
			}
		case <-s.quitCh:
			fmt.Println("Server shutdown")
			return
		case <-ticker.C:
			if s.isValidator {
				if err := s.CreateNewBlock(); err != nil {
					logrus.Error("Failed to create new block:", err)
				}
			}
		}
	}
}

// ProcessMessage handles incoming messages based on their type.
func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	default:
		logrus.Warn("Unknown message type")
		return fmt.Errorf("unknown message type: %T", t)
	}
}

// processTransaction processes incoming transactions.
func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Has(hash) {
		logrus.WithField("hash", hash).Info("Transaction already in mempool")
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	logrus.WithFields(logrus.Fields{
		"hash":           hash,
		"mempool length": s.memPool.Len(),
	}).Info("Adding new transaction to mempool")

	return s.memPool.Add(tx)
}

// CreateNewBlock creates a new block (consensus logic).
func (s *Server) CreateNewBlock() error {
	fmt.Println("Creating a new block")
	return nil
}

// initTransports initializes the network transports.
func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
