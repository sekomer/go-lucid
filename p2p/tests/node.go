package p2p_test

import (
	"context"
	coreBlock "go-lucid/core/block"
	"go-lucid/database"
	"go-lucid/node"
	"go-lucid/rpc/block"
	"log"
	"time"

	"golang.org/x/exp/rand"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func StartTestBootNode(c *node.FullNodeConfig, ctx context.Context) {
	database.InitDB("/tmp/bootnode.db")
	db := database.GetDB()

	test_block := coreBlock.BlockModel{
		BlockHeaderModel: coreBlock.BlockHeaderModel{
			PrevBlock: []byte("test"),
			Height:    1,
		},
	}
	db.Save(&test_block)

	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(uint64(0))))
	if err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	n := node.CreateHost(priv, c)
	defer n.Close()

	n.InitPeers()

	blockService := block.NewBlockService(n.Host)
	if err := n.Rpc.RegisterService(blockService, block.ProtocolID); err != nil {
		log.Fatalf("failed to register block service: %v", err)
	}

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
