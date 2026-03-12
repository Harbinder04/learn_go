package eventsubscriber

import (
	"context"
	"encoding/json"
	"go-user-service/internal/ws"

	"github.com/redis/go-redis/v9"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func Listen(ctx context.Context, hub *ws.Hub, rdb *redis.Client) {
	//todo: remove empty context and use actuall context
	pubsub := rdb.Subscribe(ctx, "myconfirmationChannel")

	defer pubsub.Close()
	// verifying subscription is created or not
	_, err := pubsub.Receive(ctx)
	if err != nil {
		panic(err)
	}

	ch := pubsub.Channel()

	for {
		for msg := range ch {
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
		}
	}
}
