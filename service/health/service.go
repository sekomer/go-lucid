package health

import (
	"context"
	"crypto/rand"
	"go-lucid/rpc/ping"
	"go-lucid/service"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
)

type HealthService struct {
	service.BaseService
	pingClient *ping.PingClient
}

func NewHealthService(h host.Host) *HealthService {
	return &HealthService{
		BaseService: *service.NewBaseService(h, "HealthService"),
		pingClient:  ping.NewPingClient(h),
	}
}

func (hs *HealthService) Start(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			for _, peer := range hs.GetHost().Network().Peers() {
				if peer == hs.GetHost().ID() {
					continue
				}
				data := make([]byte, 128)
				if _, err := rand.Read(data); err != nil {
					return err
				}
				if err := hs.pingClient.Call(ctx, peer, "Ping", &ping.PingArgs{Data: data}, &ping.PingReply{}); err != nil {
					hs.GetLogger().Println("Ping Error:", err)
				} else {
					// hs.GetLogger().Println("Ping reply from:", peer)
				}
			}
		}
	}
}
