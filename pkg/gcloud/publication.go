package gcloud

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/mlhamel/rougecombien/pkg/config"
)

type PubSubEmission struct {
	cfg    *config.Config
	client *pubsub.Client
}

func NewPubSubEmission(ctx context.Context, cfg *config.Config) (*PubSubEmission, error) {
	client, err := pubsub.NewClient(ctx, cfg.GoogleCloudProject())

	if err != nil {
		return nil, err
	}

	return &PubSubEmission{cfg: cfg, client: client}, nil
}

func (p *PubSubEmission) Publish(ctx context.Context, topicName string, data []byte) error {
	topic, err := initializePubsubTopic(ctx, p.client, topicName)

	if err != nil {
		return err
	}

	msg := pubsub.Message{Data: data}

	p.cfg.Logger().
		Info().
		Str("TopicName", topicName).
		Time("EmitedAt", time.Now().UTC()).
		Msg("Emitting message to pubsub")

	_, err = topic.Publish(ctx, &msg).Get(ctx)
	return err
}
