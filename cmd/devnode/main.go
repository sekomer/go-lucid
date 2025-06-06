package cmd

import (
	"context"
	"encoding/json"
	"go-lucid/api/routes"
	block_core "go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/database"
	"go-lucid/node"
	tx_p2p "go-lucid/p2p/transaction"
	"go-lucid/rpc/block"
	"go-lucid/rpc/ping"
	"go-lucid/service/health"
	"log"
	"net/http"
	"time"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
)

var rpcProtocolID = protocol.ID("/p2p/rpc/ping")

func StartRpcClient(client host.Host) *gorpc.Client {
	return gorpc.NewClient(client, rpcProtocolID)
}

func main(c *node.FullNodeConfig) {
	log.Println("devnode starting...")

	database.InitDB("/tmp/devnode.db")

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
		for {
			time.Sleep(3 * time.Second)

			blockClient := block.NewBlockClient(n.Host)
			blockRpcArgs := block.BlockRpcArgs{
				Method: "GetBlock",
				Args:   []any{1},
			}
			reply := block.BlockRpcReply{}
			err := blockClient.Call(context.Background(), n.Host.Network().Peers()[0], "GetBlock", &blockRpcArgs, &reply)
			if err != nil {
				log.Println("error calling get block rpc:", err)
				return
			}

			log.Println("xo dev called the peer:", n.Host.Network().Peers()[0])
			if !reply.Success {
				log.Println("xo dev success:", reply.Success)
				log.Println("xo dev error:", reply.Error)
				log.Println("xo dev result:", reply.Result)
			}

			block := block_core.Block{}
			err = json.Unmarshal(reply.Result, &block)
			if err != nil {
				log.Println("error unmarshalling block:", err)
				return
			}

			log.Println("init test: devnode block:", block)
		}
	}()

	healthService := health.NewHealthService(n.Host)
	go func() {
		err := healthService.Start(context.Background())
		if err != nil {
			log.Println("error starting health service:", err)
		}
	}()

	transactionService, err := tx_p2p.NewTransactionService(n.Host, n.PubSub)
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
		block_number := uint32(0)
		for range time.Tick(3 * time.Second) {
			err := yeniTxService.Broadcast(context.Background(), transaction.RawTransaction{
				Version:   32,
				BlockID:   block_number,
				Hash:      []byte("devnode hash"),
				TxInCount: 333333,
			})
			if err != nil {
				log.Println("error broadcasting block:", err)
			}
			block_number++
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
		}
	}()

	mux := http.NewServeMux()
	routes.RegisterTransactionRoutes(mux, transactionService)
	go func() {
		err := http.ListenAndServe(":8080", mux)
		if err != nil {
			log.Println("error starting http server:", err)
		}
	}()

	// go func() {
	// 	for range time.Tick(3 * time.Second) {
	// 		log.Println("mempool size:", mempool.GetMempool().Size())
	// 		for _, tx := range mempool.GetMempool().GetTxs() {
	// 			log.Printf("tx: %+v", tx)
	// 		}
	// 		log.Println("--------------------------------")
	// 	}
	// }()

	go func() {
		for range time.Tick(10 * time.Second) {
			peers := n.Host.Network().Peers()
			log.Println("devnode pubsub peers:", peers)
		}
	}()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
