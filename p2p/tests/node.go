package p2p_test

import (
	"context"
	"go-lucid/node"
	"log"
	"time"

	"golang.org/x/exp/rand"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func StartTestBootNode(c *node.FullNodeConfig, ctx context.Context) {
	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(uint64(0))))
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	n := node.CreateHost(priv, c)
	defer n.Close()

	n.InitPeers()

	log.Println("Test node started")
	log.Printf("Test node ID: %s\n", n.Host.ID())

	<-ctx.Done()
	log.Println("Test node shutting down")
}

func CreateTestNode(c *node.FullNodeConfig) (*node.Node, context.CancelFunc) {
	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(uint64(time.Now().UnixNano()))))
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	n := node.CreateHost(priv, c)
	n.InitPeers()

	return n, func() {
		n.Close()
		log.Println("Test node closed")
	}
}
