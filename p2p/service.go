package p2p

import (
	"context"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type Message struct {
	From    string
	Payload []byte
	// Signature []byte
}

type P2PService interface {
	Publish(ctx context.Context, data []byte) error
	Subscribe(ctx context.Context) (<-chan Message, error)
	GetHost() host.Host
	GetPubSub() *pubsub.PubSub
	GetTopic() *pubsub.Topic
	GetSubscription() *pubsub.Subscription
	GetChannel() string
	Close() error
	Name() string
}

type BaseService struct {
	host    host.Host
	ps      *pubsub.PubSub
	topic   *pubsub.Topic
	sub     *pubsub.Subscription
	channel string
	name    string
}

// Ensure BaseService implements P2PService interface
var _ P2PService = (*BaseService)(nil)

func NewBaseService(name string, h host.Host, ps *pubsub.PubSub, channel string) (*BaseService, error) {
	topic, err := ps.Join(channel)
	if err != nil {
		return nil, err
	}

	sub, err := topic.Subscribe()
	if err != nil {
		return nil, err
	}

	return &BaseService{
		host:    h,
		ps:      ps,
		topic:   topic,
		sub:     sub,
		channel: channel,
	}, nil
}

func (s *BaseService) Publish(ctx context.Context, data []byte) error {
	return s.topic.Publish(ctx, data)
}

func (s *BaseService) Subscribe(ctx context.Context) (<-chan Message, error) {
	ch := make(chan Message, 1024)

	go func() {
		defer close(ch)
		for {
			msg, err := s.sub.Next(ctx)
			if err != nil {
				return
			}
			if msg.ReceivedFrom == s.host.ID() {
				continue
			}
			ch <- Message{
				From:    msg.ReceivedFrom.String(),
				Payload: msg.Data,
			}
		}
	}()

	return ch, nil
}

func (s *BaseService) Close() error {
	s.sub.Cancel()
	return s.topic.Close()
}

// GetHost returns the host of the BaseService
func (s *BaseService) GetHost() host.Host {
	return s.host
}

// GetPubSub returns the PubSub instance of the BaseService
func (s *BaseService) GetPubSub() *pubsub.PubSub {
	return s.ps
}

// GetTopic returns the Topic of the BaseService
func (s *BaseService) GetTopic() *pubsub.Topic {
	return s.topic
}

// GetSubscription returns the Subscription of the BaseService
func (s *BaseService) GetSubscription() *pubsub.Subscription {
	return s.sub
}

// GetChannel returns the channel name of the BaseService
func (s *BaseService) GetChannel() string {
	return s.channel
}

// GetName returns the name of the BaseService
func (s *BaseService) Name() string {
	return s.name
}
