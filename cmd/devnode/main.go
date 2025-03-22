package cmd

import (
	"context"
	"encoding/json"
	block_core "go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/node"
	tx_p2p "go-lucid/p2p/transaction"
	"go-lucid/rpc/block"
	"go-lucid/rpc/ping"
	"go-lucid/service/health"
	"log"
	"time"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

var rpcProtocolID = protocol.ID("/p2p/rpc/ping")

func StartRpcClient(client host.Host) *gorpc.Client {
	return gorpc.NewClient(client, rpcProtocolID)
}

func main(c *node.FullNodeConfig) {
	log.Println("dev node starting...")

	n := node.CreateHost(nil, c)
	n.InitPeers()
	defer n.Close()

	log.Printf("Hello World, hosts ID is %s\n", n.Host.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())

	pingService := ping.NewPingService(n.Host)
	err := n.Rpc.RegisterService(pingService, ping.ProtocolID)
	if err != nil {
		panic(err)
	}

	blockService := block.NewBlockService(n.Host)
	err = n.Rpc.RegisterService(blockService, block.ProtocolID)
	if err != nil {
		panic(err)
	}
	go func() {
		for range time.Tick(3 * time.Second) {
			// get a random peer thats not the same as the current node
			peers := n.Host.Network().Peers()
			var randomPeer peer.ID
			for _, peer := range peers {
				if peer != n.Host.ID() {
					randomPeer = peer
					break
				}
			}
			if randomPeer == "" {
				log.Println("no peers found")
				continue
			}

			blockClient := block.NewBlockClient(n.Host)
			blockRpcArgs := block.GetBlockRpcArgs{
				Method: "GetBlock",
				Args:   []interface{}{5, 10},
			}
			blockRpcReply := block.GetBlockRpcReply{}
			err := blockClient.Call(context.Background(), randomPeer, "GetBlock", &blockRpcArgs, &blockRpcReply)
			if err != nil {
				log.Println("error calling get block rpc:", err)
				continue
			}
			log.Printf("[ >>>>>>>>>> GET BLOCK RPC CALL RESULT]")

			// 2025/03/04 00:57:03 main.go:77: get block result: map[Bits:0 Hash:<nil> Height:123456789 MerkleRoot:<nil> Nonce:99999999 PrevBlock:[112 114 101 118 32 98 108 111 99 107] Timestamp:<nil> TxCount:0 Txs:<nil> Version:1]
			// how to get block type from this result?
			block := block_core.Block{}
			err = json.Unmarshal(blockRpcReply.Result, &block)
			if err != nil {
				log.Println("error unmarshalling block:", err)
				continue
			}
			log.Printf("get block result: %v", block)
			log.Printf("get block result nonce: %v", block.Nonce)
		}
	}()

	healthService := health.NewHealthService(n.Host)
	go func() {
		err := healthService.Start(context.Background())
		if err != nil {
			log.Println("error starting health service:", err)
		}
	}()

	ps, err := pubsub.NewGossipSub(context.Background(), n.Host)
	if err != nil {
		panic(err)
	}
	transactionService, err := tx_p2p.NewTransactionService(n.Host, ps)
	if err != nil {
		panic(err)
	}
	n.AddService(transactionService)

	yeniTxService := n.GetService(transactionService.Name()).(*tx_p2p.TransactionService)

	ch, err := yeniTxService.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		for range time.Tick(3 * time.Second) {
			log.Println("broadcasting block... devnode")
			err := yeniTxService.Broadcast(context.Background(), transaction.RawTransaction{Version: 32})
			if err != nil {
				log.Println("error broadcasting block:", err)
			}
		}
	}()
	go func() {
		for msg := range ch {
			tx := transaction.RawTransaction{}
			err := tx.Deserialize(msg.Payload)
			if err != nil {
				log.Println("error deserializing block:", err)
				continue
			}
			log.Printf("from: %s, block: %+v\n", msg.From, tx)
		}
	}()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
