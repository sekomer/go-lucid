package p2p_test

import (
	"fmt"
	"go-lucid/node"
	"io"
	"os"
	"testing"
	"time"

	bootnode "go-lucid/cmd/bootnode"

	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v3"
)

func TestNodeInit(t *testing.T) {
	c := &node.FullNodeConfig{
		Node: node.NodeConfig{
			Id:      "test-node",
			Type:    node.NodeType(node.FullNode),
			Network: "test",
			Rpc: node.RpcConfig{
				Enabled: true,
				Port:    17289,
				Cors:    "*",
				Apis:    []string{"/ping"},
			},
			P2p: node.P2pConfig{
				ListenPort:  17289,
				MinPeers:    1,
				MaxPeers:    10,
				GracePeriod: 10,
			},
			Data: node.DataConfig{
				Dir: "./data",
			},
			Logging: node.LoggingConfig{
				Level: "debug",
				File:  "./logs/node.log",
			},
			Sync: node.SyncConfig{
				Mode: "full",
			},
			Mining: node.MiningConfig{
				Enabled:  true,
				MinerUrl: "http://localhost:8080",
			},
			Genesis: node.GenesisConfig{
				File: "./genesis.json",
			},
			Peers: []string{},
			Debug: node.DebugConfig{
				Enabled: false,
			},
		},
	}

	fmt.Println("bootnode starting...")

	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(c.Node.Debug.Seed)))
	if err != nil {
		panic(err)
	}

	n := node.CreateHost(priv, c)
	defer n.Close()

	n.InitPeers()

	fmt.Printf("[lucid-go], hosts ID is %s\n", n.Host.ID())
	fmt.Printf("connection address of this node is: %s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())
}

func TestTwoNodesCommunication(t *testing.T) {
	var c *node.FullNodeConfig
	var buf []byte
	var err error

	yamlReader, err := os.Open("./config/boot.yaml")
	if err != nil {
		panic(err)
	}
	defer yamlReader.Close()

	buf, err = io.ReadAll(yamlReader)
	if err != nil {
		panic(err)
	}

	c = &node.FullNodeConfig{}
	c.Node.Peers = []string{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		panic(err)
	}

	go bootnode.Start(c)
	time.Sleep(2 * time.Second)

	yamlReader, err = os.Open("./config/dev.yaml")
	if err != nil {
		panic(err)
	}
	defer yamlReader.Close()

	buf, err = io.ReadAll(yamlReader)
	if err != nil {
		panic(err)
	}

	c = &node.FullNodeConfig{}
	c.Node.Peers = []string{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		panic(err)
	}

	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(uint64(time.Now().UnixNano()))))
	if err != nil {
		panic(err)
	}
	devnode := node.CreateHost(priv, c)
	defer devnode.Close()

	fmt.Println("[devnode] started")
	fmt.Println("[devnode] connection address of this node is: ", devnode.Host.Addrs()[0], "/p2p/", devnode.Host.ID())

	devnode.InitPeers()

	fmt.Printf("[lucid-go], hosts ID is %s\n", devnode.Host.ID())
	fmt.Printf("connection address of this node is: %s/p2p/%s\n", devnode.Host.Addrs()[0], devnode.Host.ID())

	// check if the devnode is connected to the bootnode
	peerCount := len(devnode.Host.Network().Peers())
	if peerCount == 0 {
		t.Fatalf("devnode is not connected to the bootnode")
	}
}
