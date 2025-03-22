package p2p_test

import (
	"context"
	"go-lucid/config"
	"testing"
	"time"
)

func TestStartBootNode(t *testing.T) {
	t.Parallel()

	c := config.MustReadConfig("./config/test-init.yaml")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go StartTestBootNode(c, ctx)
	time.Sleep(1 * time.Second)
}

func TestNodeConnectsToPeers(t *testing.T) {
	t.Parallel()

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
