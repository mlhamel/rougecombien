package consumer

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/mlhamel/rougecombien/pkg/config"
	"github.com/mlhamel/rougecombien/pkg/gcloud"
)

type Consumer struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Consumer {
	return &Consumer{cfg: cfg}
}

func (c *Consumer) Run(ctx context.Context) error {
	c.cfg.Logger().Info().Str("test", "testing").Msg("teeeest")

	subscription, err := gcloud.NewPubSubSubscription(ctx, c.cfg)

	if err != nil {
		return err
	}

	c.cfg.Logger().Info().Msg("Waiting for message")

	return subscription.Subscribe(ctx, c.cfg.TopicName(), func(ctx context.Context, message *pubsub.Message) {
		c.cfg.Logger().Info().
			Str("ID", message.ID).
			Time("PublishTime", message.PublishTime)
		message.Ack()
	})
}
