package p2p_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go-lucid/config"
	coreBlock "go-lucid/core/block"
	"go-lucid/rpc/block"
	"log"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

func TestStartBootNode(t *testing.T) {
	c := config.MustReadConfig("./config/test-init.yaml")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartTestBootNode(c, ctx)
	time.Sleep(1 * time.Second)
}

func TestNodeConnectsToPeers(t *testing.T) {
	c := config.MustReadConfig("./config/test-boot.yaml")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartTestBootNode(c, ctx)
	time.Sleep(1 * time.Second)

	c = config.MustReadConfig("./config/dev.yaml")
	n, cancel := CreateTestNode(c)
	defer cancel()

	peerCount := len(n.Host.Network().Peers())
	t.Logf("peer count: %d\n", peerCount)
	t.Logf("peers: %v\n", n.Host.Network().Peers())
}

func TestNodeRpc(t *testing.T) {
	c := config.MustReadConfig("./config/test-boot.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartTestBootNode(c, ctx)
	time.Sleep(1 * time.Second)

	c = config.MustReadConfig("./config/dev.yaml")
	n, cancel := CreateTestNode(c)
	defer cancel()

	blockClient := block.NewBlockClient(n.Host)

	// get random peer
	peers := n.Host.Network().Peers()
	var randomPeer peer.ID
	for _, p := range peers {
		if p != n.Host.ID() {
			randomPeer = p
			break
		}
	}

	args := block.BlockRpcArgs{
		Method: "GetBlock",
		Args:   []any{1},
	}
	reply := block.BlockRpcReply{}
	err := blockClient.Call(context.TODO(), randomPeer, block.MethodBlock, &args, &reply)
	if err != nil {
		t.Fatalf("failed to get block: %v", err)
	}

	log.Println("rpc call success:", reply.Success)
	log.Println("rpc call error:", reply.Error)
	log.Println("rpc call result:", reply.Result)

	res1 := coreBlock.Block{}
	err = json.Unmarshal(reply.Result, &res1)
	if err != nil {
		t.Fatalf("failed to unmarshal block: %v", err)
	}
	prettyJSON, err := json.MarshalIndent(res1, "", "    ")
	if err != nil {
		t.Fatalf("failed to marshal block for pretty printing: %v", err)
	}
	fmt.Printf("block: \n%s\n", string(prettyJSON))
}
