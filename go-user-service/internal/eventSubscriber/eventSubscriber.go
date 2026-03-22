package eventsubscriber

import (
	"context"
	"encoding/json"
	"go-user-service/internal/ws"
	"time"

	"github.com/redis/go-redis/v9"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func Listen(ctx context.Context, hub *ws.Hub, rdb *redis.Client) {
	for {
		if ctx.Err() != nil {
			return
		}

		pubsub := rdb.Subscribe(ctx, "myconfirmationChannel")

		// Receive is a kind of handshake
		if _, err := pubsub.Receive(ctx); err != nil {
			pubsub.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		eventLoop(ctx, pubsub, hub)

		pubsub.Close()
	}
}

func eventLoop(ctx context.Context, pubsub *redis.PubSub, hub *ws.Hub) {
	ch := pubsub.Channel()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}

			var e Event
			if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
				continue
			}

			data, err := json.Marshal(e)
			if err != nil {
				continue
			}

			select {
			case hub.Broadcast <- data:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
