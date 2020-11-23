package gcloud

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/mlhamel/rougecombien/pkg/config"
)

type PubSubSubscription struct {
	cfg    *config.Config
	client *pubsub.Client
}

func NewPubSubSubscription(ctx context.Context, cfg *config.Config) (*PubSubSubscription, error) {
	cfg.Logger().Info().
		Str("PUBSUB_EMULATOR_HOST", os.Getenv("PUBSUB_EMULATOR_HOST")).
		Msg("Initializing gcloud client")
	client, err := pubsub.NewClient(ctx, cfg.GoogleCloudProject())

	if err != nil {
		return nil, err
	}

	subscription := PubSubSubscription{
		cfg:    cfg,
		client: client,
	}

	return &subscription, nil
}

func (p *PubSubSubscription) FindSubscription(ctx context.Context, topicName string) (*pubsub.Subscription, error) {
	it := p.client.Subscriptions(ctx)

	for {
		subscription, err := it.Next()

		if err != nil {
			return nil, nil
		}

		if subscription.ID() == topicName {
			return subscription, nil
		}
	}

	return nil, nil
}

func (p *PubSubSubscription) CreateSubscription(ctx context.Context, topicName string) (*pubsub.Subscription, error) {
	topic, err := initializePubsubTopic(ctx, p.client, topicName)

	if err != nil {
		return nil, fmt.Errorf("could not open topic subscription: %v", err)
	}

	sub, err := p.client.CreateSubscription(ctx, topicName, pubsub.SubscriptionConfig{
		Topic:       topic,
		AckDeadline: 10 * time.Second,
	})

	return sub, err
}

func (p *PubSubSubscription) Subscribe(ctx context.Context, topicName string, callback func(context.Context, *pubsub.Message)) error {
	var sub *pubsub.Subscription

	sub, err := p.FindSubscription(ctx, topicName)

	if err != nil {
		return fmt.Errorf("could not search for subscription: %v", err)
	}

	if sub == nil {
		sub, err = p.CreateSubscription(ctx, topicName)

		if err != nil {
			return fmt.Errorf("could not create subscription: %v", err)
		}
	}

	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		p.cfg.Logger().Info().Str("Msg", string(msg.Data)).Msg("Received message")
		msg.Ack()
	})
}
